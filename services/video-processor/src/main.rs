mod config;
mod database;
mod error;
mod handlers;
mod models;
mod services;

use std::sync::Arc;
use anyhow::Result;
use tracing::{info, error};
use tracing_subscriber;

use crate::config::Config;
use crate::database::{Database, RedisClient, NatsClient};
use crate::handlers::MessageHandler;
use crate::services::{VideoProcessor, S3Service};

#[tokio::main]
async fn main() -> Result<()> {
    // Initialize tracing
    tracing_subscriber::init();

    // Load configuration
    let config = Config::load()?;
    info!("Starting ASSOS Video Processor");

    // Initialize database connections
    let database = Database::connect(&config.database_url).await?;
    let redis = RedisClient::connect(&config.redis_url).await?;
    let nats = NatsClient::connect(&config.nats_url).await?;

    // Initialize S3 service
    let s3_service = S3Service::new(&config).await?;

    // Initialize video processor
    let video_processor = Arc::new(VideoProcessor::new(
        database.clone(),
        s3_service,
    ));

    // Initialize message handler
    let message_handler = MessageHandler::new(
        database,
        redis,
        video_processor,
    );

    // Subscribe to NATS subjects
    info!("Subscribing to NATS subjects...");
    
    let video_sub = nats.subscribe("video.process").await?;
    let thumbnail_sub = nats.subscribe("video.thumbnail").await?;
    let upload_sub = nats.subscribe("video.upload").await?;

    // Start message processing loops
    let handler = Arc::new(message_handler);
    
    let h1 = handler.clone();
    let video_task = tokio::spawn(async move {
        info!("Starting video processing handler");
        while let Some(message) = video_sub.next().await {
            if let Err(e) = h1.handle_video_processing(message).await {
                error!("Error processing video message: {}", e);
            }
        }
    });

    let h2 = handler.clone();
    let thumbnail_task = tokio::spawn(async move {
        info!("Starting thumbnail processing handler");
        while let Some(message) = thumbnail_sub.next().await {
            if let Err(e) = h2.handle_thumbnail_generation(message).await {
                error!("Error processing thumbnail message: {}", e);
            }
        }
    });

    let h3 = handler.clone();
    let upload_task = tokio::spawn(async move {
        info!("Starting upload handler");
        while let Some(message) = upload_sub.next().await {
            if let Err(e) = h3.handle_video_upload(message).await {
                error!("Error processing upload message: {}", e);
            }
        }
    });

    info!("Video Processor is running. Press Ctrl+C to stop.");

    // Wait for all tasks or shutdown signal
    tokio::select! {
        _ = video_task => {
            error!("Video processing task ended");
        }
        _ = thumbnail_task => {
            error!("Thumbnail processing task ended");
        }
        _ = upload_task => {
            error!("Upload processing task ended");
        }
        _ = tokio::signal::ctrl_c() => {
            info!("Shutdown signal received");
        }
    }

    info!("Shutting down Video Processor");
    Ok(())
}