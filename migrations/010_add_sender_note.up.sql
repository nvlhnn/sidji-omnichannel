-- Add 'note' to sender_type check constraint
ALTER TABLE messages DROP CONSTRAINT IF EXISTS messages_sender_type_check;
ALTER TABLE messages ADD CONSTRAINT messages_sender_type_check CHECK (sender_type IN ('contact', 'agent', 'system', 'ai', 'note'));
