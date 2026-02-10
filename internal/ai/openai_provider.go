package ai

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sidji-omnichannel/internal/models"
	"github.com/sashabaranov/go-openai"
)

// OpenAIProvider implements the Provider interface using OpenAI
type OpenAIProvider struct {
	client *openai.Client
}

func NewOpenAIProvider(apiKey string) (*OpenAIProvider, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}
	client := openai.NewClient(apiKey)
	return &OpenAIProvider{client: client}, nil
}

func (p *OpenAIProvider) Name() string {
	return string(models.AIProviderOpenAI)
}

func (p *OpenAIProvider) EmbedText(ctx context.Context, text string) ([]float32, error) {
	resp, err := p.client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: []string{text},
		Model: openai.AdaEmbeddingV2, // Standard embedding model
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned from OpenAI")
	}

	return resp.Data[0].Embedding, nil
}

func (p *OpenAIProvider) GenerateReply(ctx context.Context, config *models.AIConfig, userQuery string, contextDocs []string, history []models.Message) (string, error) {
	contextStr := strings.Join(contextDocs, "\n\n")

	// Format history
	var messages []openai.ChatCompletionMessage
	
	// System Instructions
	systemMsg := fmt.Sprintf(`You are a polite, patient, and professional business assistant. Use simple, easy-to-understand language. Provide direct, relevant, and concise responses that are clear and to the point.
Additional rules: %s

CONTEXT FROM KNOWLEDGE BASE:
%s`, config.Persona, contextStr)

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: systemMsg,
	})

	// History
	for _, m := range history {
		role := openai.ChatMessageRoleUser
		if m.SenderType == models.SenderAgent || m.SenderType == models.SenderAI {
			role = openai.ChatMessageRoleAssistant
		}
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    role,
			Content: m.Content,
		})
	}

	// Current Query
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userQuery,
	})

	// Request with retries
	var resp openai.ChatCompletionResponse
	var lastErr error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		resp, lastErr = p.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model:       openai.GPT4oMini, // High performance, low cost (similar to Gemini Flash)
			Messages:    messages,
			MaxTokens:   400, // Approx 800 characters
			Temperature: 0.7,
		})

		if lastErr == nil {
			break
		}

		// Handle rate limits (429)
		if strings.Contains(lastErr.Error(), "429") {
			log.Printf("OpenAI Rate Limit Hit (429), retry %d/%d in 2s...", i+1, maxRetries)
			time.Sleep(time.Duration(2*(i+1)) * time.Second)
			continue
		}
		return "", lastErr
	}

	if lastErr != nil {
		return "", lastErr
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("empty response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}
