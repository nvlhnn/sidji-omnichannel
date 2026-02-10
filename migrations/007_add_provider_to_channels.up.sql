-- Add provider column to channels table to support abstraction (Twilio/Meta)
ALTER TABLE channels ADD COLUMN provider VARCHAR(50) NOT NULL DEFAULT 'meta';
