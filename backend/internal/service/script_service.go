package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/aimerneige/auto-you-koma/internal/llm"
	"github.com/aimerneige/auto-you-koma/internal/model"
	"github.com/aimerneige/auto-you-koma/internal/repository"
)

type ScriptService struct {
	repo repository.ScriptRepository
	llm  llm.TextGenerator
}

func NewScriptService(repo repository.ScriptRepository, textGen llm.TextGenerator) *ScriptService {
	return &ScriptService{repo: repo, llm: textGen}
}

func (s *ScriptService) Create(ctx context.Context, script *model.Script) error {
	script.ID = uuid.New().String()
	return s.repo.Create(ctx, script)
}

func (s *ScriptService) GetByID(ctx context.Context, id string) (*model.Script, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ScriptService) Update(ctx context.Context, script *model.Script) error {
	return s.repo.Update(ctx, script)
}

func (s *ScriptService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *ScriptService) ListByProject(ctx context.Context, projectID string) ([]*model.Script, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *ScriptService) GenerateScript(ctx context.Context, prompt string) (string, error) {
	sysPrompt := "You are a professional comic scriptwriter. Create a short, engaging 4-panel comic script based on the following prompt:\n\n" + prompt
	
	res, err := s.llm.Generate(ctx, llm.TextRequest{Prompt: sysPrompt})
	if err != nil {
		return "", err
	}
	return res.Content, nil
}

func (s *ScriptService) ParseToPanels(ctx context.Context, scriptID string) error {
	script, err := s.repo.GetByID(ctx, scriptID)
	if err != nil {
		return err
	}

	sysPrompt := `You are an AI assistant that parses a story into a structured JSON array for a 4-panel comic. 
Please read the story below, and extract precisely 4 distinct panels. Reply ONLY with valid JSON in this exact format (do not use markdown blocks, just raw JSON array starting with [ and ending with ]):
[
  { "panel": 1, "visual_desc": "description of the scene", "dialog": "character dialog" },
  { "panel": 2, "visual_desc": "...", "dialog": "..." },
  { "panel": 3, "visual_desc": "...", "dialog": "..." },
  { "panel": 4, "visual_desc": "...", "dialog": "..." }
]
Story:
` + script.Content

	res, err := s.llm.Generate(ctx, llm.TextRequest{Prompt: sysPrompt})
	if err != nil {
		return err
	}

	content := strings.TrimSpace(res.Content)
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)
	}

	script.ParsedData = content
	return s.repo.Update(ctx, script)
}

type PanelData struct {
	Panel      int    `json:"panel"`
	VisualDesc string `json:"visual_desc"`
	Dialog     string `json:"dialog"`
	Locked     bool   `json:"locked"`
	LayoutType string `json:"layout_type"` // Agent 7 feature
}

func (s *ScriptService) UpdatePanel(ctx context.Context, scriptID string, panelIndex int, panel PanelData) error {
	script, err := s.repo.GetByID(ctx, scriptID)
	if err != nil {
		return err
	}

	var panels []PanelData
	if script.ParsedData != "" {
		if err := json.Unmarshal([]byte(script.ParsedData), &panels); err != nil {
			return err
		}
	}

	if panelIndex < 0 || panelIndex >= len(panels) {
		return errors.New("panel index out of bounds")
	}

	panels[panelIndex] = panel

	out, err := json.Marshal(panels)
	if err != nil {
		return err
	}

	script.ParsedData = string(out)
	return s.repo.Update(ctx, script)
}

func (s *ScriptService) RegeneratePanel(ctx context.Context, scriptID string, panelIndex int, instructions string) error {
	script, err := s.repo.GetByID(ctx, scriptID)
	if err != nil {
		return err
	}

	var panels []PanelData
	if script.ParsedData == "" {
		return errors.New("no parsed data available to regenerate")
	}
	if err := json.Unmarshal([]byte(script.ParsedData), &panels); err != nil {
		return err
	}

	if panelIndex < 0 || panelIndex >= len(panels) {
		return errors.New("panel index out of bounds")
	}
	
	if panels[panelIndex].Locked {
		return errors.New("panel is locked")
	}

	targetPanel, _ := json.Marshal(panels[panelIndex])

	sysPrompt := "You are an AI revising a specific panel for a 4-panel comic. " +
		"I will provide the entire story context, the original panel data, and instructions for how to change it. " +
		"Return ONLY the revised JSON for this single panel in the exact format: " +
		`{"panel": X, "visual_desc": "...", "dialog": "...", "locked": false}` + "\n\n" +
		"Story Context: " + script.Content + "\n" +
		"Original Panel JSON: " + string(targetPanel) + "\n" +
		"Instructions for Revision: " + instructions

	res, err := s.llm.Generate(ctx, llm.TextRequest{Prompt: sysPrompt})
	if err != nil {
		return err
	}
	
	content := strings.TrimSpace(res.Content)
	if strings.HasPrefix(content, "```json") {
		content = strings.TrimPrefix(content, "```json")
		content = strings.TrimSuffix(content, "```")
		content = strings.TrimSpace(content)
	}

	var newPanel PanelData
	if err := json.Unmarshal([]byte(content), &newPanel); err != nil {
		return errors.New("failed to parse AI output into Panel JSON: " + err.Error())
	}
	
	newPanel.Panel = panels[panelIndex].Panel
	newPanel.Locked = false
	panels[panelIndex] = newPanel

	out, err := json.Marshal(panels)
	if err != nil {
		return err
	}

	script.ParsedData = string(out)
	return s.repo.Update(ctx, script)
}
