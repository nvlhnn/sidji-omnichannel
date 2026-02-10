package ai

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/sidji-omnichannel/internal/models"
	"google.golang.org/api/option"
)

// GeminiProvider implements the Provider interface using Google Gemini
type GeminiProvider struct {
	client *genai.Client
}

func NewGeminiProvider(apiKey string) (*GeminiProvider, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &GeminiProvider{client: client}, nil
}

func (p *GeminiProvider) Name() string {
	return string(models.AIProviderGemini)
}

func (p *GeminiProvider) EmbedText(ctx context.Context, text string) ([]float32, error) {
	emModel := p.client.EmbeddingModel(models.ModelGeminiEmbedding)
	res, err := emModel.EmbedContent(ctx, genai.Text(text))
	if err != nil {
		return nil, err
	}
	if res.Embedding == nil {
		return nil, fmt.Errorf("no embedding returned")
	}
	return res.Embedding.Values, nil
}

func (p *GeminiProvider) GenerateReply(ctx context.Context, config *models.AIConfig, userQuery string, contextDocs []string, history []models.Message) (string, error) {
	model := p.client.GenerativeModel(models.ModelGeminiFlash)
	model.SetTemperature(0.7)

	contextStr := strings.Join(contextDocs, "\n\n")

	// Format history
	historyStr := ""
	if len(history) > 0 {
		var lines []string
		for _, m := range history {
			role := "Customer"
			if m.SenderType == models.SenderAgent {
				role = "Human Agent"
			} else if m.SenderType == models.SenderAI {
				role = "AI Assistant"
			}
			lines = append(lines, fmt.Sprintf("%s: %s", role, m.Content))
		}
		historyStr = strings.Join(lines, "\n")
	}

	prompt := fmt.Sprintf(`
SYSTEM INSTRUCTION: You are a polite, patient, and professional business assistant. Use simple, easy-to-understand language. Provide direct, relevant, and concise responses that are clear and to the point.
Additional rules: %s

CONTEXT FROM KNOWLEDGE BASE:
%s

RECENT CONVERSATION HISTORY:
%s

USER QUERY: %s

Answer the user's query based PRIMARILY on the context provided and the recent conversation history. 
Keep the tone consistent with the system instruction.

CRITICAL: Your response MUST NOT exceed 800 characters. Be concise and direct.
`, config.Persona, contextStr, historyStr, userQuery)

	// Retry logic
	var resp *genai.GenerateContentResponse
	var lastErr error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		resp, lastErr = model.GenerateContent(ctx, genai.Text(prompt))
		if lastErr == nil {
			break
		}

		if strings.Contains(lastErr.Error(), "429") || strings.Contains(lastErr.Error(), "quota") {
			log.Printf("Gemini Rate Limit Hit (429), retry %d/%d in 2s...", i+1, maxRetries)
			time.Sleep(time.Duration(2*(i+1)) * time.Second)
			continue
		}
		return "", lastErr
	}

	if lastErr != nil {
		return "", lastErr
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from AI")
	}

	if txt, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
		return string(txt), nil
	}

	return "", fmt.Errorf("unexpected response format")
}
