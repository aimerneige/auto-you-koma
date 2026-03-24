package handler

import (
	"net/http"

	"github.com/aimerneige/auto-you-koma/internal/model"
	"github.com/aimerneige/auto-you-koma/internal/service"
	"github.com/gin-gonic/gin"
)

type ScriptHandler struct {
	svc *service.ScriptService
}

func NewScriptHandler(svc *service.ScriptService) *ScriptHandler {
	return &ScriptHandler{svc: svc}
}

func (h *ScriptHandler) Create(c *gin.Context) {
	var script model.Script
	if err := c.ShouldBindJSON(&script); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.Create(c.Request.Context(), &script); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, script)
}

func (h *ScriptHandler) Get(c *gin.Context) {
	id := c.Param("id")
	script, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "script not found"})
		return
	}
	c.JSON(http.StatusOK, script)
}

func (h *ScriptHandler) ListByProject(c *gin.Context) {
	projectID := c.Query("project_id")
	if projectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "project_id is required"})
		return
	}

	scripts, err := h.svc.ListByProject(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, scripts)
}

func (h *ScriptHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var script model.Script
	if err := c.ShouldBindJSON(&script); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	script.ID = id
	if err := h.svc.Update(c.Request.Context(), &script); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, script)
}

type ReqGenerate struct {
	Prompt string `json:"prompt" binding:"required"`
}

func (h *ScriptHandler) GenerateStream(c *gin.Context) {
	// Simple wrapped blocking call for now; robust SSE can be implemented later
	var req ReqGenerate
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	text, err := h.svc.GenerateScript(c.Request.Context(), req.Prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"content": text})
}

func (h *ScriptHandler) Parse(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.ParseToPanels(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	script, _ := h.svc.GetByID(c.Request.Context(), id)
	c.JSON(http.StatusOK, script)
}

func (h *ScriptHandler) UpdatePanel(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		PanelIndex int              `json:"panel_index" binding:"required"`
		PanelData  service.PanelData `json:"panel_data" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.UpdatePanel(c.Request.Context(), id, req.PanelIndex, req.PanelData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ScriptHandler) RegeneratePanel(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		PanelIndex   int    `json:"panel_index" binding:"required"`
		Instructions string `json:"instructions" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.RegeneratePanel(c.Request.Context(), id, req.PanelIndex, req.Instructions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return updated full script directly
	script, _ := h.svc.GetByID(c.Request.Context(), id)
	c.JSON(http.StatusOK, script)
}
