use crate::database::Database;
use crate::error::ProcessorError;
use crate::models::{Video, VideoScript, Scene, RenderSettings, ProcessingJob};
use crate::services::S3Service;
use anyhow::Result;
use std::path::{Path, PathBuf};
use std::process::Command;
use std::sync::Arc;
use tokio::fs;
use tracing::{info, error, warn};
use uuid::Uuid;
use chrono::Utc;
use rayon::prelude::*;

#[derive(Debug)]
pub struct VideoProcessor {
    database: Database,
    s3_service: S3Service,
    temp_dir: PathBuf,
}

impl VideoProcessor {
    pub fn new(database: Database, s3_service: S3Service) -> Self {
        let temp_dir = std::env::temp_dir().join("assos_video_processor");
        
        Self {
            database,
            s3_service,
            temp_dir,
        }
    }

    pub async fn initialize(&self) -> Result<()> {
        // Create temp directory
        fs::create_dir_all(&self.temp_dir).await?;
        info!("Initialized video processor with temp dir: {:?}", self.temp_dir);
        Ok(())
    }

    pub async fn process_video(&self, video_id: Uuid) -> Result<()> {
        info!("Starting video processing for video: {}", video_id);

        // Get video from database
        let video = self.get_video(video_id).await?;
        if video.script.is_none() {
            return Err(ProcessorError::Custom("Video script not found".to_string()).into());
        }

        // Parse video script
        let script: VideoScript = serde_json::from_value(video.script.unwrap())?;
        
        // Update video status
        self.update_video_status(video_id, "producing").await?;

        // Create processing job
        let job_id = self.create_processing_job(video_id, "video_assembly").await?;

        match self.process_video_internal(&video, &script).await {
            Ok(video_url) => {
                // Update video with final URL
                self.update_video_url(video_id, &video_url).await?;
                self.update_video_status(video_id, "completed").await?;
                self.complete_processing_job(job_id, None).await?;
                
                info!("Video processing completed: {}", video_id);
            }
            Err(e) => {
                error!("Video processing failed: {}", e);
                self.update_video_status(video_id, "failed").await?;
                self.complete_processing_job(job_id, Some(&e.to_string())).await?;
                return Err(e);
            }
        }

        Ok(())
    }

    async fn process_video_internal(&self, video: &Video, script: &VideoScript) -> Result<String> {
        let video_dir = self.temp_dir.join(video.id.to_string());
        fs::create_dir_all(&video_dir).await?;

        // Generate scenes in parallel
        info!("Processing {} scenes", script.scenes.len());
        let scene_paths: Result<Vec<_>> = script.scenes
            .par_iter()
            .enumerate()
            .map(|(i, scene)| self.process_scene(scene, i, &video_dir))
            .collect();

        let scene_paths = scene_paths?;

        // Combine scenes into final video
        let output_path = video_dir.join("final_video.mp4");
        self.combine_scenes(&scene_paths, &output_path, script).await?;

        // Upload to S3
        let s3_key = format!("videos/{}/final.mp4", video.id);
        let video_url = self.s3_service.upload_file(&output_path, &s3_key).await?;

        // Cleanup temp files
        if let Err(e) = fs::remove_dir_all(&video_dir).await {
            warn!("Failed to cleanup temp directory: {}", e);
        }

        Ok(video_url)
    }

    fn process_scene(&self, scene: &Scene, index: usize, video_dir: &Path) -> Result<PathBuf> {
        let scene_path = video_dir.join(format!("scene_{:03d}.mp4", index));
        
        match scene.content_type.as_str() {
            "text" => self.create_text_scene(scene, &scene_path),
            "image" => self.create_image_scene(scene, &scene_path),
            "video" => self.create_video_scene(scene, &scene_path),
            _ => Err(ProcessorError::Custom(format!("Unknown scene type: {}", scene.content_type)).into()),
        }
    }

