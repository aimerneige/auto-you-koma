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
	qcReviewer llm.VisionAnalyzer
	compositor *compositor.Compositor
}

func NewGenerationService(g repository.GenerationRepository, s repository.ScriptRepository, ig llm.ImageGenerator, qc llm.VisionAnalyzer, comp *compositor.Compositor) *GenerationService {
	return &GenerationService{
		genRepo:    g,
		scriptRepo: s,
		imageGen:   ig,
		qcReviewer: qc,
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
		
		var confirmedUrl string
		maxRetries := 2
		
		for attempt := 0; attempt <= maxRetries; attempt++ {
			req := llm.ImageRequest{
				Prompt: prompt,
				Width:  512,
				Height: 512,
			}
			
			// If we wanted to lock seed, we could set: req.Seed = pointerToSeed
			res, err := s.imageGen.GenerateImage(ctx, req)
			
			if err != nil {
				log.Printf("Image gen fail on panel %d attempt %d: %v\n", p.Panel, attempt, err)
				continue
			}
			
			// Call QC Reviewer Agent 6
			qcErr := s.qcReviewer.AnalyzeImage(ctx, res.ImageURL, prompt)
			if qcErr != nil {
				log.Printf("QC Reject panel %d, redraw triggered: %v\n", p.Panel, qcErr)
				continue // retry
			}
			
			confirmedUrl = res.ImageURL
			break
		}
		
		if confirmedUrl == "" {
			log.Printf("Giving up on panel %d after %d retries\n", p.Panel, maxRetries)
			continue // Accept missing or handle explicitly
		}

		imageUrls = append(imageUrls, confirmedUrl)
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

	rawBytes, _ := json.Marshal(imageUrls)

	gen.Status = model.GenerationDone
	gen.ResultImageURL = resultUrl
	gen.RawImageURLs = string(rawBytes)
	_ = s.genRepo.Update(ctx, gen)
}
