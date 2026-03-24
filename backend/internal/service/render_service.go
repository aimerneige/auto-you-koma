package service

import (
	"context"
	"encoding/json"

	"github.com/aimerneige/auto-you-koma/internal/llm"
	"github.com/aimerneige/auto-you-koma/internal/models"
	"github.com/aimerneige/auto-you-koma/internal/repository"
)

// RenderService handles image rendering operations
type RenderService struct {
	renderTaskRepo *repository.SQLiteRenderTaskRepo
	charRepo       *repository.SQLiteCharacterRepo
	storyboardRepo *repository.SQLiteStoryboardRepo
	projectRepo    *repository.SQLiteProjectRepo
	scriptRepo     *repository.SQLiteScriptRepo
	imageGenerator llm.ImageGenerator
}

func NewRenderService(
	renderTaskRepo *repository.SQLiteRenderTaskRepo,
	charRepo *repository.SQLiteCharacterRepo,
	storyboardRepo *repository.SQLiteStoryboardRepo,
	projectRepo *repository.SQLiteProjectRepo,
	scriptRepo *repository.SQLiteScriptRepo,
	imgGen llm.ImageGenerator,
) *RenderService {
	return &RenderService{
		renderTaskRepo: renderTaskRepo,
		charRepo:       charRepo,
		storyboardRepo: storyboardRepo,
		projectRepo:    projectRepo,
		scriptRepo:     scriptRepo,
		imageGenerator: imgGen,
	}
}

type RenderRequest struct {
	ProjectID    string `json:"project_id"`
	ScriptID     string `json:"script_id"`
	ExportType   string `json:"export_type"`   // native_text / clean_plate
	Layout       string `json:"layout"`         // 2x2 / 1x4
	ImageWidth   int    `json:"image_width"`   // default 1024
	ImageHeight  int    `json:"image_height"`  // default 1024
}

type PanelRenderResult struct {
	PanelNumber int    `json:"panel_number"`
	ImageURL    string `json:"image_url"`
	Seed        int64  `json:"seed"`
	Success     bool   `json:"success"`
	Error       string `json:"error,omitempty"`
}

// StartRender initiates the rendering process for a project
func (s *RenderService) StartRender(ctx context.Context, req RenderRequest) (*models.RenderTask, error) {
	// Get storyboard
	storyboard, err := s.storyboardRepo.GetByScriptID(ctx, req.ScriptID)
	if err != nil {
		return nil, err
	}

	// Parse storyboard content
	var storyboardContent StoryboardContent
	if err := json.Unmarshal([]byte(storyboard.Content), &storyboardContent); err != nil {
		return nil, err
	}

	// Get characters for reference images
	characters, err := s.charRepo.List(ctx, "", 10, 0)
	if err != nil {
		characters = []*models.Character{}
	}

	// Create render task
	task := &models.RenderTask{
		ProjectID:    req.ProjectID,
		StoryboardID: storyboard.ID,
		ExportType:   req.ExportType,
		Layout:       req.Layout,
		ImageWidth:   req.ImageWidth,
		ImageHeight:  req.ImageHeight,
		Status:       "rendering",
	}

	err = s.renderTaskRepo.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	// Render each panel
	results := make([]PanelRenderResult, len(storyboardContent.Panels))
	for i, panel := range storyboardContent.Panels {
		results[i] = s.renderPanel(ctx, panel, characters, req.ImageWidth, req.ImageHeight)
	}

	// Save results
	resultsJSON, _ := json.Marshal(results)
	task.OutputPaths = string(resultsJSON)
	task.Status = "done"

	s.renderTaskRepo.Update(ctx, task)

	return task, nil
}

