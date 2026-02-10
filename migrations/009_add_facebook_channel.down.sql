-- Remove Facebook channel support
-- Migration: 009_add_facebook_channel.down.sql

ALTER TABLE contacts DROP COLUMN IF EXISTS facebook_id;
DROP INDEX IF EXISTS idx_contacts_facebook_id;

ALTER TABLE channels DROP COLUMN IF EXISTS facebook_page_id;
DROP INDEX IF EXISTS idx_channels_facebook_page_id;

-- Restore channels type check constraint
ALTER TABLE channels DROP CONSTRAINT IF EXISTS channels_type_check;
ALTER TABLE channels ADD CONSTRAINT channels_type_check CHECK (type IN ('whatsapp', 'instagram'));
