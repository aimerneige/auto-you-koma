package llm

import (
	"context"
	"fmt"
)

// MockTextGenerator is a mock implementation for testing
type MockTextGenerator struct {
	Response string
}

func (m *MockTextGenerator) Generate(ctx context.Context, req TextRequest) (*TextResponse, error) {
	return &TextResponse{
		Content:   m.Response,
		TokenCount: len(m.Response),
		Model:     "mock",
	}, nil
}

func (m *MockTextGenerator) GenerateStream(ctx context.Context, req TextRequest) (<-chan string, error) {
	ch := make(chan string, 1)
	go func() {
		ch <- m.Response
		close(ch)
	}()
	return ch, nil
}

// GeminiTextGenerator implements TextGenerator using Google Gemini API
type GeminiTextGenerator struct {
	APIKey string
	Model  string
}

func NewGeminiTextGenerator(apiKey string) *GeminiTextGenerator {
	return &GeminiTextGenerator{
		APIKey: apiKey,
		Model:  "gemini-2.0-flash",
	}
}

func (g *GeminiTextGenerator) Generate(ctx context.Context, req TextRequest) (*TextResponse, error) {
	// TODO: Implement actual Gemini API call
	// This is a placeholder implementation
	return &TextResponse{
		Content:   fmt.Sprintf("Mock response for: %s", req.Prompt),
		TokenCount: 100,
		Model:     g.Model,
	}, nil
}

func (g *GeminiTextGenerator) GenerateStream(ctx context.Context, req TextRequest) (<-chan string, error) {
	ch := make(chan string, 1)
	go func() {
		ch <- fmt.Sprintf("Mock stream response for: %s", req.Prompt)
		close(ch)
	}()
	return ch, nil
}