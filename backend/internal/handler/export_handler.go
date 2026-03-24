package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/aimerneige/auto-you-koma/internal/service"
)

// ExportHandler handles export API requests
type ExportHandler struct {
	service *service.CompositorService
}

func NewExportHandler(svc *service.CompositorService) *ExportHandler {
	return &ExportHandler{service: svc}
}

// GetTextCoords handles GET /api/v1/projects/:id/text-coords
func (h *ExportHandler) GetTextCoords(c *gin.Context) {
	projectID := c.Param("id")

	coords, err := h.service.GenerateTextCoords(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coords)
}

// Composite handles POST /api/v1/projects/:id/composite
func (h *ExportHandler) Composite(c *gin.Context) {
	projectID := c.Param("id")

	var req struct {
		ExportType string `json:"export_type"`
		Layout     string `json:"layout"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		req.ExportType = "clean_plate"
		req.Layout = "2x2"
	}

	// Get text coordinates
	coords, err := h.service.GenerateTextCoords(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	compositeReq := service.CompositorRequest{
		ProjectID:  projectID,
		ExportType: req.ExportType,
		Layout:     req.Layout,
		TextCoords: coords,
	}

	result, err := h.service.Composite(c.Request.Context(), compositeReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ExportProject handles GET /api/v1/projects/:id/export
func (h *ExportHandler) ExportProject(c *gin.Context) {
	projectID := c.Param("id")

	data, err := h.service.ExportProject(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

// RegisterRoutes registers export routes
func (h *ExportHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/projects/:id/text-coords", h.GetTextCoords)
	r.POST("/projects/:id/composite", h.Composite)
	r.GET("/projects/:id/export", h.ExportProject)
}