    fn create_text_scene(&self, scene: &Scene, output_path: &Path) -> Result<PathBuf> {
        info!("Creating text scene: {}", scene.content);

        // Use FFmpeg to create a video with text overlay
        let mut cmd = Command::new("ffmpeg");
        cmd.args([
            "-f", "lavfi",
            "-i", &format!("color=c=black:size=1920x1080:duration={}", scene.duration),
            "-vf", &format!("drawtext=text='{}':fontcolor=white:fontsize=60:x=(w-text_w)/2:y=(h-text_h)/2", 
                          scene.content.replace("'", "\\'")),
            "-c:v", "libx264",
            "-pix_fmt", "yuv420p",
            "-y",
            output_path.to_str().unwrap(),
        ]);

        let output = cmd.output()?;
        if !output.status.success() {
            let error = String::from_utf8_lossy(&output.stderr);
            return Err(ProcessorError::Custom(format!("FFmpeg error: {}", error)).into());
        }

        Ok(output_path.to_path_buf())
    }

    fn create_image_scene(&self, scene: &Scene, output_path: &Path) -> Result<PathBuf> {
        info!("Creating image scene: {}", scene.content);

        // For now, we'll create a simple colored background
        // In a real implementation, you would download the image from scene.content URL
        let mut cmd = Command::new("ffmpeg");
        cmd.args([
            "-f", "lavfi",
            "-i", &format!("color=c=blue:size=1920x1080:duration={}", scene.duration),
            "-c:v", "libx264",
            "-pix_fmt", "yuv420p",
            "-y",
            output_path.to_str().unwrap(),
        ]);

        let output = cmd.output()?;
        if !output.status.success() {
            let error = String::from_utf8_lossy(&output.stderr);
            return Err(ProcessorError::Custom(format!("FFmpeg error: {}", error)).into());
        }

        Ok(output_path.to_path_buf())
    }

    fn create_video_scene(&self, scene: &Scene, output_path: &Path) -> Result<PathBuf> {
        info!("Creating video scene: {}", scene.content);

        // For now, create a simple test pattern
        // In a real implementation, you would download the video from scene.content URL
        let mut cmd = Command::new("ffmpeg");
        cmd.args([
            "-f", "lavfi",
            "-i", &format!("testsrc=duration={}:size=1920x1080:rate=30", scene.duration),
            "-c:v", "libx264",
            "-pix_fmt", "yuv420p",
            "-y",
            output_path.to_str().unwrap(),
        ]);

        let output = cmd.output()?;
        if !output.status.success() {
            let error = String::from_utf8_lossy(&output.stderr);
            return Err(ProcessorError::Custom(format!("FFmpeg error: {}", error)).into());
        }

        Ok(output_path.to_path_buf())
    }

    async fn combine_scenes(&self, scene_paths: &[PathBuf], output_path: &Path, script: &VideoScript) -> Result<()> {
        info!("Combining {} scenes into final video", scene_paths.len());

        // Create concat file for FFmpeg
        let concat_file = output_path.parent().unwrap().join("concat_list.txt");
        let mut concat_content = String::new();
        
        for path in scene_paths {
            concat_content.push_str(&format!("file '{}'\n", path.to_str().unwrap()));
        }
        
        fs::write(&concat_file, concat_content).await?;

        // Use FFmpeg to concatenate videos
        let mut cmd = Command::new("ffmpeg");
        cmd.args([
            "-f", "concat",
            "-safe", "0",
            "-i", concat_file.to_str().unwrap(),
            "-c", "copy",
            "-y",
            output_path.to_str().unwrap(),
        ]);

        // Add audio if available
        if let Some(audio_url) = &script.audio_url {
            info!("Adding audio track: {}", audio_url);
            // In a real implementation, you would download and mix the audio
        }

        let output = cmd.output()?;
        if !output.status.success() {
            let error = String::from_utf8_lossy(&output.stderr);
            return Err(ProcessorError::Custom(format!("FFmpeg concatenation error: {}", error)).into());
        }

        // Cleanup concat file
        if let Err(e) = fs::remove_file(&concat_file).await {
            warn!("Failed to cleanup concat file: {}", e);
        }

        Ok(())
    }

