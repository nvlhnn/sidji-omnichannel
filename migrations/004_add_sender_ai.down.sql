-- Revert sender_type check constraint
-- WARNING: This will delete messages with sender_type = 'ai' to make it compatible
DELETE FROM messages WHERE sender_type = 'ai';
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_sender_type_check;
ALTER TABLE messages ADD CONSTRAINT messages_sender_type_check CHECK (sender_type IN ('contact', 'agent', 'system'));
