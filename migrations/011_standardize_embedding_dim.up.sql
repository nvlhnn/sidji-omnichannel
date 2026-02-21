-- Standardize embedding dimension to 1536 (OpenAI ada-002 compatible)
-- Both OpenAI and Gemini providers now output 1536 dimensions
-- so they can be used interchangeably.

-- Clear existing embeddings since dimension is changing
DELETE FROM knowledge_base;

ALTER TABLE knowledge_base ALTER COLUMN embedding TYPE vector(1536);
