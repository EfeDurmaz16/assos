use crate::config::Config;
use anyhow::Result;
use aws_config::meta::region::RegionProviderChain;
use aws_config::BehaviorVersion;
use aws_sdk_s3::{Client, Config as S3Config};
use aws_sdk_s3::config::{Credentials, Region};
use std::path::Path;
use tokio::fs::File;
use tokio::io::AsyncReadExt;
use tracing::{info, error};

#[derive(Debug, Clone)]
pub struct S3Service {
    client: Client,
    bucket: String,
}

impl S3Service {
    pub async fn new(config: &Config) -> Result<Self> {
        let credentials = Credentials::new(
            &config.s3_access_key,
            &config.s3_secret_key,
            None,
            None,
            "assos",
        );

        let s3_config = S3Config::builder()
            .region(Region::new(config.s3_region.clone()))
            .credentials_provider(credentials)
            .endpoint_url(&config.s3_endpoint)
            .force_path_style(true)
            .behavior_version(BehaviorVersion::latest())
            .build();

        let client = Client::from_conf(s3_config);

        // Ensure bucket exists
        let s3_service = S3Service {
            client,
            bucket: config.s3_bucket.clone(),
        };

        s3_service.ensure_bucket_exists().await?;

        Ok(s3_service)
    }

    async fn ensure_bucket_exists(&self) -> Result<()> {
        match self.client.head_bucket()
            .bucket(&self.bucket)
            .send()
            .await
        {
            Ok(_) => {
                info!("S3 bucket '{}' exists", self.bucket);
            }
            Err(_) => {
                info!("Creating S3 bucket '{}'", self.bucket);
                self.client.create_bucket()
                    .bucket(&self.bucket)
                    .send()
                    .await?;
            }
        }

        Ok(())
    }

    pub async fn upload_file(&self, file_path: &Path, s3_key: &str) -> Result<String> {
        let mut file = File::open(file_path).await?;
        let mut contents = Vec::new();
        file.read_to_end(&mut contents).await?;

        let content_type = match file_path.extension().and_then(|ext| ext.to_str()) {
            Some("mp4") => "video/mp4",
            Some("jpg") | Some("jpeg") => "image/jpeg",
            Some("png") => "image/png",
            Some("webm") => "video/webm",
            _ => "application/octet-stream",
        };

        info!("Uploading file to S3: {} -> {}", file_path.display(), s3_key);

        self.client.put_object()
            .bucket(&self.bucket)
            .key(s3_key)
            .body(contents.into())
            .content_type(content_type)
            .send()
            .await?;

        let url = format!("http://localhost:9000/{}/{}", self.bucket, s3_key);
        info!("File uploaded successfully: {}", url);

        Ok(url)
    }

    pub async fn upload_bytes(&self, data: Vec<u8>, s3_key: &str, content_type: &str) -> Result<String> {
        info!("Uploading bytes to S3: {} ({} bytes)", s3_key, data.len());

        self.client.put_object()
            .bucket(&self.bucket)
            .key(s3_key)
            .body(data.into())
            .content_type(content_type)
            .send()
            .await?;

        let url = format!("http://localhost:9000/{}/{}", self.bucket, s3_key);
        info!("Bytes uploaded successfully: {}", url);

        Ok(url)
    }

    pub async fn delete_file(&self, s3_key: &str) -> Result<()> {
        info!("Deleting file from S3: {}", s3_key);

        self.client.delete_object()
            .bucket(&self.bucket)
            .key(s3_key)
            .send()
            .await?;

        Ok(())
    }

    pub async fn file_exists(&self, s3_key: &str) -> bool {
        self.client.head_object()
            .bucket(&self.bucket)
            .key(s3_key)
            .send()
            .await
            .is_ok()
    }

    pub fn get_public_url(&self, s3_key: &str) -> String {
        format!("http://localhost:9000/{}/{}", self.bucket, s3_key)
    }
}