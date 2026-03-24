package service

import (
	"context"
	"encoding/json"

	"github.com/aimerneige/auto-you-koma/internal/repository"
)

// CompositorService handles final image compositing and export
type CompositorService struct {
	storyboardRepo *repository.SQLiteStoryboardRepo
	scriptRepo     *repository.SQLiteScriptRepo
	renderTaskRepo *repository.SQLiteRenderTaskRepo
	projectRepo    *repository.SQLiteProjectRepo
}

// NewCompositorService creates a new CompositorService
func NewCompositorService(
	storyboardRepo *repository.SQLiteStoryboardRepo,
	scriptRepo *repository.SQLiteScriptRepo,
	renderTaskRepo *repository.SQLiteRenderTaskRepo,
	projectRepo *repository.SQLiteProjectRepo,
) *CompositorService {
	return &CompositorService{
		storyboardRepo: storyboardRepo,
		scriptRepo:     scriptRepo,
		renderTaskRepo: renderTaskRepo,
		projectRepo:    projectRepo,
	}
}

// TextBubble represents a text bubble position
type TextBubble struct {
	Character  string  `json:"character"`
	Text       string  `json:"text"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	Width      float64 `json:"width"`
	Height     float64 `json:"height"`
	BubbleType string  `json:"bubble_type"` // speech, thought, shout
}

// PanelTextCoords represents text coordinates for a panel
type PanelTextCoords struct {
	PanelNumber int          `json:"panel_number"`
	Bubbles     []TextBubble `json:"bubbles"`
}

// GenerateTextCoords generates text bubble coordinates for each panel
func (s *CompositorService) GenerateTextCoords(ctx context.Context, projectID string) ([]PanelTextCoords, error) {
	// Get script
	script, err := s.scriptRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var scriptContent ScriptContent
	if err := json.Unmarshal([]byte(script.Content), &scriptContent); err != nil {
		return nil, err
	}

	// Generate coordinates for each panel
	coords := make([]PanelTextCoords, len(scriptContent.Panels))
	for i, panel := range scriptContent.Panels {
		bubbles := generateBubblesForPanel(panel, i)
		coords[i] = PanelTextCoords{
			PanelNumber: panel.PanelNumber,
			Bubbles:     bubbles,
		}
	}

	return coords, nil
}

func generateBubblesForPanel(panel ScriptPanel, panelIndex int) []TextBubble {
	bubbles := []TextBubble{}

	// Simple coordinate logic based on panel position
	// Panel 1: left side, Panel 2: right side, etc.
	baseX := []float64{0.1, 0.5, 0.1, 0.5}
	baseY := float64(0.7)

	if panel.Dialogue != "" {
		bubbles = append(bubbles, TextBubble{
			Character:  panel.Characters,
			Text:       panel.Dialogue,
			X:          baseX[panelIndex],
			Y:          baseY,
			Width:      0.35,
			Height:     0.15,
			BubbleType: "speech",
		})
	}

	if panel.Narration != "" {
		bubbles = append(bubbles, TextBubble{
			Character:  "",
			Text:       panel.Narration,
			X:          0.3,
			Y:          0.1,
			Width:      0.4,
			Height:     0.1,
			BubbleType: "narration",
		})
	}

	return bubbles
}

// CompositorRequest contains the final compositing parameters
type CompositorRequest struct {
	ProjectID     string            `json:"project_id"`
	ExportType    string            `json:"export_type"` // native_text / clean_plate
	Layout        string            `json:"layout"`      // 2x2 / 1x4
	TextCoords    []PanelTextCoords `json:"text_coords"`
}

// CompositeResult contains the final composited image paths
type CompositeResult struct {
	FinalImageURL string `json:"final_image_url"`
	Layout        string `json:"layout"`
	PanelImages   []string `json:"panel_images"`
}

// Composite creates the final composited image
func (s *CompositorService) Composite(ctx context.Context, req CompositorRequest) (*CompositeResult, error) {
	// Get render task to get the panel images
	task, err := s.renderTaskRepo.GetByProjectID(ctx, req.ProjectID)
	if err != nil {
		return nil, err
	}

	// Parse output paths
	var panelImages []string
	if task.OutputPaths != "" {
		var results []PanelRenderResult
		json.Unmarshal([]byte(task.OutputPaths), &results)
		for _, r := range results {
			panelImages = append(panelImages, r.ImageURL)
		}
	}

	// For now, return the panel images as the result
	// In a real implementation, this would:
	// 1. Load each panel image
	// 2. Add text bubbles if export_type is native_text
	// 3. Composite into final layout (2x2 or 1x4)
	// 4. Save the final image

	result := &CompositeResult{
		Layout:        req.Layout,
		PanelImages:   panelImages,
		FinalImageURL: panelImages[0], // Placeholder - would be final composited image
	}

	return result, nil
}

// ExportProject exports the final project as a downloadable file
func (s *CompositorService) ExportProject(ctx context.Context, projectID string) (map[string]interface{}, error) {
	// Get project
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Get render task
	task, err := s.renderTaskRepo.GetByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	// Get text coordinates
	textCoords, _ := s.GenerateTextCoords(ctx, projectID)

	exportData := map[string]interface{}{
		"project":    project,
		"render_task": task,
		"text_coords": textCoords,
		"layout":      task.Layout,
		"export_type": task.ExportType,
	}

	return exportData, nil
}