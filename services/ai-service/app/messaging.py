import asyncio
import json
import logging
from typing import Dict, Any
import nats
from nats.aio.client import Client as NATS

from .config import settings
from .agents import ManusAgent, ContentStrategistAgent, TrendPredictorAgent, ResearchAgent
from .models import VideoProcessingRequest, AgentTaskRequest

logger = logging.getLogger(__name__)


class MessageProcessor:
    """Processes messages from NATS message queue"""
    
    def __init__(self):
        self.nats_client: NATS = None
        self.agents = {
            "manus": ManusAgent(),
            "content_strategist": ContentStrategistAgent(),
            "trend_predictor": TrendPredictorAgent(),
            "research_agent": ResearchAgent(),
        }
    
    async def start(self):
        """Start the message processor"""
        try:
            # Connect to NATS
            self.nats_client = await nats.connect(settings.NATS_URL)
            logger.info("Connected to NATS")
            
            # Subscribe to relevant subjects
            await self._setup_subscriptions()
            
            logger.info("Message processor started")
            
        except Exception as e:
            logger.error(f"Failed to start message processor: {e}")
            raise
    
    async def _setup_subscriptions(self):
        """Set up NATS subscriptions"""
        
        # Video processing messages
        await self.nats_client.subscribe(
            "video.process", 
            cb=self._handle_video_processing
        )
        
        # AI task messages
        await self.nats_client.subscribe(
            "ai.task",
            cb=self._handle_ai_task
        )
        
        # Research requests
        await self.nats_client.subscribe(
            "ai.research",
            cb=self._handle_research_request
        )
        
        # Content generation requests
        await self.nats_client.subscribe(
            "ai.content.generate",
            cb=self._handle_content_generation
        )
        
        logger.info("NATS subscriptions set up")
    
    async def _handle_video_processing(self, msg):
        """Handle video processing messages"""
        try:
            data = json.loads(msg.data.decode())
            request = VideoProcessingRequest(**data)
            
            logger.info(f"Processing video: {request.video_id}")
            
            if request.action == "start_processing":
                # Use Manus agent to orchestrate the process
                response = await self.agents["manus"].process_task(
                    task_id=f"video_{request.video_id}",
                    task_type="orchestrate_video_creation",
                    input_data={
                        "video_id": request.video_id,
                        "user_id": request.user_id,
                        "channel_config": {}  # This would come from the database
                    }
                )
                
                # Send response back
                await self._send_response("ai.video.response", response.model_dump())
            
        except Exception as e:
            logger.error(f"Error handling video processing: {e}")
    
    async def _handle_ai_task(self, msg):
        """Handle AI task messages"""
        try:
            data = json.loads(msg.data.decode())
            request = AgentTaskRequest(**data)
            
            logger.info(f"Processing AI task: {request.task_id} for agent: {request.agent_id}")
            
            # Route to appropriate agent
            agent = self._get_agent_by_type(request.agent_id)
            if agent:
                response = await agent.process_task(
                    task_id=request.task_id,
                    task_type=request.task_type,
                    input_data=request.input_data
                )
                
                # Send response back
                await self._send_response("ai.task.response", response.model_dump())
            else:
                logger.error(f"Unknown agent: {request.agent_id}")
                
        except Exception as e:
            logger.error(f"Error handling AI task: {e}")
    
    async def _handle_research_request(self, msg):
        """Handle research requests"""
        try:
            data = json.loads(msg.data.decode())
            
            # Use research agent
            response = await self.agents["research_agent"].process_task(
                task_id=data.get("task_id", "research_task"),
                task_type="comprehensive_research",
                input_data=data
            )
            
            await self._send_response("ai.research.response", response.model_dump())
            
        except Exception as e:
            logger.error(f"Error handling research request: {e}")
    
    async def _handle_content_generation(self, msg):
        """Handle content generation requests"""
        try:
            data = json.loads(msg.data.decode())
            
            # Use content strategist agent
            response = await self.agents["content_strategist"].process_task(
                task_id=data.get("task_id", "content_task"),
                task_type="script_generation",
                input_data=data
            )
            
            await self._send_response("ai.content.response", response.model_dump())
            
        except Exception as e:
            logger.error(f"Error handling content generation: {e}")
    
    def _get_agent_by_type(self, agent_type: str):
        """Get agent by type or ID"""
        # Map agent types to instances
        agent_mapping = {
            "manus": self.agents["manus"],
            "content_strategist": self.agents["content_strategist"],
            "trend_predictor": self.agents["trend_predictor"],
            "research_agent": self.agents["research_agent"],
        }
        
        return agent_mapping.get(agent_type)
    
    async def _send_response(self, subject: str, data: Dict[str, Any]):
        """Send response message"""
        try:
            message = json.dumps(data).encode()
            await self.nats_client.publish(subject, message)
            logger.debug(f"Sent response to {subject}")
        except Exception as e:
            logger.error(f"Failed to send response: {e}")
    
    async def stop(self):
        """Stop the message processor"""
        if self.nats_client:
            await self.nats_client.close()
            logger.info("Message processor stopped")