use anyhow::Result;
use sqlx::{PgPool, Pool, Postgres};
use redis::aio::Connection;
use async_nats::{Client, Subscriber};
use std::sync::Arc;

#[derive(Debug, Clone)]
pub struct Database {
    pub pool: Arc<PgPool>,
}

impl Database {
    pub async fn connect(database_url: &str) -> Result<Self> {
        let pool = PgPool::connect(database_url).await?;
        
        Ok(Database {
            pool: Arc::new(pool),
        })
    }

    pub fn pool(&self) -> &PgPool {
        &self.pool
    }
}

#[derive(Debug)]
pub struct RedisClient {
    pub connection: Connection,
}

impl RedisClient {
    pub async fn connect(redis_url: &str) -> Result<Self> {
        let client = redis::Client::open(redis_url)?;
        let connection = client.get_async_connection().await?;
        
        Ok(RedisClient { connection })
    }
}

#[derive(Debug)]
pub struct NatsClient {
    pub client: Client,
}

impl NatsClient {
    pub async fn connect(nats_url: &str) -> Result<Self> {
        let client = async_nats::connect(nats_url).await?;
        
        Ok(NatsClient { client })
    }

    pub async fn subscribe(&self, subject: &str) -> Result<Subscriber> {
        Ok(self.client.subscribe(subject).await?)
    }

    pub async fn publish(&self, subject: &str, payload: Vec<u8>) -> Result<()> {
        self.client.publish(subject, payload.into()).await?;
        Ok(())
    }
}