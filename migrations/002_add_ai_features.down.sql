-- Drop AI features

DROP TRIGGER IF EXISTS update_knowledge_base_updated_at ON knowledge_base;
DROP TRIGGER IF EXISTS update_ai_configs_updated_at ON ai_configs;

ALTER TABLE conversations DROP COLUMN IF EXISTS ai_paused_until;
ALTER TABLE conversations DROP COLUMN IF EXISTS last_human_reply_at;

DROP TABLE IF EXISTS knowledge_base;
DROP TABLE IF EXISTS ai_configs;
DROP TYPE IF EXISTS ai_mode;

-- DROP EXTENSION IF EXISTS vector; -- Optional: Keep extension if used elsewhere, but safe to remove if this was the only user