func (s *RenderService) renderPanel(ctx context.Context, panel StoryboardPanel, characters []*models.Character, width, height int) PanelRenderResult {
	// Find reference images for characters in this panel
	referenceImage := ""
	referenceWeight := llm.ImagePromptWeightMedium

	// Look for character reference sheets
	for _, char := range characters {
		if char.ReferenceSheetURL != "" && containsChar(panel.Characters, char.Name) {
			referenceImage = char.ReferenceSheetURL
			break
		}
	}

	// Build the prompt
	prompt := panel.PositivePrompt
	if referenceImage != "" {
		prompt = llm.FormatPromptWithReference(prompt, referenceImage, referenceWeight)
	}

	// Generate image
	req := llm.ImageRequest{
		Prompt:         prompt,
		NegativePrompt: panel.NegativePrompt,
		Width:          width,
		Height:         height,
		Seed:           panel.Seed,
		ReferenceImage: referenceImage,
		ImagePromptWeight: referenceWeight,
	}

	resp, err := s.imageGenerator.GenerateImage(ctx, req)
	if err != nil {
		return PanelRenderResult{
			PanelNumber: panel.PanelNumber,
			Success:     false,
			Error:       err.Error(),
		}
	}

	return PanelRenderResult{
		PanelNumber: panel.PanelNumber,
		ImageURL:    resp.ImageURL,
		Seed:        resp.Seed,
		Success:     true,
	}
}

// RegenerateSinglePanel regenerates a single panel
func (s *RenderService) RegenerateSinglePanel(ctx context.Context, taskID string, panelNumber int) (*PanelRenderResult, error) {
	task, err := s.renderTaskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Get storyboard
	storyboard, err := s.storyboardRepo.GetByID(ctx, task.StoryboardID)
	if err != nil {
		return nil, err
	}

	var storyboardContent StoryboardContent
	json.Unmarshal([]byte(storyboard.Content), &storyboardContent)

	// Get characters
	characters, _ := s.charRepo.List(ctx, "", 10, 0)

	// Find the panel
	var panel StoryboardPanel
	for _, p := range storyboardContent.Panels {
		if p.PanelNumber == panelNumber {
			panel = p
			break
		}
	}

	// Render with new seed
	result := s.renderPanel(ctx, panel, characters, task.ImageWidth, task.ImageHeight)
	result.Seed = result.Seed + 1 // Change seed for regeneration

	// Update task
	var results []PanelRenderResult
	json.Unmarshal([]byte(task.OutputPaths), &results)
	for i := range results {
		if results[i].PanelNumber == panelNumber {
			results[i] = result
			break
		}
	}
	resultsJSON, _ := json.Marshal(results)
	task.OutputPaths = string(resultsJSON)
	s.renderTaskRepo.Update(ctx, task)

	return &result, nil
}

// GetRenderTask retrieves the render task for a project
func (s *RenderService) GetRenderTask(ctx context.Context, projectID string) (*models.RenderTask, error) {
	return s.renderTaskRepo.GetByProjectID(ctx, projectID)
}

// ConfirmRender confirms the render and updates project status
func (s *RenderService) ConfirmRender(ctx context.Context, projectID string) error {
	task, err := s.renderTaskRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return err
	}
	task.Status = "done"
	err = s.renderTaskRepo.Update(ctx, task)
	if err != nil {
		return err
	}
	// Update project status
	project, _ := s.projectRepo.GetByID(ctx, projectID)
	if project != nil {
		project.Status = "done"
		s.projectRepo.Update(ctx, project)
	}
	return nil
}

// GetProject retrieves a project by ID
func (s *RenderService) GetProject(ctx context.Context, id string) (*models.Project, error) {
	return s.projectRepo.GetByID(ctx, id)
}

// GetScript retrieves a script by project ID
func (s *RenderService) GetScript(ctx context.Context, projectID string) (*models.Script, error) {
	return s.scriptRepo.GetByProjectID(ctx, projectID)
}

func containsChar(text, name string) bool {
	return len(name) > 0 && len(text) > 0 &&
		(len(text) >= len(name) && (text == name || contains(text, name)))
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}