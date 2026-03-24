package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aimerneige/auto-you-koma/internal/llm"
	"github.com/aimerneige/auto-you-koma/internal/models"
	"github.com/aimerneige/auto-you-koma/internal/repository"
)

// ProjectService provides project business logic
type ProjectService struct {
	repo           *repository.SQLiteProjectRepo
	scriptRepo     *repository.SQLiteScriptRepo
	charRepo       *repository.SQLiteCharacterRepo
	storyboardRepo *repository.SQLiteStoryboardRepo
	textGenerator  llm.TextGenerator
}

func NewProjectService(
	repo *repository.SQLiteProjectRepo,
	scriptRepo *repository.SQLiteScriptRepo,
	charRepo *repository.SQLiteCharacterRepo,
	storyboardRepo *repository.SQLiteStoryboardRepo,
	textGen llm.TextGenerator,
) *ProjectService {
	return &ProjectService{
		repo:           repo,
		scriptRepo:     scriptRepo,
		charRepo:        charRepo,
		storyboardRepo:  storyboardRepo,
		textGenerator:   textGen,
	}
}

type CreateProjectRequest struct {
	UserID   string
	Title    string
	Mode     string
	Synopsis string
}

// ScriptPanel represents one panel in the script
type ScriptPanel struct {
	PanelNumber    int    `json:"panel_number"`
	Structure      string `json:"structure"` // 起承转合
	Scene          string `json:"scene"`
	Characters     string `json:"characters"`
	Dialogue       string `json:"dialogue"`
	Narration      string `json:"narration,omitempty"`
}

// ScriptContent represents the full script structure
type ScriptContent struct {
	ProjectID   string       `json:"project_id"`
	Title       string       `json:"title"`
	Synopsis    string       `json:"synopsis"`
	Mode        string       `json:"mode"`
	Panels      []ScriptPanel `json:"panels"`
}

func (s *ProjectService) Create(ctx context.Context, req CreateProjectRequest) (*models.Project, error) {
	project := &models.Project{
		UserID:   req.UserID,
		Title:    req.Title,
		Mode:     req.Mode,
		Synopsis: req.Synopsis,
		Status:   "draft",
	}
	err := s.repo.Create(ctx, project)
	return project, err
}

func (s *ProjectService) GetByID(ctx context.Context, id string) (*models.Project, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProjectService) List(ctx context.Context, userID string, limit, offset int) ([]*models.Project, error) {
	return s.repo.List(ctx, userID, limit, offset)
}

func (s *ProjectService) Update(ctx context.Context, id string, req CreateProjectRequest) (*models.Project, error) {
	project, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	project.Title = req.Title
	project.Mode = req.Mode
	project.Synopsis = req.Synopsis
	err = s.repo.Update(ctx, project)
	return project, err
}

