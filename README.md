# ASSOS - YouTube Automation Platform

AI-powered YouTube content creation and automation platform using advanced AI agents (Manus/AutoGPT) to create, optimize, and manage YouTube content at scale.

## 🚀 Features

- **AI Agent Orchestration**: Manus and AutoGPT agents for autonomous content creation
- **Multi-Language Architecture**: Rust for performance, Go for concurrency, Python for AI/ML
- **Scalable Processing**: Handle 1000+ videos per day
- **Intelligent Research**: Automated trend analysis and content optimization
- **Real-time Analytics**: Performance tracking and optimization
- **Multi-Channel Management**: Manage multiple YouTube channels simultaneously

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   API Gateway   │    │   AI Service    │
│   (Next.js)     │◄──►│   (Go/Fiber)    │◄──►│   (Python)      │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ Video Processor │    │   PostgreSQL    │    │   Redis Cache   │
│   (Rust)        │◄──►│   Database      │◄──►│   & Queue       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🛠️ Quick Start

### Prerequisites
- Docker & Docker Compose
- Git

### Setup
```bash
# Clone the repository
git clone <repository-url>
cd assos

# Copy environment variables
cp .env.example .env
# Edit .env with your API keys

# Start the development environment
make dev

# Build and start all services
make build
make up
```

### Accessing Services
- Frontend: http://localhost:3000
- API Gateway: http://localhost:8080
- MinIO (S3): http://localhost:9001
- Grafana: http://localhost:3001
- Prometheus: http://localhost:9090

## 🔧 Development

### Service Commands
```bash
# View logs
make logs

# Run tests
make test

# Reset database
make db-reset

# Service-specific logs
make api-logs
make ai-logs
make video-logs
```

### Service Development

#### API Gateway (Go)
```bash
cd services/api-gateway
go run ./cmd/gateway
```

#### AI Service (Python)
```bash
cd services/ai-service
pip install -r requirements.txt
python -m uvicorn main:app --reload
```

#### Video Processor (Rust)
```bash
cd services/video-processor
cargo run
```

#### Frontend (Next.js)
```bash
cd frontend
npm install
npm run dev
```

## 📊 Performance Targets

- **Processing**: 1000+ videos/day
- **API Latency**: <20ms p99
- **Video Processing**: <30 minutes per video
- **Quality Metrics**: 10%+ CTR, 50%+ AVD

## 🔐 Environment Variables

Required API keys:
- `OPENAI_API_KEY`: OpenAI GPT-4 API
- `ANTHROPIC_API_KEY`: Claude API
- `ELEVENLABS_API_KEY`: Voice synthesis
- `YOUTUBE_API_KEY`: YouTube Data API

## 🏛️ Database Schema

Core tables:
- `users`: User accounts and subscriptions
- `channels`: YouTube channel configurations
- `videos`: Video metadata and status
- `content_pipeline`: Processing stages
- `ai_agents`: AI agent configurations
- `agent_tasks`: Agent task management

## 🤖 AI Agent System

### Agent Types
- **Manus Orchestrator**: Primary strategy and coordination
- **Content Strategist**: Script writing and optimization
- **Trend Predictor**: Market analysis and opportunities
- **Performance Analyst**: Analytics and A/B testing
- **Research Agent**: Comprehensive content research

### Agent Communication
- **Message Bus**: NATS for inter-agent communication
- **Protocol**: JSON-RPC over WebSocket
- **Memory**: Vector database for persistent learning

## 📈 Monitoring

- **Metrics**: Prometheus + Grafana
- **Logs**: Docker Compose logs
- **Health Checks**: Built-in service monitoring
- **Alerts**: Configurable thresholds

## 🔒 Security

- **Authentication**: JWT + OAuth 2.0
- **API Security**: Rate limiting and request signing
- **Data Protection**: Encryption at rest and in transit
- **Content Compliance**: Automated policy checking

## 📚 Documentation

- [Architecture Guide](docs/architecture.md)
- [API Documentation](docs/api.md)
- [AI Agent Guide](docs/agents.md)
- [Deployment Guide](docs/deployment.md)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `make test`
5. Submit a pull request

## 📄 License

MIT License - see LICENSE file for details.

## 🆘 Support

- Documentation: Check the docs/ directory
- Issues: Create GitHub issues
- Discussions: GitHub Discussions tab

---

**ASSOS** - "Create. Automate. Dominate."