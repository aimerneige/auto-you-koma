package llm

import (
	"context"
	"fmt"
)

// MockImageGenerator is a mock implementation for testing
type MockImageGenerator struct{}

func (m *MockImageGenerator) GenerateImage(ctx context.Context, req ImageRequest) (*ImageResponse, error) {
	// Return mock response
	return &ImageResponse{
		ImageURL: "https://example.com/mock-image.png",
		Seed:     req.Seed,
		Model:    "mock",
		Width:    req.Width,
		Height:   req.Height,
	}, nil
}

// NanoBananaImageGenerator implements ImageGenerator using Nano Banana 2 API
type NanoBananaImageGenerator struct {
	APIKey string
	Model  string
}

func NewNanoBananaImageGenerator(apiKey string) *NanoBananaImageGenerator {
	return &NanoBananaImageGenerator{
		APIKey: apiKey,
		Model:  "nano-banana-2",
	}
}

func (n *NanoBananaImageGenerator) GenerateImage(ctx context.Context, req ImageRequest) (*ImageResponse, error) {
	// TODO: Implement actual Nano Banana 2 API call
	// This is a placeholder implementation
	// Key: Use ReferenceImage and ImagePromptWeight for character consistency
	return &ImageResponse{
		ImageURL: "https://example.com/generated-image.png",
		Seed:     req.Seed,
		Model:    n.Model,
		Width:    req.Width,
		Height:   req.Height,
	}, nil
}

// ImagePromptWeight constant values
const (
	ImagePromptWeightLow    = 0.3
	ImagePromptWeightMedium = 0.6
	ImagePromptWeightHigh   = 0.9
)

// FormatPromptWithReference formats the prompt with reference image for character consistency
func FormatPromptWithReference(prompt string, referenceImage string, weight float64) string {
	if referenceImage == "" {
		return prompt
	}
	return fmt.Sprintf("%s [reference:%s:%.2f]", prompt, referenceImage, weight)
}