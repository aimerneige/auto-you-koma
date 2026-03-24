package llm

import (
	"context"
)

// TextRequest represents a text generation request
type TextRequest struct {
	Prompt      string
	SystemPrompt string
	Model       string
	MaxTokens   int
	Temperature float64
}

// TextResponse represents a text generation response
type TextResponse struct {
	Content   string
	TokenCount int
	Model     string
}

// TextGenerator interface for text generation (Gemini, OpenAI, Deepseek, etc.)
type TextGenerator interface {
	Generate(ctx context.Context, req TextRequest) (*TextResponse, error)
	GenerateStream(ctx context.Context, req TextRequest) (<-chan string, error)
}