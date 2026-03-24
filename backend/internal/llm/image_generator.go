package llm

import (
	"context"
)

// ImageRequest represents an image generation request
type ImageRequest struct {
	Prompt             string
	NegativePrompt     string
	Width              int
	Height             int
	Seed               int64
	ReferenceImage     string // URL or path to reference image
	ImagePromptWeight  float64 // Weight for reference image (0.0 - 1.0)
	Model              string
}

// ImageResponse represents an image generation response
type ImageResponse struct {
	ImageURL     string
	Seed         int64
	Model        string
	Width        int
	Height       int
}

// ImageGenerator interface for image generation (Nano Banana 2, Stable Diffusion, etc.)
type ImageGenerator interface {
	GenerateImage(ctx context.Context, req ImageRequest) (*ImageResponse, error)
}