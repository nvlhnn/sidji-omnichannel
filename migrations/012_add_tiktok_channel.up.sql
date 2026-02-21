-- Add TikTok channel support
-- Migration: 012_add_tiktok_channel.up.sql

-- Add tiktok_open_id to channels
ALTER TABLE channels ADD COLUMN IF NOT EXISTS tiktok_open_id VARCHAR(100);
CREATE INDEX IF NOT EXISTS idx_channels_tiktok_open_id ON channels(tiktok_open_id);

-- Update channels type check constraint
ALTER TABLE channels DROP CONSTRAINT IF EXISTS channels_type_check;
ALTER TABLE channels ADD CONSTRAINT channels_type_check CHECK (type IN ('whatsapp', 'instagram', 'facebook', 'tiktok'));

-- Add tiktok_id to contacts
ALTER TABLE contacts ADD COLUMN IF NOT EXISTS tiktok_id VARCHAR(100);
CREATE INDEX IF NOT EXISTS idx_contacts_tiktok_id ON contacts(tiktok_id);
