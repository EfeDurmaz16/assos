# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ASSOS is a YouTube automation platform that uses AI agents (Manus/AutoGPT) to create, optimize, and manage YouTube content at scale. The system processes 1000+ videos per day using a multi-language microservices architecture optimized for performance.

## Architecture & Tech Stack

### Core Services by Language
- **Rust**: Video processing, analytics engine, file operations (performance-critical)
- **Go**: API Gateway, queue workers, microservice communication (high concurrency)  
- **Python 3.11+**: AI/ML operations, data science, research algorithms
- **Node.js**: Real-time features, webhooks, frontend SSR
- **Frontend**: Next.js 14 with TypeScript, Tailwind CSS, shadcn/ui

### Database Stack
- **PostgreSQL 15**: Primary database with JSONB for flexible schemas
- **Redis 7.0**: Caching and message queuing
- **Qdrant**: Vector database for AI embeddings
- **TimescaleDB**: Time-series analytics data
- **S3/MinIO**: Object storage for videos and assets

### Message Queue
- **NATS**: Go-based messaging system
- **Redpanda**: Kafka-compatible event streaming

## Development Commands

### Docker Environment
```bash
# Start all services
docker-compose up -d

# Start specific services
docker-compose up postgres redis minio

# View service logs
docker-compose logs -f [service-name]

# Rebuild services
docker-compose build --no-cache
```

### Service-Specific Commands

#### Rust Services (Video Processing)
```bash
# Build video processor
cd services/video-processor
cargo build --release

# Run with optimizations
RUSTFLAGS="-C target-cpu=native" cargo run --release

# Run tests
cargo test -- --nocapture
```

#### Go Services (API Gateway)
```bash
# Build gateway
cd services/api-gateway  
go build -o gateway ./cmd/gateway

# Run with race detection
go run -race ./cmd/gateway

# Run tests
go test -v ./...
```

#### Python Services (AI/ML)
```bash
# Install dependencies
pip install -r requirements.txt

# Run AI service
python -m services.ai_service

# Run specific agent
python -m agents.manus_agent --channel-id=<id>
```

#### Frontend
```bash
# Install dependencies
npm install

# Development server
npm run dev

# Build for production
npm run build

# Run type checking
npm run typecheck

# Run linting
npm run lint
```

## AI Agent System

### Manus Agent Integration
The platform uses Manus and AutoGPT agents for autonomous content creation:

- **ManusYouTubeAgent**: Primary orchestrator for channel management
- **ContentStrategistAgent**: Content planning and optimization  
- **TrendPredictorAgent**: Trend analysis and opportunity detection
- **PerformanceAnalystAgent**: Analytics and A/B testing

### Agent Communication
Agents communicate via NATS message bus using JSON-RPC over WebSocket protocol.

## Core Workflows

### Video Production Pipeline
1. **Research**: AI agents analyze trends and identify content opportunities
2. **Content Generation**: GPT-4 creates optimized scripts with retention hooks
3. **Voice Synthesis**: ElevenLabs/Azure Speech converts scripts to audio
4. **Video Assembly**: Rust service combines audio, visuals, and effects
5. **Publishing**: Automated upload to YouTube with optimal scheduling

### Performance Optimization
- Video processing: **5-8x faster** with Rust vs Python
- API Gateway: **Sub-20ms p99 latency** with Go vs 150ms with Python
- Memory usage: **80% reduction** with compiled languages

## Database Schema

### Key Tables
- `users`: User accounts and subscription tiers
- `channels`: YouTube channel configurations and settings
- `videos`: Video metadata and production status  
- `content_pipeline`: Processing stages and performance metrics

### Status Values
- Video statuses: `research`, `scripting`, `producing`, `published`
- Pipeline stages: Track processing through each step

## Testing

### Unit Tests
```bash
# Rust services
cargo test

# Go services  
go test ./...

# Python services
pytest tests/

# Frontend
npm test
```

### Integration Tests
```bash
# Full pipeline test
./scripts/test-pipeline.sh

# Load testing
./scripts/load-test.sh
```

## Monitoring & Debugging

### Service Health
```bash
# Check service status
docker-compose ps

# View resource usage
docker stats

# Monitor queue depth
redis-cli llen video:processing
```

### Logs
- Application logs: `./logs/`
- Error tracking: Sentry integration
- Metrics: Prometheus + Grafana dashboard
- Distributed tracing: OpenTelemetry

## Environment Configuration

### Required Environment Variables
```bash
# AI Services
OPENAI_API_KEY=
ANTHROPIC_API_KEY=  
ELEVENLABS_API_KEY=

# YouTube API
YOUTUBE_API_KEY=
YOUTUBE_CLIENT_ID=
YOUTUBE_CLIENT_SECRET=

# Database
DATABASE_URL=postgresql://user:pass@localhost/assos
REDIS_URL=redis://localhost:6379

# Storage
S3_BUCKET=assos-videos
S3_ACCESS_KEY=
S3_SECRET_KEY=
```

## Performance Targets

### Processing Capacity
- **Current**: 1000+ videos/day
- **Target**: 5000+ videos/day on same hardware
- **Video processing**: <30 minutes per video
- **API response**: <20ms p99 latency

### Quality Metrics
- **CTR Target**: 10%+ (thumbnail optimization)
- **AVD Target**: 50%+ (retention optimization)  
- **RPM Target**: $3+ (content selection)
- **Upload frequency**: 2-6 videos/channel/day

## Deployment

### Production Stack
- **Orchestration**: Kubernetes with Helm charts
- **Service mesh**: Istio for traffic management
- **Secrets**: Kubernetes secrets + external secret operator
- **Monitoring**: Prometheus, Grafana, ELK stack

### Scaling Configuration
- Video processing: Horizontal pod autoscaling based on queue depth
- API Gateway: Load balancer with health checks
- Database: Read replicas for analytics queries

## Security & Compliance

### API Security
- OAuth 2.0 + JWT authentication
- Rate limiting per user tier
- Request signing for sensitive operations

### Content Compliance  
- Automated content policy checking
- Copyright detection pipeline
- Community guidelines validation
- Age-appropriate content rating

## Brand Guidelines

### Color Palette
- YouTube Red: `#FF0000` (primary actions)
- AI Purple: `#8B5CF6` (AI-powered features)
- Automation Green: `#10B981` (success states)
- Off-white: `#FAFAFA` (backgrounds)

This platform represents a next-generation YouTube automation system that combines cutting-edge AI agents with performance-optimized microservices to deliver autonomous content creation at unprecedented scale.