use std::fmt;

#[derive(Debug)]
pub enum ProcessorError {
    Database(sqlx::Error),
    Redis(redis::RedisError),
    Nats(async_nats::Error),
    S3(aws_sdk_s3::Error),
    Ffmpeg(ffmpeg_next::Error),
    Io(std::io::Error),
    Serde(serde_json::Error),
    Custom(String),
}

impl fmt::Display for ProcessorError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        match self {
            ProcessorError::Database(e) => write!(f, "Database error: {}", e),
            ProcessorError::Redis(e) => write!(f, "Redis error: {}", e),
            ProcessorError::Nats(e) => write!(f, "NATS error: {}", e),
            ProcessorError::S3(e) => write!(f, "S3 error: {}", e),
            ProcessorError::Ffmpeg(e) => write!(f, "FFmpeg error: {}", e),
            ProcessorError::Io(e) => write!(f, "IO error: {}", e),
            ProcessorError::Serde(e) => write!(f, "Serialization error: {}", e),
            ProcessorError::Custom(e) => write!(f, "Error: {}", e),
        }
    }
}

impl std::error::Error for ProcessorError {}

impl From<sqlx::Error> for ProcessorError {
    fn from(e: sqlx::Error) -> Self {
        ProcessorError::Database(e)
    }
}

impl From<redis::RedisError> for ProcessorError {
    fn from(e: redis::RedisError) -> Self {
        ProcessorError::Redis(e)
    }
}

impl From<async_nats::Error> for ProcessorError {
    fn from(e: async_nats::Error) -> Self {
        ProcessorError::Nats(e)
    }
}

impl From<std::io::Error> for ProcessorError {
    fn from(e: std::io::Error) -> Self {
        ProcessorError::Io(e)
    }
}

impl From<serde_json::Error> for ProcessorError {
    fn from(e: serde_json::Error) -> Self {
        ProcessorError::Serde(e)
    }
}