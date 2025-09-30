import asyncio
import logging
import time
from abc import ABC, abstractmethod
from typing import Dict, Any, Optional
from uuid import UUID, uuid4
import json

from ..models import AgentResponse
from ..services.llm_service import LLMService
from ..database import get_redis, get_qdrant
from qdrant_client.models import PointStruct

logger = logging.getLogger(__name__)


class BaseAgent(ABC):
    """Base class for all AI agents"""
    
    def __init__(self, agent_id: str, name: str, agent_type: str):
        self.agent_id = agent_id
        self.name = name
        self.agent_type = agent_type
        self.llm_service = LLMService()
        self.memory = {}
        self.performance_metrics = {
            "tasks_completed": 0,
            "success_rate": 0.0,
            "avg_execution_time": 0.0,
            "confidence_scores": []
        }
    
    @abstractmethod
    async def execute_task(self, task_type: str, input_data: Dict[str, Any]) -> AgentResponse:
        """Execute a specific task"""
        pass
    
    @abstractmethod
    async def get_capabilities(self) -> Dict[str, Any]:
        """Get agent capabilities and supported task types"""
        pass
    
    async def process_task(self, task_id: str, task_type: str, input_data: Dict[str, Any]) -> AgentResponse:
        """Process a task with timing and error handling"""
        start_time = time.time()
        
        try:
            logger.info(f"Agent {self.name} starting task {task_id} of type {task_type}")
            
            # Load context from memory if needed
            input_data["context"] = await self._load_context(task_id, input_data)
            
            # Execute the task
            response = await self.execute_task(task_type, input_data)
            
            # Calculate execution time
            execution_time = time.time() - start_time
            response.execution_time = execution_time
            
            # Update performance metrics
            await self._update_metrics(execution_time, response.confidence_score or 0.8)
            
            # Store results in memory
            await self._store_results(task_id, input_data, response)
            
            logger.info(f"Agent {self.name} completed task {task_id} in {execution_time:.2f}s")
            return response
            
        except Exception as e:
            execution_time = time.time() - start_time
            error_response = AgentResponse(
                agent_id=self.agent_id,
                task_id=task_id,
                status="failed",
                error=str(e),
                execution_time=execution_time
            )
            
            logger.error(f"Agent {self.name} failed task {task_id}: {e}", exc_info=True)
            return error_response
    
    def _get_collection_name(self, task_type: str) -> Optional[str]:
        if "research" in task_type or "ideation" in task_type:
            return "research_data"
        if "script" in task_type or "content" in task_type:
            return "video_scripts"
        if "performance" in task_type or "analytics" in task_type:
            return "performance_data"
        return None

    async def _load_context(self, task_id: str, input_data: Dict[str, Any]) -> Dict[str, Any]:
        """Load relevant context from memory and vector database"""
        context = {}
        try:
            # Load from Redis cache
            redis = await get_redis()
            cached_context = await redis.get(f"agent_context:{self.agent_id}:{task_id}")
            
            if cached_context:
                context['redis_cache'] = json.loads(cached_context)
            
            # Load relevant embeddings from Qdrant
            collection_name = self._get_collection_name(input_data.get("task_type", ""))
            if collection_name:
                qdrant = await get_qdrant()
                query_text = json.dumps(input_data)
                query_vector = await self.llm_service.generate_embeddings(query_text)

                search_result = await qdrant.search(
                    collection_name=collection_name,
                    query_vector=query_vector,
                    limit=5
                )
                context['vector_search_results'] = [hit.payload for hit in search_result]
                
        except Exception as e:
            logger.warning(f"Failed to load context for task {task_id}: {e}")

        return context

    async def _store_results(self, task_id: str, input_data: Dict[str, Any], response: AgentResponse):
        """Store task results in memory and vector database"""
        try:
            # Store in Redis cache
            redis = await get_redis()
            await redis.setex(
                f"agent_result:{self.agent_id}:{task_id}",
                3600,  # 1 hour TTL
                response.model_dump_json()
            )
            
            # Store embeddings in Qdrant if we have text content
            collection_name = self._get_collection_name(input_data.get("task_type", ""))
            if collection_name and response.result and isinstance(response.result, dict):
                qdrant = await get_qdrant()

                # Create embedding from result
                text_to_embed = json.dumps(response.result)
                vector = await self.llm_service.generate_embeddings(text_to_embed)

                point = PointStruct(
                    id=str(uuid4()),
                    vector=vector,
                    payload={
                        "task_id": task_id,
                        "agent_id": self.agent_id,
                        "agent_name": self.name,
                        "task_type": input_data.get("task_type"),
                        "input_data": input_data,
                        "result": response.result,
                        "timestamp": time.time()
                    }
                )

                await qdrant.upsert(collection_name=collection_name, points=[point], wait=True)
                
        except Exception as e:
            logger.warning(f"Failed to store results for task {task_id}: {e}")
    
    async def _update_metrics(self, execution_time: float, confidence_score: float):
        """Update agent performance metrics"""
        self.performance_metrics["tasks_completed"] += 1
        
        # Update average execution time
        current_avg = self.performance_metrics["avg_execution_time"]
        task_count = self.performance_metrics["tasks_completed"]
        self.performance_metrics["avg_execution_time"] = (
            (current_avg * (task_count - 1) + execution_time) / task_count
        )
        
        # Track confidence scores
        self.performance_metrics["confidence_scores"].append(confidence_score)
        
        # Keep only last 100 confidence scores
        if len(self.performance_metrics["confidence_scores"]) > 100:
            self.performance_metrics["confidence_scores"] = (
                self.performance_metrics["confidence_scores"][-100:]
            )
        
        # Calculate success rate (confidence > 0.7 is considered success)
        successful_tasks = sum(
            1 for score in self.performance_metrics["confidence_scores"] 
            if score > 0.7
        )
        self.performance_metrics["success_rate"] = (
            successful_tasks / len(self.performance_metrics["confidence_scores"])
        )
    
    async def get_performance_metrics(self) -> Dict[str, Any]:
        """Get current performance metrics"""
        return {
            **self.performance_metrics,
            "agent_id": self.agent_id,
            "name": self.name,
            "type": self.agent_type
        }
    
    def _create_system_prompt(self, task_type: str) -> str:
        """Create system prompt for the agent"""
        base_prompt = f"""
        You are {self.name}, a specialized AI agent for YouTube content automation.
        Your role is to {self.agent_type} with high accuracy and efficiency.
        
        Current task type: {task_type}
        
        Guidelines:
        - Provide detailed, actionable responses
        - Include confidence scores for your recommendations
        - Consider YouTube algorithm preferences
        - Focus on audience engagement and retention
        - Ensure content is monetization-friendly
        """
        return base_prompt