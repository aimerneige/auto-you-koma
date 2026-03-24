package llm

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"
)

type MockVisionAnalyzer struct{}

func NewMockVisionAnalyzer() *MockVisionAnalyzer {
	return &MockVisionAnalyzer{}
}

func (m *MockVisionAnalyzer) AnalyzeImage(ctx context.Context, imageURL, originalPrompt string) error {
	log.Printf("[Agent 6 QC Reviewer] Analyzing generated asset: %s", imageURL)
	time.Sleep(1 * time.Second) // Simulate AI vision processing time

	// Randomly reject 30% of generations to simulate AI hallucination detection
	if rand.Float32() < 0.3 {
		log.Println("[Agent 6 QC Reviewer] REJECTED. Hallucination or anatomy error detected. Requesting redraw.")
		return errors.New("QC Failed: Anatomy defect or character unfaithful to prompt")
	}

	log.Println("[Agent 6 QC Reviewer] APPROVED.")
	return nil
}
