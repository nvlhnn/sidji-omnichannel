-- Revert TikTok channel support
-- Migration: 012_add_tiktok_channel.down.sql

DROP INDEX IF EXISTS idx_contacts_tiktok_id;
ALTER TABLE contacts DROP COLUMN IF EXISTS tiktok_id;

DROP INDEX IF EXISTS idx_channels_tiktok_open_id;
ALTER TABLE channels DROP COLUMN IF EXISTS tiktok_open_id;

-- Revert channels type check constraint
ALTER TABLE channels DROP CONSTRAINT IF EXISTS channels_type_check;
ALTER TABLE channels ADD CONSTRAINT channels_type_check CHECK (type IN ('whatsapp', 'instagram', 'facebook'));
