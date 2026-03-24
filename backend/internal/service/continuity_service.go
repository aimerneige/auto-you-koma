package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aimerneige/auto-you-koma/internal/llm"
	"github.com/aimerneige/auto-you-koma/internal/model"
	"github.com/aimerneige/auto-you-koma/internal/repository"
)

type ContinuityService struct {
	stateRepo  repository.CharacterStateRepository
	scriptRepo repository.ScriptRepository
	llm        llm.TextGenerator
}

func NewContinuityService(sr repository.CharacterStateRepository, scr repository.ScriptRepository, tgen llm.TextGenerator) *ContinuityService {
	return &ContinuityService{
		stateRepo:  sr,
		scriptRepo: scr,
		llm:        tgen,
	}
}

// ExtractContext parses what happened to a character in the last script
func (s *ContinuityService) SynthesizeMemory(ctx context.Context, seriesID string, scriptID string, characterID string) error {
	// Look up the script to analyze
	script, err := s.scriptRepo.GetByID(ctx, scriptID)
	if err != nil {
		return err
	}

	state, err := s.stateRepo.Get(ctx, seriesID, characterID)
	if err != nil {
         // Create default if not found
		state = &model.CharacterState{
			SeriesID:    seriesID,
			CharacterID: characterID,
			Health:      100,
			Sanity:      100,
		}
	}

	// Form prompt for memory digestion
	sysPrompt := "You are a Story Continuity Agent. Summarize how the events of the following script affect the target character. Reply in JSON: {\"health_change\": 0, \"items_acquired\": [\"...\"], \"items_lost\": [\"...\"], \"memory_summary\": \"...\"}\nTarget Character ID/Name relative to story: " + characterID + "\nScript Data: " + script.ParsedData

	res, err := s.llm.Generate(ctx, llm.TextRequest{Prompt: sysPrompt})
	if err != nil {
		log.Printf("Continuity memory generation failed: %v", err)
		return err
	}

	var aiAnalysis struct {
		HealthChange  int      `json:"health_change"`
		ItemsAcquired []string `json:"items_acquired"`
		ItemsLost     []string `json:"items_lost"`
		MemorySummary string   `json:"memory_summary"`
	}

	_ = json.Unmarshal([]byte(res.Content), &aiAnalysis)

	state.Health += aiAnalysis.HealthChange
	state.MemorySummary = state.MemorySummary + " | " + aiAnalysis.MemorySummary
	// Skip advanced inventory union sets logic for this MVP to avoid complex array traversals

	return s.stateRepo.Save(ctx, state)
}
