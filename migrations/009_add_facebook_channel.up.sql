-- Add Facebook channel support
-- Migration: 009_add_facebook_channel.up.sql

-- Add fb_page_id to channels
ALTER TABLE channels ADD COLUMN facebook_page_id VARCHAR(100);
CREATE INDEX idx_channels_facebook_page_id ON channels(facebook_page_id);

-- Update channels type check constraint
ALTER TABLE channels DROP CONSTRAINT IF EXISTS channels_type_check;
ALTER TABLE channels ADD CONSTRAINT channels_type_check CHECK (type IN ('whatsapp', 'instagram', 'facebook'));

-- Add facebook_id to contacts
ALTER TABLE contacts ADD COLUMN facebook_id VARCHAR(100);
CREATE INDEX idx_contacts_facebook_id ON contacts(facebook_id);
