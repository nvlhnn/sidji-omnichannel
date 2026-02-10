-- Add message usage fields to organizations table
ALTER TABLE organizations
ADD COLUMN message_usage_limit INTEGER NOT NULL DEFAULT 1000,
ADD COLUMN message_usage_used INTEGER NOT NULL DEFAULT 0;
