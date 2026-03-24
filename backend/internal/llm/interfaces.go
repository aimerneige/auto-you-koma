package llm

import "context"

type ImageRequest struct {
	Prompt             string
	Width              int
	Height             int
	Seed               *int64   // Agent 4 Seed locker
	ReferenceImageURLs []string // Image-to-Image or ControlNet references
	ControlNetWeight   float64  // Strength of reference preservation
}

type ImageResponse struct {
	ImageURL string
}

type VisionAnalyzer interface {
	AnalyzeImage(ctx context.Context, imageURL, originalPrompt string) error
}

type ImageGenerator interface {
	GenerateImage(ctx context.Context, req ImageRequest) (*ImageResponse, error)
}

// LLM 文本生成接口
type TextRequest struct {
	Prompt string
}

type TextResponse struct {
	Content string
}

type TextChunk struct {
	Content string
}

type TextGenerator interface {
	Generate(ctx context.Context, req TextRequest) (*TextResponse, error)
	GenerateStream(ctx context.Context, req TextRequest) (<-chan TextChunk, error)
}
