package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/ai"
	"github.com/sidji-omnichannel/internal/config"
	"github.com/sidji-omnichannel/internal/domain/ports/repository"
	"github.com/sidji-omnichannel/internal/models"
)

type AIService struct {
	repo     repository.AIRepository
	provider ai.Provider
}

func NewAIService(repo repository.AIRepository, cfg *config.Config) (*AIService, error) {
	var provider ai.Provider
	var err error

	pType := models.AIProviderType(cfg.AI.Provider)
	if pType == "" {
		pType = models.AIProviderGemini
	}

	switch pType {
	case models.AIProviderGemini:
		if cfg.AI.GeminiAPIKey != "" {
			provider, err = ai.NewGeminiProvider(cfg.AI.GeminiAPIKey)
		} else {
			log.Println("⚠️ GEMINI_API_KEY missing.")
		}
	case models.AIProviderOpenAI:
		if cfg.AI.OpenAIAPIKey != "" {
			provider, err = ai.NewOpenAIProvider(cfg.AI.OpenAIAPIKey)
		} else {
			log.Println("⚠️ OPENAI_API_KEY missing.")
		}
	default:
		log.Printf("⚠️ Unknown AI provider: %s", pType)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to init AI provider %s: %w", pType, err)
	}

	return &AIService{
		repo:     repo,
		provider: provider,
	}, nil
}

func (s *AIService) getActiveProvider() (ai.Provider, error) {
	if s.provider == nil {
		return nil, fmt.Errorf("AI provider not initialized. Check your global config/env.")
	}
	return s.provider, nil
}

// GetConfig returns the AI configuration for a channel
func (s *AIService) GetConfig(channelID uuid.UUID) (*models.AIConfig, error) {
	cfg, err := s.repo.GetConfig(channelID)
	if err == repository.ErrNotFound {
		return &models.AIConfig{
			ChannelID:              channelID,
			IsEnabled:              false,
			Mode:                   models.AIModeManual,
			HandoverTimeoutMinutes: 15,
		}, nil
	}
	return cfg, err
}

// UpdateConfig updates AI configuration
func (s *AIService) UpdateConfig(cfg *models.AIConfig) error {
	return s.repo.UpdateConfig(cfg)
}

// DetermineAction checks what AI should do: auto-reply, suggest draft, or nothing
func (s *AIService) DetermineAction(channelID uuid.UUID, lastHumanReplyAt *time.Time) (bool, bool, *models.AIConfig, error) {
	cfg, err := s.GetConfig(channelID)
	if err != nil {
		return false, false, nil, err
	}

	if !cfg.IsEnabled {
		return false, false, cfg, nil
	}

	if cfg.Mode == models.AIModeManual {
		return false, true, cfg, nil
	}

	if cfg.Mode == models.AIModeHybrid && lastHumanReplyAt != nil {
		now := time.Now()
		handoverDuration := time.Duration(cfg.HandoverTimeoutMinutes) * time.Minute
		elapsed := now.Sub(*lastHumanReplyAt)
		
		log.Printf("[AI Debug] mode=hybrid now=%v last_reply=%v elapsed=%v timeout=%v", 
			now.Format(time.RFC3339), lastHumanReplyAt.Format(time.RFC3339), elapsed, handoverDuration)
		
		if elapsed >= 0 && elapsed < handoverDuration {
			log.Printf("[AI Debug] Staying silent: elapsed %v < timeout %v", elapsed, handoverDuration)
			return false, false, cfg, nil
		}
		
		if elapsed < 0 {
			log.Printf("[AI Debug] WARNING: Negative elapsed time (Clock drift?). AI taking over.")
		} else {
			log.Printf("[AI Debug] Handover timeout expired (%v). AI taking over.", elapsed)
		}
	} else if cfg.Mode == models.AIModeHybrid {
		log.Printf("[AI Debug] mode=hybrid last_human_reply=nil. AI taking over.")
	}

	return true, false, cfg, nil
}


// AddKnowledge adds text to the knowledge base
func (s *AIService) AddKnowledge(channelID uuid.UUID, content string) error {
	provider, err := s.getActiveProvider()
	if err != nil {
		return err
	}

	// 2. Get embedding
	em, err := provider.EmbedText(context.Background(), content)
	if err != nil {
		return err
	}

	// 3. Save to DB
	return s.repo.SaveKnowledge(channelID, content, em)
}

// ListKnowledge retrieves all knowledge base items for a channel
func (s *AIService) ListKnowledge(channelID uuid.UUID) ([]*models.KnowledgeBaseItem, error) {
	return s.repo.ListKnowledge(channelID)
}

// UpdateKnowledge updates a knowledge base item's content and its embedding
func (s *AIService) UpdateKnowledge(id uuid.UUID, content string) error {
	provider, err := s.getActiveProvider()
	if err != nil {
		return err
	}

	// 1. Get new embedding
	em, err := provider.EmbedText(context.Background(), content)
	if err != nil {
		return err
	}

	// 2. Update DB
	return s.repo.UpdateKnowledge(id, content, em)
}

// DeleteKnowledge removes a text chunk from the knowledge base
func (s *AIService) DeleteKnowledge(id uuid.UUID) error {
	return s.repo.DeleteKnowledge(id)
}

// EmbedText is now a helper that uses the active provider
func (s *AIService) EmbedText(channelID uuid.UUID, text string) ([]float32, error) {
	provider, err := s.getActiveProvider()
	if err != nil {
		return nil, err
	}
	return provider.EmbedText(context.Background(), text)
}

// SearchKnowledge finds relevant context
func (s *AIService) SearchKnowledge(channelID uuid.UUID, queryVector []float32, limit int) ([]string, error) {
	return s.repo.SearchKnowledge(channelID, queryVector, limit)
}

// GenerateReply generates a response using the selected provider
func (s *AIService) GenerateReply(ctx context.Context, config *models.AIConfig, userQuery string, contextDocs []string, history []models.Message) (string, error) {
	provider, err := s.getActiveProvider()
	if err != nil {
		return "", err
	}
	return provider.GenerateReply(ctx, config, userQuery, contextDocs, history)
}

// CheckCredits verifies if the organization has sufficient credits
func (s *AIService) CheckCredits(channelID uuid.UUID) error {
	used, limit, err := s.repo.GetCredits(channelID)
	if err != nil {
		return fmt.Errorf("failed to check credits: %w", err)
	}

	// -1 means unlimited
	if limit >= 0 && used >= limit {
		return fmt.Errorf("insufficient AI credits (%d/%d)", used, limit)
	}
	return nil
}

// DeductCredit increments the used credits for the organization
func (s *AIService) DeductCredit(channelID uuid.UUID) error {
	if err := s.repo.DeductCredit(channelID); err != nil {
		return fmt.Errorf("failed to deduct credit: %w", err)
	}
	return nil
}
