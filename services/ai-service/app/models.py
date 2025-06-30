from datetime import datetime
from typing import Optional, Dict, Any, List
from uuid import UUID
from pydantic import BaseModel


class VideoProcessingRequest(BaseModel):
    video_id: str
    user_id: str
    action: str


class AgentTaskRequest(BaseModel):
    task_id: str
    agent_id: str
    video_id: str
    task_type: str
    priority: int
    input_data: Dict[str, Any]


class ScriptGenerationRequest(BaseModel):
    topic: str
    niche: str
    target_duration: int  # in minutes
    style: str = "educational"
    tone: str = "professional"
    target_audience: str = "general"


class ResearchRequest(BaseModel):
    topic: str
    niche: str
    depth: str = "comprehensive"  # basic, comprehensive, deep
    sources: List[str] = ["youtube", "google_trends", "reddit"]


class VideoScript(BaseModel):
    title: str
    description: str
    scenes: List[Dict[str, Any]]
    total_duration: float
    keywords: List[str]
    tags: List[str]
    hooks: List[str]
    call_to_actions: List[str]


class ResearchData(BaseModel):
    topic: str
    findings: List[Dict[str, Any]]
    trends: List[Dict[str, Any]]
    competitors: List[Dict[str, Any]]
    keywords: List[Dict[str, Any]]
    content_gaps: List[Dict[str, Any]]
    confidence_score: float


class TrendAnalysis(BaseModel):
    keyword: str
    trend_score: float
    search_volume: int
    competition_level: str
    related_keywords: List[str]
    seasonal_patterns: Dict[str, Any]


class ContentIdea(BaseModel):
    title: str
    description: str
    estimated_views: int
    difficulty_score: float
    viral_potential: float
    monetization_potential: float
    target_keywords: List[str]
    content_type: str
    estimated_duration: int


class AgentResponse(BaseModel):
    agent_id: str
    task_id: str
    status: str
    result: Optional[Dict[str, Any]] = None
    error: Optional[str] = None
    execution_time: Optional[float] = None
    confidence_score: Optional[float] = None