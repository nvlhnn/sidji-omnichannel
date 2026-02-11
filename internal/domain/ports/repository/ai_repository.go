package repository

import (
	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

// AIRepository defines the outbound port for AI-related data operations
type AIRepository interface {
	GetConfig(channelID uuid.UUID) (*models.AIConfig, error)
	UpdateConfig(cfg *models.AIConfig) error
	SaveKnowledge(channelID uuid.UUID, content string, embedding []float32) error
	ListKnowledge(channelID uuid.UUID) ([]*models.KnowledgeBaseItem, error)
	UpdateKnowledge(id uuid.UUID, content string, embedding []float32) error
	DeleteKnowledge(id uuid.UUID) error
	SearchKnowledge(channelID uuid.UUID, queryVector []float32, limit int) ([]string, error)
	GetCredits(channelID uuid.UUID) (used, limit int, err error)
	DeductCredit(channelID uuid.UUID) error
}