func (s *ProjectService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// GenerateScript generates a 4-panel script from the project synopsis and characters
func (s *ProjectService) GenerateScript(ctx context.Context, projectID string) (*models.Script, error) {
	project, err := s.repo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Get characters for this project
	characters, err := s.charRepo.List(ctx, project.UserID, 10, 0)
	if err != nil {
		characters = []*models.Character{}
	}

	// Build the prompt for script generation
	charInfo := ""
	for _, c := range characters {
		charInfo += fmt.Sprintf("- %s (%s): %s\n", c.Name, c.NameJP, c.Personality)
	}

	prompt := fmt.Sprintf(`Generate a 4-panel comic script based on the following:
Title: %s
Synopsis: %s

Character Information:
%s

Generate a JSON with exactly 4 panels following the structure:
{
  "project_id": "%s",
  "title": "...",
  "synopsis": "...",
  "mode": "%s",
  "panels": [
    {"panel_number": 1, "structure": "起", "scene": "...", "characters": "...", "dialogue": "..."},
    {"panel_number": 2, "structure": "承", "scene": "...", "characters": "...", "dialogue": "..."},
    {"panel_number": 3, "structure": "转", "scene": "...", "characters": "...", "dialogue": "..."},
    {"panel_number": 4, "structure": "合", "scene": "...", "characters": "...", "dialogue": "..."}
  ]
}

Return ONLY the JSON, no other text.`, project.Title, project.Synopsis, charInfo, projectID, project.Mode)

	// Call the LLM to generate the script
	resp, err := s.textGenerator.Generate(ctx, llm.TextRequest{
		Prompt:      prompt,
		MaxTokens:   2000,
		Temperature: 0.7,
	})
	if err != nil {
		return nil, err
	}

	// Parse the response
	var scriptContent ScriptContent
	if err := json.Unmarshal([]byte(resp.Content), &scriptContent); err != nil {
		// If JSON parsing fails, create a simple default script
		scriptContent = ScriptContent{
			ProjectID: projectID,
			Title:     project.Title,
			Synopsis:  project.Synopsis,
			Mode:      project.Mode,
			Panels: []ScriptPanel{
				{PanelNumber: 1, Structure: "起", Scene: "Scene 1", Characters: "", Dialogue: ""},
				{PanelNumber: 2, Structure: "承", Scene: "Scene 2", Characters: "", Dialogue: ""},
				{PanelNumber: 3, Structure: "转", Scene: "Scene 3", Characters: "", Dialogue: ""},
				{PanelNumber: 4, Structure: "合", Scene: "Scene 4", Characters: "", Dialogue: ""},
			},
		}
	}

	// Save the script
	contentJSON, err := json.Marshal(scriptContent)
	if err != nil {
		return nil, err
	}

	script := &models.Script{
		ProjectID: projectID,
		Content:   string(contentJSON),
	}

	err = s.scriptRepo.Create(ctx, script)
	if err != nil {
		return nil, err
	}

	// Update project status
	project.Status = "scripted"
	s.repo.Update(ctx, project)

	return script, nil
}

// GetScript retrieves the script for a project
func (s *ProjectService) GetScript(ctx context.Context, projectID string) (*models.Script, error) {
	return s.scriptRepo.GetByProjectID(ctx, projectID)
}

// UpdateScript updates the script content
func (s *ProjectService) UpdateScript(ctx context.Context, scriptID string, content string) (*models.Script, error) {
	script, err := s.scriptRepo.GetByID(ctx, scriptID)
	if err != nil {
		return nil, err
	}
	script.Content = content
	script.Version++
	err = s.scriptRepo.Update(ctx, script)
	return script, err
}

// StoryboardPanel represents one panel in the storyboard
type StoryboardPanel struct {
	PanelNumber    int    `json:"panel_number"`
	ShotType       string `json:"shot_type"`       // wide_shot, medium_shot, close_up, etc.
	Angle          string `json:"angle"`           // eye_level, high_angle, low_angle
	Atmosphere     string `json:"atmosphere"`      // Description of lighting/mood
	Characters     string `json:"characters"`      // Characters in this panel
	Action         string `json:"action"`         // Action description
	Expression     string `json:"expression"`      // Facial expressions
	PositivePrompt string `json:"positive_prompt"` // AI image generation prompt
	NegativePrompt string `json:"negative_prompt"` // Negative prompt
	Dialogue       string `json:"dialogue"`       // Dialogue from script
	Seed           int64  `json:"seed"`
}

// StoryboardContent represents the full storyboard structure
type StoryboardContent struct {
	StoryboardID string            `json:"storyboard_id"`
	ScriptID     string           `json:"script_id"`
	Panels       []StoryboardPanel `json:"panels"`
}

// GenerateStoryboard generates detailed storyboard from script
func (s *ProjectService) GenerateStoryboard(ctx context.Context, projectID string) (*models.Storyboard, error) {
	// Get the script
	script, err := s.scriptRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Parse script content
	var scriptContent ScriptContent
	if err := json.Unmarshal([]byte(script.Content), &scriptContent); err != nil {
		return nil, err
	}

	// Build prompt for storyboard generation
	prompt := fmt.Sprintf(`Convert the following 4-panel script into a detailed storyboard JSON with visual rendering parameters:

Script Title: %s
Synopsis: %s

Panels:
%s

Generate a JSON with exactly 4 panels following this structure:
{
  "storyboard_id": "uuid",
  "script_id": "%s",
  "panels": [
    {
      "panel_number": 1,
      "shot_type": "medium_shot",
      "angle": "eye_level",
      "atmosphere": "bright, sunny day",
      "characters": "character names",
      "action": "action description",
      "expression": "facial expression",
      "positive_prompt": "detailed positive prompt for AI image generation",
      "negative_prompt": "bad anatomy, extra limbs, blurry",
      "dialogue": "dialogue text",
      "seed": 12345
    }
  ]
}

Return ONLY the JSON, no other text.`, scriptContent.Title, scriptContent.Synopsis, formatScriptPanels(scriptContent.Panels), script.ID)

	// Call the LLM
	resp, err := s.textGenerator.Generate(ctx, llm.TextRequest{
		Prompt:      prompt,
		MaxTokens:   3000,
		Temperature: 0.7,
	})
	if err != nil {
		return nil, err
	}

	// Parse the response
	var storyboardContent StoryboardContent
	if err := json.Unmarshal([]byte(resp.Content), &storyboardContent); err != nil {
		// Create default storyboard if parsing fails
		storyboardContent = createDefaultStoryboard(script.ID)
	}

	// Save the storyboard
	contentJSON, err := json.Marshal(storyboardContent)
	if err != nil {
		return nil, err
	}

	storyboard := &models.Storyboard{
		ScriptID: script.ID,
		Content:  string(contentJSON),
	}

	err = s.storyboardRepo.Create(ctx, storyboard)
	if err != nil {
		return nil, err
	}

	// Update project status
	project, _ := s.repo.GetByID(ctx, projectID)
	if project != nil {
		project.Status = "previewed"
		s.repo.Update(ctx, project)
	}

	return storyboard, nil
}

func formatScriptPanels(panels []ScriptPanel) string {
	result := ""
	for _, p := range panels {
		result += fmt.Sprintf("Panel %d (%s): Scene: %s, Characters: %s, Dialogue: %s\n",
			p.PanelNumber, p.Structure, p.Scene, p.Characters, p.Dialogue)
	}
	return result
}

func createDefaultStoryboard(scriptID string) StoryboardContent {
	panels := make([]StoryboardPanel, 4)
	for i := 0; i < 4; i++ {
		panels[i] = StoryboardPanel{
			PanelNumber:    i + 1,
			ShotType:       "medium_shot",
			Angle:          "eye_level",
			Atmosphere:     "normal",
			PositivePrompt: "anime style, clean background",
			NegativePrompt: "bad anatomy, blurry",
			Seed:           int64(1000 + i*100),
		}
	}
	return StoryboardContent{
		ScriptID: scriptID,
		Panels:   panels,
	}
}

// GetStoryboard retrieves the storyboard for a project
func (s *ProjectService) GetStoryboard(ctx context.Context, projectID string) (*models.Storyboard, error) {
	script, err := s.scriptRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return s.storyboardRepo.GetByScriptID(ctx, script.ID)
}

// UpdateStoryboard updates the storyboard content
func (s *ProjectService) UpdateStoryboard(ctx context.Context, storyboardID string, content string) (*models.Storyboard, error) {
	storyboard, err := s.storyboardRepo.GetByID(ctx, storyboardID)
	if err != nil {
		return nil, err
	}
	storyboard.Content = content
	storyboard.Version++
	err = s.storyboardRepo.Update(ctx, storyboard)
	return storyboard, err
}