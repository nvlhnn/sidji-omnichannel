INSERT INTO schema_migrations (version) VALUES
('001_initial_schema'),
('002_add_ai_features'),
('003_fix_embedding_dim'),
('004_add_sender_ai'),
('005_add_ai_provider'),
('006_add_subscription_fields'),
('007_add_provider_to_channels'),
('008_add_message_usage_fields'),
('009_add_facebook_channel'),
('010_add_sender_note'),
('011_standardize_embedding_dim')
ON CONFLICT DO NOTHING;
