package postgres

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"github.com/sidji-omnichannel/internal/domain/ports/repository"
	"github.com/sidji-omnichannel/internal/models"
)

type aiRepository struct {
	db *sql.DB
}

func NewAIRepository(db *sql.DB) repository.AIRepository {
	return &aiRepository{db: db}
}

func (r *aiRepository) GetConfig(channelID uuid.UUID) (*models.AIConfig, error) {
	config := &models.AIConfig{ChannelID: channelID}
	err := r.db.QueryRow(`
		SELECT id, is_enabled, mode, persona, handover_timeout_minutes, created_at, updated_at
		FROM ai_configs WHERE channel_id = $1
	`, channelID).Scan(&config.ID, &config.IsEnabled, &config.Mode, &config.Persona, &config.HandoverTimeoutMinutes, &config.CreatedAt, &config.UpdatedAt)
	if err == sql.ErrNoRows { return nil, repository.ErrNotFound }
	return config, err
}

func (r *aiRepository) UpdateConfig(cfg *models.AIConfig) error {
	_, err := r.db.Exec(`
		INSERT INTO ai_configs (channel_id, is_enabled, mode, persona, handover_timeout_minutes, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (channel_id) DO UPDATE SET
			is_enabled = EXCLUDED.is_enabled,
			mode = EXCLUDED.mode,
			persona = EXCLUDED.persona,
			handover_timeout_minutes = EXCLUDED.handover_timeout_minutes,
			updated_at = NOW()
	`, cfg.ChannelID, cfg.IsEnabled, cfg.Mode, cfg.Persona, cfg.HandoverTimeoutMinutes)
	return err
}

func (r *aiRepository) SaveKnowledge(channelID uuid.UUID, content string, embedding []float32) error {
	_, err := r.db.Exec(`
		INSERT INTO knowledge_base (channel_id, content, embedding)
		VALUES ($1, $2, $3)
	`, channelID, content, pgvector.NewVector(embedding))
	return err
}

func (r *aiRepository) ListKnowledge(channelID uuid.UUID) ([]*models.KnowledgeBaseItem, error) {
	rows, err := r.db.Query(`SELECT id, channel_id, content, created_at FROM knowledge_base WHERE channel_id = $1 ORDER BY created_at DESC`, channelID)
	if err != nil { return nil, err }
	defer rows.Close()
	var items []*models.KnowledgeBaseItem
	for rows.Next() {
		item := &models.KnowledgeBaseItem{}
		if err := rows.Scan(&item.ID, &item.ChannelID, &item.Content, &item.CreatedAt); err != nil { return nil, err }
		items = append(items, item)
	}
	return items, nil
}

func (r *aiRepository) UpdateKnowledge(id uuid.UUID, content string, embedding []float32) error {
	_, err := r.db.Exec(`UPDATE knowledge_base SET content = $1, embedding = $2, updated_at = NOW() WHERE id = $3`, content, pgvector.NewVector(embedding), id)
	return err
}

func (r *aiRepository) DeleteKnowledge(id uuid.UUID) error {
	_, err := r.db.Exec("DELETE FROM knowledge_base WHERE id = $1", id)
	return err
}

func (r *aiRepository) SearchKnowledge(channelID uuid.UUID, queryVector []float32, limit int) ([]string, error) {
	rows, err := r.db.Query(`SELECT content FROM knowledge_base WHERE channel_id = $1 ORDER BY embedding <-> $2 LIMIT $3`, channelID, pgvector.NewVector(queryVector), limit)
	if err != nil { return nil, err }
	defer rows.Close()
	var results []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil { return nil, err }
		results = append(results, content)
	}
	return results, nil
}

func (r *aiRepository) GetCredits(channelID uuid.UUID) (used, limit int, err error) {
	err = r.db.QueryRow(`
		SELECT o.ai_credits_used, o.ai_credits_limit
		FROM organizations o
		JOIN channels c ON c.organization_id = o.id
		WHERE c.id = $1
	`, channelID).Scan(&used, &limit)
	return used, limit, err
}

func (r *aiRepository) DeductCredit(channelID uuid.UUID) error {
	_, err := r.db.Exec(`
		UPDATE organizations o SET ai_credits_used = ai_credits_used + 1, updated_at = NOW()
		FROM channels c WHERE c.organization_id = o.id AND c.id = $1
	`, channelID)
	return err
}
