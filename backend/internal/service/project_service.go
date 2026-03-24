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
	repo         *repository.SQLiteProjectRepo
	scriptRepo   *repository.SQLiteScriptRepo
	charRepo     *repository.SQLiteCharacterRepo
	textGenerator llm.TextGenerator
}

func NewProjectService(
	repo *repository.SQLiteProjectRepo,
	scriptRepo *repository.SQLiteScriptRepo,
	charRepo *repository.SQLiteCharacterRepo,
	textGen llm.TextGenerator,
) *ProjectService {
	return &ProjectService{
		repo:         repo,
		scriptRepo:   scriptRepo,
		charRepo:     charRepo,
		textGenerator: textGen,
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