package llm

import "context"

type TextRequest struct {
	Prompt string
}

type TextResponse struct {
	Content string
}

type TextChunk struct {
	Content string
}

type ImageRequest struct {
	Prompt         string
	NegativePrompt string
	Width          int
	Height         int
	Seed           int64
}

type ImageResponse struct {
	ImageURL string
}

type VisionRequest struct {
	ImageURL string
	Prompt   string
}

type VisionResponse struct {
	Analysis string
}

// LLM 文本生成接口
type TextGenerator interface {
	Generate(ctx context.Context, req TextRequest) (*TextResponse, error)
	GenerateStream(ctx context.Context, req TextRequest) (<-chan TextChunk, error)
}

// LLM 图像生成接口
type ImageGenerator interface {
	GenerateImage(ctx context.Context, req ImageRequest) (*ImageResponse, error)
}

// 多模态分析接口（用于 QC Agent）
type VisionAnalyzer interface {
	Analyze(ctx context.Context, req VisionRequest) (*VisionResponse, error)
}
