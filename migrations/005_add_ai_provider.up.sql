-- Add provider column to ai_configs
ALTER TABLE ai_configs ADD COLUMN provider TEXT DEFAULT 'gemini';
