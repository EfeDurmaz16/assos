from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession, async_sessionmaker
from sqlalchemy.orm import declarative_base
import redis.asyncio as redis
from qdrant_client.async_qdrant_client import AsyncQdrantClient
from qdrant_client.models import Distance, VectorParams
import logging

from .config import settings

logger = logging.getLogger(__name__)

# SQLAlchemy setup
engine = create_async_engine(
    settings.DATABASE_URL,
    echo=settings.DEBUG,
    pool_pre_ping=True,
    pool_recycle=300
)

async_session_maker = async_sessionmaker(
    engine,
    class_=AsyncSession,
    expire_on_commit=False
)

Base = declarative_base()

# Redis client
redis_client = None

# Qdrant client
qdrant_client = None


async def get_db():
    """Get database session"""
    async with async_session_maker() as session:
        try:
            yield session
        finally:
            await session.close()


async def get_redis():
    """Get Redis client"""
    global redis_client
    if redis_client is None:
        redis_client = redis.from_url(settings.REDIS_URL)
    return redis_client


async def get_qdrant():
    """Get Qdrant client"""
    global qdrant_client
    if qdrant_client is None:
        qdrant_client = AsyncQdrantClient(url=settings.QDRANT_URL)
    return qdrant_client


async def init_db():
    """Initialize database and create tables if needed"""
    try:
        # Test database connection
        async with engine.begin() as conn:
            logger.info("Database connection established")
        
        # Initialize Redis
        redis_conn = await get_redis()
        await redis_conn.ping()
        logger.info("Redis connection established")
        
        # Initialize Qdrant and create collections
        qdrant = await get_qdrant()
        
        # Create collections for different types of embeddings
        collections = [
            {
                "name": "video_scripts",
                "size": 1536,  # OpenAI ada-002 embedding size
                "description": "Embeddings for video scripts"
            },
            {
                "name": "research_data", 
                "size": 1536,
                "description": "Embeddings for research data"
            },
            {
                "name": "performance_data",
                "size": 1536,
                "description": "Embeddings for performance analytics"
            }
        ]
        
        for collection in collections:
            try:
                await qdrant.create_collection(
                    collection_name=collection["name"],
                    vectors_config=VectorParams(
                        size=collection["size"],
                        distance=Distance.COSINE
                    )
                )
                logger.info(f"Created Qdrant collection: {collection['name']}")
            except Exception as e:
                if "already exists" in str(e):
                    logger.info(f"Qdrant collection already exists: {collection['name']}")
                else:
                    logger.error(f"Error creating collection {collection['name']}: {e}")
        
        logger.info("Database initialization completed")
        
    except Exception as e:
        logger.error(f"Database initialization failed: {e}")
        raise