package handler

import (
	"net/http"

	"github.com/aimerneige/auto-you-koma/internal/service"
	"github.com/gin-gonic/gin"
)

type GenerationHandler struct {
	svc *service.GenerationService
}

func NewGenerationHandler(svc *service.GenerationService) *GenerationHandler {
	return &GenerationHandler{svc: svc}
}

type ReqStartGeneration struct {
	ProjectID string `json:"project_id" binding:"required"`
	ScriptID  string `json:"script_id" binding:"required"`
	Layout    string `json:"layout" binding:"required"` // '2x2' or '1x4'
}

func (h *GenerationHandler) Start(c *gin.Context) {
	var req ReqStartGeneration
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gen, err := h.svc.StartGeneration(c.Request.Context(), req.ProjectID, req.ScriptID, req.Layout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gen)
}

func (h *GenerationHandler) Get(c *gin.Context) {
	id := c.Param("id")
	gen, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "generation task not found"})
		return
	}
	c.JSON(http.StatusOK, gen)
}
