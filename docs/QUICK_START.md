# ASSOS Quick Start Guide

## Overview

ASSOS is an AI-powered YouTube automation platform that creates, optimizes, and manages YouTube content at scale using advanced AI agents.

## Prerequisites

- Docker & Docker Compose
- Git
- At least 8GB RAM
- 20GB free disk space

## Quick Setup

1. **Clone and Setup**
   ```bash
   cd assos
   cp .env.example .env
   ```

2. **Configure API Keys**
   Edit `.env` file with your API keys:
   ```bash
   OPENAI_API_KEY=your_openai_api_key_here
   ANTHROPIC_API_KEY=your_anthropic_api_key_here
   ELEVENLABS_API_KEY=your_elevenlabs_api_key_here
   YOUTUBE_API_KEY=your_youtube_api_key_here
   ```

3. **Start Development Environment**
   ```bash
   ./scripts/start-dev.sh
   ```

## Access Points

- **Frontend Dashboard**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **MinIO S3 Console**: http://localhost:9001
- **Grafana Monitoring**: http://localhost:3001
- **Prometheus Metrics**: http://localhost:9090

## Default Credentials

- **Demo User**: demo@assos.ai / demo123
- **MinIO**: assos_admin / assos_password
- **Grafana**: admin / admin

## Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Next.js       │    │   Go Fiber      │    │   Python        │
│   Frontend      │◄──►│   API Gateway   │◄──►│   AI Service    │
│   (React)       │    │   (High Perf)   │    │   (GPT/Claude)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Rust Video    │    │   PostgreSQL    │    │   Redis Cache   │
│   Processor     │◄──►│   Database      │◄──►│   & Queue       │
│   (Performance) │    │   (JSONB)       │    │   (NATS)        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Key Features

### AI Agents
- **Manus Orchestrator**: Primary strategic AI agent
- **Content Strategist**: Script writing and optimization
- **Trend Predictor**: Market analysis and forecasting
- **Research Agent**: Comprehensive content research
- **Performance Analyst**: Analytics and A/B testing

### Video Pipeline
1. **Research**: AI analyzes trends and opportunities
2. **Script Generation**: GPT-4 creates optimized scripts
3. **Voice Synthesis**: ElevenLabs converts text to speech
4. **Video Assembly**: Rust service combines all elements
5. **Publishing**: Automated upload to YouTube

### Performance
- **Processing**: 1000+ videos/day
- **API Latency**: <20ms p99
- **Video Processing**: <30 minutes per video
- **Languages**: Rust (speed), Go (concurrency), Python (AI), React (UI)

## Development Commands

```bash
# Start development databases only
make dev

# Build all services
make build

# Start all services
make up

# View logs
make logs

# Stop all services
make down

# Clean everything
make clean

# Run tests
make test

# Service-specific logs
make api-logs      # API Gateway logs
make ai-logs       # AI Service logs
make video-logs    # Video Processor logs
```

## Service Development

### API Gateway (Go)
```bash
cd services/api-gateway
go run ./cmd/gateway
```

### AI Service (Python)
```bash
cd services/ai-service
pip install -r requirements.txt
python main.py
```

### Video Processor (Rust)
```bash
cd services/video-processor
cargo run
```

### Frontend (Next.js)
```bash
cd frontend
npm install
npm run dev
```

## Database Schema

### Core Tables
- `users`: User accounts and subscriptions
- `channels`: YouTube channel configurations
- `videos`: Video metadata and processing status
- `content_pipeline`: Processing stage tracking
- `ai_agents`: AI agent configurations
- `agent_tasks`: Task management and results

## AI Integration

### Manus Agent Examples
```python
# Orchestrate video creation
await manus_agent.process_task(
    task_type="orchestrate_video_creation",
    input_data={
        "video_id": "uuid",
        "channel_config": {...},
        "target_niche": "technology"
    }
)

# Strategic planning
await manus_agent.process_task(
    task_type="strategic_planning",
    input_data={
        "channel_data": {...},
        "goals": {"subscribers": 100000}
    }
)
```

## Environment Variables

### Required
```bash
# Database
DATABASE_URL=postgres://assos_user:assos_password@localhost:5432/assos
REDIS_URL=redis://localhost:6379

# AI APIs
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=sk-ant-...
ELEVENLABS_API_KEY=...

# YouTube
YOUTUBE_API_KEY=...
```

### Optional
```bash
# Security
JWT_SECRET=your-secret-key

# Storage
S3_ENDPOINT=http://localhost:9000
S3_BUCKET=assos-videos

# Environment
NODE_ENV=development
LOG_LEVEL=info
```

## Monitoring

### Health Checks
- API Gateway: http://localhost:8080/health
- AI Service: http://localhost:8000/health

### Metrics
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3001

### Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api-gateway
```

## Troubleshooting

### Common Issues

1. **Port Conflicts**
   ```bash
   # Check port usage
   netstat -tulpn | grep :8080
   
   # Kill process
   sudo kill -9 $(lsof -t -i:8080)
   ```

2. **Database Connection Issues**
   ```bash
   # Check database status
   docker-compose ps postgres
   
   # Reset database
   make db-reset
   ```

3. **Out of Memory**
   ```bash
   # Check Docker memory
   docker system df
   
   # Clean up
   docker system prune -f
   ```

4. **API Key Issues**
   ```bash
   # Test OpenAI API
   curl https://api.openai.com/v1/models \
     -H "Authorization: Bearer $OPENAI_API_KEY"
   ```

## Next Steps

1. **Configure API Keys**: Add your AI service API keys
2. **Create First Video**: Use the dashboard to create your first automated video
3. **Set Up Channels**: Connect your YouTube channels
4. **Configure AI Agents**: Customize agent behavior for your niche
5. **Monitor Performance**: Use Grafana dashboards to track metrics

## Support

- **Documentation**: Check `/docs` directory
- **Issues**: Create GitHub issues for bugs
- **Discussions**: Use GitHub Discussions for questions

---

**ASSOS** - "Create. Automate. Dominate."