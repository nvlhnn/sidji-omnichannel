-- Add AI features

-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- AI Configuration for Channels
CREATE TYPE ai_mode AS ENUM ('manual', 'auto', 'hybrid');

CREATE TABLE ai_configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    channel_id UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    is_enabled BOOLEAN DEFAULT false,
    mode ai_mode DEFAULT 'manual',
    persona TEXT DEFAULT 'You are a helpful assistant for this business. Answer politely and concisely.',
    handover_timeout_minutes INT DEFAULT 15,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(channel_id)
);

-- AI Knowledge Base (RAG)
CREATE TABLE knowledge_base (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    channel_id UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    embedding vector(768), -- Dimension for Gemini Embedding model (text-embedding-004) or similar
    metadata JSONB DEFAULT '{}', -- Store source info (filename, url, etc)
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index for vector search (IVFFlat is good for starter, HNSW is better for scale but slower build)
-- We use validation to ensure enough data exists before creating index if needed, 
-- but Postgres allows creating empty index.
CREATE INDEX idx_knowledge_base_embedding ON knowledge_base USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);

-- Track last human interaction for hybrid mode
ALTER TABLE conversations ADD COLUMN last_human_reply_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE conversations ADD COLUMN ai_paused_until TIMESTAMP WITH TIME ZONE; -- Explicit pause until time

-- Update trigger for ai_configs
CREATE TRIGGER update_ai_configs_updated_at BEFORE UPDATE ON ai_configs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Update trigger for knowledge_base
CREATE TRIGGER update_knowledge_base_updated_at BEFORE UPDATE ON knowledge_base FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
