from fastapi import APIRouter, HTTPException, Depends
from typing import Dict, Any, List
import logging

from .agents import ManusAgent, ContentStrategistAgent, TrendPredictorAgent, ResearchAgent
from .models import ScriptGenerationRequest, ResearchRequest, AgentResponse

logger = logging.getLogger(__name__)

router = APIRouter()

# Initialize agents
manus_agent = ManusAgent()
content_agent = ContentStrategistAgent()
trend_agent = TrendPredictorAgent()
research_agent = ResearchAgent()


@router.get("/agents")
async def get_agents():
    """Get list of available AI agents"""
    agents = [
        await manus_agent.get_capabilities(),
        await content_agent.get_capabilities(),
        await trend_agent.get_capabilities(),
        await research_agent.get_capabilities(),
    ]
    
    return {
        "agents": agents,
        "count": len(agents)
    }


@router.post("/agents/manus/orchestrate")
async def orchestrate_video_creation(request: Dict[str, Any]):
    """Orchestrate video creation using Manus agent"""
    try:
        response = await manus_agent.process_task(
            task_id=request.get("task_id", "orchestrate_task"),
            task_type="orchestrate_video_creation",
            input_data=request
        )
        
        return response.model_dump()
        
    except Exception as e:
        logger.error(f"Orchestration failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/agents/manus/strategy")
async def create_strategy(request: Dict[str, Any]):
    """Create strategic plan using Manus agent"""
    try:
        response = await manus_agent.process_task(
            task_id=request.get("task_id", "strategy_task"),
            task_type="strategic_planning",
            input_data=request
        )
        
        return response.model_dump()
        
    except Exception as e:
        logger.error(f"Strategy creation failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/content/script")
async def generate_script(request: ScriptGenerationRequest):
    """Generate video script"""
    try:
        response = await content_agent.process_task(
            task_id="script_generation",
            task_type="script_generation",
            input_data=request.model_dump()
        )
        
        return response.model_dump()
        
    except Exception as e:
        logger.error(f"Script generation failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/content/ideas")
async def generate_content_ideas(request: Dict[str, Any]):
    """Generate content ideas"""
    try:
        response = await manus_agent.process_task(
            task_id="content_ideation",
            task_type="content_ideation",
            input_data=request
        )
        
        return response.model_dump()
        
    except Exception as e:
        logger.error(f"Content ideation failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/research/comprehensive")
async def conduct_research(request: ResearchRequest):
    """Conduct comprehensive research"""
    try:
        response = await research_agent.process_task(
            task_id="research_task",
            task_type="comprehensive_research",
            input_data=request.model_dump()
        )
        
        return response.model_dump()
        
    except Exception as e:
        logger.error(f"Research failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/trends/analyze")
async def analyze_trends(request: Dict[str, Any]):
    """Analyze trends for given topic/niche"""
    try:
        response = await trend_agent.process_task(
            task_id="trend_analysis",
            task_type="trend_analysis",
            input_data=request
        )
        
        return response.model_dump()
        
    except Exception as e:
        logger.error(f"Trend analysis failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.get("/agents/{agent_id}/performance")
async def get_agent_performance(agent_id: str):
    """Get performance metrics for specific agent"""
    agent_map = {
        "manus": manus_agent,
        "content_strategist": content_agent,
        "trend_predictor": trend_agent,
        "research_agent": research_agent,
    }
    
    agent = agent_map.get(agent_id)
    if not agent:
        raise HTTPException(status_code=404, detail="Agent not found")
    
    try:
        metrics = await agent.get_performance_metrics()
        return metrics
        
    except Exception as e:
        logger.error(f"Failed to get performance metrics: {e}")
        raise HTTPException(status_code=500, detail=str(e))


@router.post("/optimize/performance")
async def optimize_performance(request: Dict[str, Any]):
    """Optimize content performance using Manus agent"""
    try:
        response = await manus_agent.process_task(
            task_id="performance_optimization",
            task_type="performance_optimization",
            input_data=request
        )
        
        return response.model_dump()
        
    except Exception as e:
        logger.error(f"Performance optimization failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))