#!/bin/bash

# ASSOS Development Startup Script
echo "🚀 Starting ASSOS YouTube Automation Platform..."

# Check if Docker is running
if ! docker info >/dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker first."
    exit 1
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp .env.example .env
    echo "⚠️  Please edit .env file with your API keys before running again."
    exit 1
fi

# Build and start services
echo "🔨 Building services..."
make build

echo "🗄️  Starting database services..."
make dev

echo "⏳ Waiting for databases to be ready..."
sleep 10

echo "🚀 Starting all services..."
make up

echo "✅ ASSOS Platform is starting up!"
echo ""
echo "📊 Service URLs:"
echo "  • Frontend:     http://localhost:3000"
echo "  • API Gateway:  http://localhost:8080"
echo "  • MinIO (S3):   http://localhost:9001"
echo "  • Grafana:      http://localhost:3001"
echo "  • Prometheus:   http://localhost:9090"
echo ""
echo "📋 Default Credentials:"
echo "  • MinIO:        assos_admin / assos_password"
echo "  • Grafana:      admin / admin"
echo "  • Demo User:    demo@assos.ai / demo123"
echo ""
echo "🔗 View logs with: make logs"
echo "🛑 Stop services with: make down"