    pub async fn generate_thumbnail(&self, video_id: Uuid, title: &str) -> Result<String> {
        info!("Generating thumbnail for video: {}", video_id);

        let thumbnail_dir = self.temp_dir.join("thumbnails");
        fs::create_dir_all(&thumbnail_dir).await?;

        let thumbnail_path = thumbnail_dir.join(format!("{}.jpg", video_id));

        // Create thumbnail with title text
        let mut cmd = Command::new("ffmpeg");
        cmd.args([
            "-f", "lavfi",
            "-i", "color=c=0x1a1a1a:size=1280x720",
            "-vf", &format!(
                "drawtext=text='{}':fontcolor=white:fontsize=48:x=(w-text_w)/2:y=(h-text_h)/2:fontfile=/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
                title.replace("'", "\\'")
            ),
            "-frames:v", "1",
            "-y",
            thumbnail_path.to_str().unwrap(),
        ]);

        let output = cmd.output()?;
        if !output.status.success() {
            let error = String::from_utf8_lossy(&output.stderr);
            return Err(ProcessorError::Custom(format!("Thumbnail generation error: {}", error)).into());
        }

        // Upload to S3
        let s3_key = format!("thumbnails/{}.jpg", video_id);
        let thumbnail_url = self.s3_service.upload_file(&thumbnail_path, &s3_key).await?;

        // Update video with thumbnail URL
        self.update_thumbnail_url(video_id, &thumbnail_url).await?;

        // Cleanup
        if let Err(e) = fs::remove_file(&thumbnail_path).await {
            warn!("Failed to cleanup thumbnail file: {}", e);
        }

        Ok(thumbnail_url)
    }

    // Database helper methods
    async fn get_video(&self, video_id: Uuid) -> Result<Video> {
        let video = sqlx::query_as::<_, Video>(
            "SELECT * FROM videos WHERE id = $1"
        )
        .bind(video_id)
        .fetch_one(self.database.pool())
        .await?;

        Ok(video)
    }

    async fn update_video_status(&self, video_id: Uuid, status: &str) -> Result<()> {
        sqlx::query(
            "UPDATE videos SET status = $1, updated_at = NOW() WHERE id = $2"
        )
        .bind(status)
        .bind(video_id)
        .execute(self.database.pool())
        .await?;

        Ok(())
    }

    async fn update_video_url(&self, video_id: Uuid, video_url: &str) -> Result<()> {
        sqlx::query(
            "UPDATE videos SET video_url = $1, processing_completed_at = NOW(), updated_at = NOW() WHERE id = $2"
        )
        .bind(video_url)
        .bind(video_id)
        .execute(self.database.pool())
        .await?;

        Ok(())
    }

    async fn update_thumbnail_url(&self, video_id: Uuid, thumbnail_url: &str) -> Result<()> {
        sqlx::query(
            "UPDATE videos SET thumbnail_url = $1, updated_at = NOW() WHERE id = $2"
        )
        .bind(thumbnail_url)
        .bind(video_id)
        .execute(self.database.pool())
        .await?;

        Ok(())
    }

    async fn create_processing_job(&self, video_id: Uuid, job_type: &str) -> Result<Uuid> {
        let job_id = Uuid::new_v4();
        
        sqlx::query(
            r#"
            INSERT INTO content_pipeline (id, video_id, stage, status, started_at)
            VALUES ($1, $2, $3, 'processing', NOW())
            "#
        )
        .bind(job_id)
        .bind(video_id)
        .bind(job_type)
        .execute(self.database.pool())
        .await?;

        Ok(job_id)
    }

    async fn complete_processing_job(&self, job_id: Uuid, error_message: Option<&str>) -> Result<()> {
        let status = if error_message.is_some() { "failed" } else { "completed" };
        
        sqlx::query(
            r#"
            UPDATE content_pipeline 
            SET status = $1, error_message = $2, completed_at = NOW()
            WHERE id = $3
            "#
        )
        .bind(status)
        .bind(error_message)
        .bind(job_id)
        .execute(self.database.pool())
        .await?;

        Ok(())
    }
}