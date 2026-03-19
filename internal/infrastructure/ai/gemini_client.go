package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

type GenerationRequest struct {
	Content string `json:"content"`
}

type GenerationResponse struct {
	Text string `json:"text"`
}

func NewGeminiClient(apiKey string) (*GeminiClient, error) {
	if apiKey == "" || apiKey == "dummy-key-for-initialization" {
		return &GeminiClient{
			client: nil,
			model:  nil,
		}, nil
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("falha ao criar cliente Gemini: %w", err)
	}

	model := client.GenerativeModel("gemini-2.5-flash-lite")

	systemInstruction := `Atue como um supervisor clínico experiente. Analise os dados abaixo e identifique padrões longitudinais. Seja breve, técnico e reflexivo. Nunca dê diagnósticos fechados. Foque em:
1. Temas Dominantes: O que mais ocupou o espaço psíquico
2. Pontos de Inflexão: Mudanças notáveis na narrativa ou humor
3. Correlações Sugeridas: Relações entre eventos clínicos
4. Provocação Clínica: Uma pergunta para o terapeuta considerar na próxima sessão`

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemInstruction)},
	}

	return &GeminiClient{
		client: client,
		model:  model,
	}, nil
}

func (c *GeminiClient) GenerateWithRetry(ctx context.Context, prompt string, maxRetries int) (string, error) {
	// Check if client is nil (dummy client for when API key is not set)
	if c.client == nil || c.model == nil {
		return "", fmt.Errorf("Gemini client não inicializado. Configure a GEMINI_API_KEY no arquivo .env")
	}

	var lastErr error

	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			backoff := time.Duration(1<<uint(i)) * time.Second
			if backoff > 30*time.Second {
				backoff = 30 * time.Second
			}

			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(backoff):
			}
		}

		result, err := c.generateSingle(ctx, prompt)
		if err == nil {
			return result, nil
		}

		lastErr = err

		if shouldRetry(err) {
			continue
		}

		break
	}

	return "", fmt.Errorf("falha após %d tentativas: %w", maxRetries, lastErr)
}

func (c *GeminiClient) generateSingle(ctx context.Context, prompt string) (string, error) {
	resp, err := c.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("erro na geração: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("resposta vazia da API")
	}

	var result strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			result.WriteString(string(text))
		}
	}

	return result.String(), nil
}

func shouldRetry(err error) bool {
	errStr := err.Error()

	if strings.Contains(errStr, "429") ||
		strings.Contains(errStr, "rate limit") ||
		strings.Contains(errStr, "quota") ||
		strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "network") {
		return true
	}

	return false
}

func (c *GeminiClient) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
