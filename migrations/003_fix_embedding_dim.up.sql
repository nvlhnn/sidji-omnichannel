-- Drop IVFFlat/HNSW index as they support max 2000 dimensions (usually). 
-- pgvector 0.7.0 increased HNSW limit to 4000, but we are likely on older version or default limit.
-- For 3072 dims, we might not be able to index it efficiently on standard pgconfig without tuning or upgrading.
-- So we just drop the index for now and rely on sequential scan (slow but works for small datasets).
DROP INDEX IF EXISTS idx_knowledge_base_embedding;

-- Update embedding dimension to 3072
ALTER TABLE knowledge_base ALTER COLUMN embedding TYPE vector(3072);
