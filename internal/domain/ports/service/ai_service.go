package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/models"
)

// AIService defines the inbound port for AI operations
type AIService interface {
	GetConfig(channelID uuid.UUID) (*models.AIConfig, error)
	UpdateConfig(cfg *models.AIConfig) error
	DetermineAction(channelID uuid.UUID, lastHumanReplyAt *time.Time) (bool, bool, *models.AIConfig, error)
	AddKnowledge(channelID uuid.UUID, content string) error
	ListKnowledge(channelID uuid.UUID) ([]*models.KnowledgeBaseItem, error)
	UpdateKnowledge(id uuid.UUID, content string) error
	DeleteKnowledge(id uuid.UUID) error
	EmbedText(channelID uuid.UUID, text string) ([]float32, error)
	SearchKnowledge(channelID uuid.UUID, queryVector []float32, limit int) ([]string, error)
	GenerateReply(ctx context.Context, config *models.AIConfig, userQuery string, contextDocs []string, history []models.Message) (string, error)
	CheckCredits(channelID uuid.UUID) error
	DeductCredit(channelID uuid.UUID) error
}
