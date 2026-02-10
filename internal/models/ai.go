package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

// AIMode defines the operating mode of the AI
type AIMode string

const (
	AIModeManual AIMode = "manual" // AI disabled (default)
	AIModeAuto   AIMode = "auto"   // AI replies to everything
	AIModeHybrid AIMode = "hybrid" // AI replies unless human took over
)

type AIProviderType string

const (
	AIProviderGemini AIProviderType = "gemini"
	AIProviderOpenAI AIProviderType = "openai"
)

const (
	ModelGeminiFlash     = "Gemini 2.5 Flash-Lite"
	ModelGeminiEmbedding = "gemini-embedding-001"
)

// AIConfig represents the AI settings for a channel
type AIConfig struct {
	ID                     uuid.UUID `json:"id"`
	ChannelID              uuid.UUID `json:"channel_id"`
	IsEnabled              bool      `json:"is_enabled"`
	Mode                   AIMode         `json:"mode"`
	Persona                string         `json:"persona"` // System prompt
	HandoverTimeoutMinutes int            `json:"handover_timeout_minutes"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
}

// KnowledgeBaseItem represents a chunk of information for RAG
type KnowledgeBaseItem struct {
	ID        uuid.UUID       `json:"id"`
	ChannelID uuid.UUID       `json:"channel_id"`
	Content   string          `json:"content"`
	Embedding pgvector.Vector `json:"-"` // Vector embedding, not exposed in JSON
	Metadata  map[string]any  `json:"metadata"`
	CreatedAt time.Time       `json:"created_at"`
}

// AIRequest represents a request to the AI service
type AIRequest struct {
	ChannelID uuid.UUID
	Query     string
	History   []Message // Context from recent chat
}
