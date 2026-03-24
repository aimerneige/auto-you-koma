package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/aimerneige/auto-you-koma/internal/service"
)

// RenderHandler handles render API requests
type RenderHandler struct {
	service *service.RenderService
}

func NewRenderHandler(svc *service.RenderService) *RenderHandler {
	return &RenderHandler{service: svc}
}

// RenderRequest represents the request body for rendering
type RenderRequest struct {
	ExportType  string `json:"export_type"`
	Layout      string `json:"layout"`
	ImageWidth  int    `json:"image_width"`
	ImageHeight int    `json:"image_height"`
}

// StartRender handles POST /api/v1/projects/:id/render
func (h *RenderHandler) StartRender(c *gin.Context) {
	projectID := c.Param("id")

	var req RenderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Use defaults
		req = RenderRequest{
			ExportType:  "clean_plate",
			Layout:      "2x2",
			ImageWidth:  1024,
			ImageHeight: 1024,
		}
	}

	// Get project to find script ID
	_, err := h.service.GetProject(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	script, err := h.service.GetScript(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Script not found"})
		return
	}

	renderReq := service.RenderRequest{
		ProjectID:   projectID,
		ScriptID:    script.ID,
		ExportType:  req.ExportType,
		Layout:      req.Layout,
		ImageWidth:  req.ImageWidth,
		ImageHeight: req.ImageHeight,
	}

	task, err := h.service.StartRender(c.Request.Context(), renderReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rendering started",
		"task":    task,
	})
}

// GetRenderStatus handles GET /api/v1/projects/:id/render-status
func (h *RenderHandler) GetRenderStatus(c *gin.Context) {
	projectID := c.Param("id")

	task, err := h.service.GetRenderTask(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Render task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// RegeneratePanel handles POST /api/v1/projects/:id/render/regenerate
func (h *RenderHandler) RegeneratePanel(c *gin.Context) {
	projectID := c.Param("id")

	var req struct {
		PanelNumber int `json:"panel_number" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "panel_number is required"})
		return
	}

	task, err := h.service.GetRenderTask(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Render task not found"})
		return
	}

	result, err := h.service.RegenerateSinglePanel(c.Request.Context(), task.ID, req.PanelNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ConfirmRender handles POST /api/v1/projects/:id/render/confirm
func (h *RenderHandler) ConfirmRender(c *gin.Context) {
	projectID := c.Param("id")

	err := h.service.ConfirmRender(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Render confirmed"})
}

// RegisterRoutes registers render routes
func (h *RenderHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/projects/:id/render", h.StartRender)
	r.GET("/projects/:id/render-status", h.GetRenderStatus)
	r.POST("/projects/:id/render/regenerate", h.RegeneratePanel)
	r.POST("/projects/:id/render/confirm", h.ConfirmRender)
}