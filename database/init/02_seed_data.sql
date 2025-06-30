-- ASSOS Seed Data
-- Insert default AI agents
INSERT INTO ai_agents (name, type, configuration) VALUES 
('Manus Orchestrator', 'manus', '{"capabilities": ["strategy", "coordination", "optimization"], "max_concurrent_tasks": 10}'),
('Content Strategist', 'content_strategist', '{"specialties": ["script_writing", "seo_optimization", "engagement"], "models": ["gpt-4"]}'),
('Trend Predictor', 'trend_predictor', '{"sources": ["youtube", "google_trends", "reddit", "twitter"], "analysis_window": "7d"}'),
('Performance Analyst', 'performance_analyst', '{"metrics": ["ctr", "avd", "engagement", "revenue"], "reporting_frequency": "daily"}'),
('Research Agent', 'research_agent', '{"search_depth": "comprehensive", "sources": ["academic", "news", "social"], "fact_check": true}');

-- Insert demo user (password: demo123)
INSERT INTO users (email, password_hash, subscription_tier) VALUES 
('demo@assos.ai', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/lewdBzpvWw0UtGkoa', 'pro');

-- Get the demo user ID for foreign key relationships
DO $$
DECLARE
    demo_user_id UUID;
BEGIN
    SELECT id INTO demo_user_id FROM users WHERE email = 'demo@assos.ai';
    
    -- Insert demo channel
    INSERT INTO channels (user_id, name, niche, description, settings, posting_schedule) VALUES 
    (demo_user_id, 'Tech Insights AI', 'technology', 'AI-powered tech content creation', 
     '{"voice": "professional", "style": "educational", "duration": "8-12min"}',
     '{"frequency": "daily", "optimal_times": ["10:00", "15:00", "19:00"], "timezone": "UTC"}');
END $$;