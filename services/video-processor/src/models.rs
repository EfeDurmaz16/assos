use chrono::{DateTime, Utc};
use serde::{Deserialize, Serialize};
use sqlx::FromRow;
use uuid::Uuid;

#[derive(Debug, Clone, Serialize, Deserialize, FromRow)]
pub struct Video {
    pub id: Uuid,
    pub channel_id: Uuid,
    pub title: Option<String>,
    pub description: Option<String>,
    pub status: String,
    pub youtube_video_id: Option<String>,
    pub thumbnail_url: Option<String>,
    pub video_url: Option<String>,
    pub script: Option<serde_json::Value>,
    pub metadata: Option<serde_json::Value>,
    pub performance_data: Option<serde_json::Value>,
    pub ai_analysis: Option<serde_json::Value>,
    pub processing_started_at: Option<DateTime<Utc>>,
    pub processing_completed_at: Option<DateTime<Utc>>,
    pub published_at: Option<DateTime<Utc>>,
    pub created_at: DateTime<Utc>,
    pub updated_at: DateTime<Utc>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct VideoProcessingMessage {
    pub video_id: String,
    pub user_id: String,
    pub action: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ThumbnailMessage {
    pub video_id: String,
    pub title: String,
    pub style: Option<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct UploadMessage {
    pub video_id: String,
    pub file_path: String,
    pub metadata: VideoMetadata,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct VideoMetadata {
    pub title: String,
    pub description: String,
    pub tags: Vec<String>,
    pub category_id: String,
    pub privacy_status: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Scene {
    pub duration: f64,
    pub content_type: String, // "text", "image", "video"
    pub content: String,
    pub transition: Option<String>,
    pub effects: Vec<String>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct VideoScript {
    pub scenes: Vec<Scene>,
    pub total_duration: f64,
    pub audio_url: Option<String>,
    pub voice_settings: Option<serde_json::Value>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ProcessingJob {
    pub id: Uuid,
    pub video_id: Uuid,
    pub job_type: String,
    pub status: String,
    pub input_data: serde_json::Value,
    pub output_data: Option<serde_json::Value>,
    pub error_message: Option<String>,
    pub started_at: DateTime<Utc>,
    pub completed_at: Option<DateTime<Utc>>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct RenderSettings {
    pub resolution: String,     // "1920x1080", "1280x720"
    pub fps: u32,              // 30, 60
    pub bitrate: String,       // "5000k", "8000k"
    pub format: String,        // "mp4", "webm"
    pub quality: String,       // "high", "medium", "low"
}