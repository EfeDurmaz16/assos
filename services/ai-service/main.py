import asyncio
import logging
from contextlib import asynccontextmanager

import uvicorn
from fastapi import FastAPI

from app.config import settings
from app.database import init_db
from app.messaging import MessageProcessor
from app.api import router

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)

logger = logging.getLogger(__name__)


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan manager"""
    logger.info("Starting ASSOS AI Service")
    
    # Initialize database
    await init_db()
    
    # Start message processor
    message_processor = MessageProcessor()
    task = asyncio.create_task(message_processor.start())
    
    yield
    
    # Cleanup
    logger.info("Shutting down ASSOS AI Service")
    task.cancel()
    try:
        await task
    except asyncio.CancelledError:
        pass


# Create FastAPI app
app = FastAPI(
    title="ASSOS AI Service",
    description="AI-powered content generation and orchestration service",
    version="1.0.0",
    lifespan=lifespan
)

# Include API router
app.include_router(router, prefix="/api/v1")


@app.get("/health")
async def health_check():
    """Health check endpoint"""
    return {
        "status": "healthy",
        "service": "assos-ai-service",
        "version": "1.0.0"
    }


if __name__ == "__main__":
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=8000,
        reload=settings.DEBUG,
        log_level="info"
    )