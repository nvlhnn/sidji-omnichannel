package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/sidji-omnichannel/internal/ai"
	"github.com/sidji-omnichannel/internal/config"
	"github.com/sidji-omnichannel/internal/models"
	"github.com/pgvector/pgvector-go"
)

type AIService struct {
	db       *sql.DB
	provider ai.Provider
}

func NewAIService(db *sql.DB, cfg *config.Config) (*AIService, error) {
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
		db:       db,
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
	config := &models.AIConfig{ChannelID: channelID}
	err := s.db.QueryRow(`
		SELECT id, is_enabled, mode, persona, handover_timeout_minutes, created_at, updated_at
		FROM ai_configs WHERE channel_id = $1
	`, channelID).Scan(
		&config.ID, &config.IsEnabled, &config.Mode, &config.Persona, &config.HandoverTimeoutMinutes, &config.CreatedAt, &config.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return &models.AIConfig{
			ChannelID:              channelID,
			IsEnabled:              false,
			Mode:                   models.AIModeManual,
			HandoverTimeoutMinutes: 15,
		}, nil
	}
	return config, err
}

// UpdateConfig updates AI configuration
func (s *AIService) UpdateConfig(cfg *models.AIConfig) error {
	_, err := s.db.Exec(`
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
	_, err = s.db.Exec(`
		INSERT INTO knowledge_base (channel_id, content, embedding)
		VALUES ($1, $2, $3)
	`, channelID, content, pgvector.NewVector(em))
	
	return err
}

// ListKnowledge retrieves all knowledge base items for a channel
func (s *AIService) ListKnowledge(channelID uuid.UUID) ([]*models.KnowledgeBaseItem, error) {
	rows, err := s.db.Query(`
		SELECT id, channel_id, content, created_at 
		FROM knowledge_base WHERE channel_id = $1
		ORDER BY created_at DESC
	`, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.KnowledgeBaseItem
	for rows.Next() {
		item := &models.KnowledgeBaseItem{}
		if err := rows.Scan(&item.ID, &item.ChannelID, &item.Content, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
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
	_, err = s.db.Exec(`
		UPDATE knowledge_base
		SET content = $1, embedding = $2, updated_at = NOW()
		WHERE id = $3
	`, content, pgvector.NewVector(em), id)
	
	return err
}

// DeleteKnowledge removes a text chunk from the knowledge base
func (s *AIService) DeleteKnowledge(id uuid.UUID) error {
	_, err := s.db.Exec("DELETE FROM knowledge_base WHERE id = $1", id)
	return err
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
	rows, err := s.db.Query(`
		SELECT content FROM knowledge_base
		WHERE channel_id = $1
		ORDER BY embedding <-> $2
		LIMIT $3
	`, channelID, pgvector.NewVector(queryVector), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			return nil, err
		}
		results = append(results, content)
	}
	return results, nil
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
	var creditsUsed, creditsLimit int
	err := s.db.QueryRow(`
		SELECT o.ai_credits_used, o.ai_credits_limit
		FROM organizations o
		JOIN channels c ON c.organization_id = o.id
		WHERE c.id = $1
	`, channelID).Scan(&creditsUsed, &creditsLimit)

	if err != nil {
		return fmt.Errorf("failed to check credits: %w", err)
	}

	// -1 means unlimited
	if creditsLimit >= 0 && creditsUsed >= creditsLimit {
		return fmt.Errorf("insufficient AI credits (%d/%d)", creditsUsed, creditsLimit)
	}
	return nil
}

// DeductCredit increments the used credits for the organization
func (s *AIService) DeductCredit(channelID uuid.UUID) error {
	_, err := s.db.Exec(`
		UPDATE organizations o
		SET ai_credits_used = ai_credits_used + 1, updated_at = NOW()
		FROM channels c
		WHERE c.organization_id = o.id AND c.id = $1
	`, channelID)
	if err != nil {
		return fmt.Errorf("failed to deduct credit: %w", err)
	}
	return nil
}
