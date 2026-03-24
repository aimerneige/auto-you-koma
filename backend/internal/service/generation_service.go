package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/aimerneige/auto-you-koma/internal/compositor"
	"github.com/aimerneige/auto-you-koma/internal/llm"
	"github.com/aimerneige/auto-you-koma/internal/model"
	"github.com/aimerneige/auto-you-koma/internal/repository"
)

type GenerationService struct {
	genRepo    repository.GenerationRepository
	scriptRepo repository.ScriptRepository
	imageGen   llm.ImageGenerator
	compositor *compositor.Compositor
}

func NewGenerationService(g repository.GenerationRepository, s repository.ScriptRepository, ig llm.ImageGenerator, comp *compositor.Compositor) *GenerationService {
	return &GenerationService{
		genRepo:    g,
		scriptRepo: s,
		imageGen:   ig,
		compositor: comp,
	}
}

func (s *GenerationService) Get(ctx context.Context, id string) (*model.Generation, error) {
	return s.genRepo.GetByID(ctx, id)
}

func (s *GenerationService) StartGeneration(ctx context.Context, projectID, scriptID, layout string) (*model.Generation, error) {
	gen := &model.Generation{
		ID:        uuid.New().String(),
		ProjectID: projectID,
		ScriptID:  scriptID,
		Status:    model.GenerationPending,
		Layout:    layout,
	}

	if err := s.genRepo.Create(ctx, gen); err != nil {
		return nil, err
	}

	// Kick off background processing
	go s.processAsync(gen.ID, scriptID, layout)

	return gen, nil
}

func (s *GenerationService) processAsync(genID, scriptID, layout string) {
	ctx := context.Background() // new context for background work

	gen, err := s.genRepo.GetByID(ctx, genID)
	if err != nil {
		log.Printf("Fail to get generation info: %v\n", err)
		return
	}

	gen.Status = model.GenerationProcessing
	_ = s.genRepo.Update(ctx, gen)

	failJob := func(errMsg string) {
		gen.Status = model.GenerationFailed
		gen.Error = errMsg
		_ = s.genRepo.Update(ctx, gen)
	}

	script, err := s.scriptRepo.GetByID(ctx, scriptID)
	if err != nil {
		failJob("Script not found: " + err.Error())
		return
	}

	var panels []PanelData // defined in script_service
	if script.ParsedData == "" {
		failJob("Script has no parsed storyboard data")
		return
	}

	if err := json.Unmarshal([]byte(script.ParsedData), &panels); err != nil {
		failJob("Failed to unmarshal parsed data: " + err.Error())
		return
	}

	var imageUrls []string
	
	for _, p := range panels {
		prompt := "A comic panel. " + p.VisualDesc
		res, err := s.imageGen.GenerateImage(ctx, llm.ImageRequest{
			Prompt: prompt,
			Width:  512,
			Height: 512,
		})
		
		if err != nil {
			// fallback/mock retry inside processing loop not handled for brevity
			log.Printf("Image partial fail: %v\n", err)
			continue
		}
		
		imageUrls = append(imageUrls, res.ImageURL)
		time.Sleep(1 * time.Second) // rate limit backoff
	}

	if len(imageUrls) == 0 {
		failJob("All image generations failed")
		return
	}

	resultUrl, err := s.compositor.Compose(imageUrls, layout)
	if err != nil {
		failJob("Compositor failed: " + err.Error())
		return
	}

	gen.Status = model.GenerationDone
	gen.ResultImageURL = resultUrl
	_ = s.genRepo.Update(ctx, gen)
}
