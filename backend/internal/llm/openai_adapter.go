package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/aimerneige/auto-you-koma/internal/config"
)

type OpenAIGenerator struct {
	cfg config.OpenAIConfig
}

func NewOpenAIGenerator(cfg config.OpenAIConfig) TextGenerator {
	return &OpenAIGenerator{cfg: cfg}
}

type openAIRequest struct {
	Model    string        `json:"model"`
	Messages []interface{} `json:"messages"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (g *OpenAIGenerator) Generate(ctx context.Context, req TextRequest) (*TextResponse, error) {
	url := g.cfg.BaseURL + "/v1/chat/completions"
	reqBody := openAIRequest{
		Model: g.cfg.Model,
		Messages: []interface{}{
			map[string]string{"role": "user", "content": req.Prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if g.cfg.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+g.cfg.APIKey)
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openai error: status %d, body %s", resp.StatusCode, string(body))
	}

	var resData openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&resData); err != nil {
		return nil, err
	}

	if len(resData.Choices) == 0 {
		return nil, errors.New("no completion choices returned")
	}

	return &TextResponse{Content: resData.Choices[0].Message.Content}, nil
}

func (g *OpenAIGenerator) GenerateStream(ctx context.Context, req TextRequest) (<-chan TextChunk, error) {
	return nil, errors.New("streaming not currently implemented")
}
