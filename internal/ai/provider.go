package ai

import (
	"context"

	"github.com/sidji-omnichannel/internal/models"
)

// Provider defines the interface for different AI backends
type Provider interface {
	EmbedText(ctx context.Context, text string) ([]float32, error)
	GenerateReply(ctx context.Context, config *models.AIConfig, userQuery string, contextDocs []string, history []models.Message) (string, error)
	Name() string
}
