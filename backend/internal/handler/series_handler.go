package handler

import (
	"net/http"

	"github.com/aimerneige/auto-you-koma/internal/model"
	"github.com/aimerneige/auto-you-koma/internal/repository"
	"github.com/aimerneige/auto-you-koma/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SeriesHandler struct {
	seriesRepo repository.SeriesRepository
	contSvc    *service.ContinuityService
}

func NewSeriesHandler(sr repository.SeriesRepository, c *service.ContinuityService) *SeriesHandler {
	return &SeriesHandler{seriesRepo: sr, contSvc: c}
}

type ReqCreateSeries struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

func (h *SeriesHandler) Create(c *gin.Context) {
	userId, _ := c.Get("user_id")

	var req ReqCreateSeries
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	s := &model.Series{
		ID:          uuid.New().String(),
		UserID:      userId.(string),
		Title:       req.Title,
		Description: req.Description,
	}

	if err := h.seriesRepo.Create(c.Request.Context(), s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, s)
}

func (h *SeriesHandler) List(c *gin.Context) {
	userId, _ := c.Get("user_id")
	list, err := h.seriesRepo.ListByUser(c.Request.Context(), userId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// Kick off continuity generation from last script
func (h *SeriesHandler) SythesizeMemory(c *gin.Context) {
	seriesId := c.Param("id")
	scriptId := c.Param("scriptId")
	characterId := c.Param("characterId")

	err := h.contSvc.SynthesizeMemory(c.Request.Context(), seriesId, scriptId, characterId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
