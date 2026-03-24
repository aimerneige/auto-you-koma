package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aimerneige/auto-you-koma/internal/config"
)

type GenericImageGenerator struct {
	cfg config.ImageLLMConfig
}

func NewGenericImageGenerator(cfg config.ImageLLMConfig) ImageGenerator {
	return &GenericImageGenerator{cfg: cfg}
}

type openAIImageReq struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

type openAIImageRes struct {
	Data []struct {
		URL string `json:"url"`
	} `json:"data"`
}

func (g *GenericImageGenerator) GenerateImage(ctx context.Context, req ImageRequest) (*ImageResponse, error) {
	// If API key is empty or missing provider configuration, simulate local generation (Mock API)
	if g.cfg.Provider == "" || g.cfg.NanoBanana.APIKey == "" {
		// Mock delay processing for 2 seconds to simulate network I/O
		time.Sleep(2 * time.Second)
		mockURL := fmt.Sprintf("https://placehold.co/%dx%d/png?text=Mock+Image", req.Width, req.Height)
		return &ImageResponse{ImageURL: mockURL}, nil
	}

	// For Nano Banana, if it follows standard OpenAI image structure, we can hit it this way:
	// Example uses OpenAI standard fallback
	url := "https://api.openai.com/v1/images/generations" 
	
	// Default size logic fallback
	size := "1024x1024"
	if req.Width == 512 && req.Height == 512 {
		size = "512x512"
	}

	reqBody := openAIImageReq{
		Model:  "dall-e-3",
		Prompt: req.Prompt,
		N:      1,
		Size:   size,
	}

	jsonData, _ := json.Marshal(reqBody)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+g.cfg.NanoBanana.APIKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("image generation error: status %d, body %s", resp.StatusCode, string(body))
	}

	var resData openAIImageRes
	if err := json.NewDecoder(resp.Body).Decode(&resData); err != nil {
		return nil, err
	}

	if len(resData.Data) == 0 {
		return nil, errors.New("no image returned")
	}

	return &ImageResponse{ImageURL: resData.Data[0].URL}, nil
}
