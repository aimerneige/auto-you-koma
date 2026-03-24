package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/aimerneige/auto-you-koma/internal/service"
)

// ProjectHandler handles project API requests
type ProjectHandler struct {
	service *service.ProjectService
}

func NewProjectHandler(svc *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: svc}
}

// CreateProjectRequest represents the request body for creating a project
type CreateProjectRequest struct {
	Title    string `json:"title" binding:"required"`
	Mode     string `json:"mode"`
	Synopsis string `json:"synopsis" binding:"required"`
}

// Create handles POST /api/v1/projects
func (h *ProjectHandler) Create(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		userID = "default-user"
	}

	result, err := h.service.Create(c.Request.Context(), service.CreateProjectRequest{
		UserID:   userID,
		Title:    req.Title,
		Mode:     req.Mode,
		Synopsis: req.Synopsis,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// List handles GET /api/v1/projects
func (h *ProjectHandler) List(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		userID = "default-user"
	}

	projects, err := h.service.List(c.Request.Context(), userID, 50, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// Get handles GET /api/v1/projects/:id
func (h *ProjectHandler) Get(c *gin.Context) {
	id := c.Param("id")

	project, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

// Update handles PUT /api/v1/projects/:id
func (h *ProjectHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Update(c.Request.Context(), id, service.CreateProjectRequest{
		Title:    req.Title,
		Mode:     req.Mode,
		Synopsis: req.Synopsis,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Delete handles DELETE /api/v1/projects/:id
func (h *ProjectHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GenerateScript handles POST /api/v1/projects/:id/generate-script
func (h *ProjectHandler) GenerateScript(c *gin.Context) {
	id := c.Param("id")

	script, err := h.service.GenerateScript(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Script generated successfully",
		"script":  script,
	})
}

// GetScript handles GET /api/v1/projects/:id/script
func (h *ProjectHandler) GetScript(c *gin.Context) {
	id := c.Param("id")

	script, err := h.service.GetScript(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Script not found"})
		return
	}

	c.JSON(http.StatusOK, script)
}

// UpdateScript handles PUT /api/v1/projects/:id/script
func (h *ProjectHandler) UpdateScript(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	script, err := h.service.GetScript(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Script not found"})
		return
	}

	result, err := h.service.UpdateScript(c.Request.Context(), script.ID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// RegisterRoutes registers project routes
func (h *ProjectHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/projects", h.Create)
	r.GET("/projects", h.List)
	r.GET("/projects/:id", h.Get)
	r.PUT("/projects/:id", h.Update)
	r.DELETE("/projects/:id", h.Delete)

	// Script endpoints (HITL Node 2)
	r.POST("/projects/:id/generate-script", h.GenerateScript)
	r.GET("/projects/:id/script", h.GetScript)
	r.PUT("/projects/:id/script", h.UpdateScript)

	// Storyboard endpoints (HITL Node 3)
	r.POST("/projects/:id/generate-storyboard", h.GenerateStoryboard)
	r.GET("/projects/:id/storyboard", h.GetStoryboard)
	r.PUT("/projects/:id/storyboard", h.UpdateStoryboard)
}

// GenerateStoryboard handles POST /api/v1/projects/:id/generate-storyboard
func (h *ProjectHandler) GenerateStoryboard(c *gin.Context) {
	id := c.Param("id")

	storyboard, err := h.service.GenerateStoryboard(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Storyboard generated successfully",
		"storyboard": storyboard,
	})
}

// GetStoryboard handles GET /api/v1/projects/:id/storyboard
func (h *ProjectHandler) GetStoryboard(c *gin.Context) {
	id := c.Param("id")

	storyboard, err := h.service.GetStoryboard(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Storyboard not found"})
		return
	}

	c.JSON(http.StatusOK, storyboard)
}

// UpdateStoryboard handles PUT /api/v1/projects/:id/storyboard
func (h *ProjectHandler) UpdateStoryboard(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storyboard, err := h.service.GetStoryboard(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Storyboard not found"})
		return
	}

	result, err := h.service.UpdateStoryboard(c.Request.Context(), storyboard.ID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}