package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/aimerneige/auto-you-koma/internal/service"
)

// CharacterHandler handles character API requests
type CharacterHandler struct {
	service *service.CharacterService
}

func NewCharacterHandler(svc *service.CharacterService) *CharacterHandler {
	return &CharacterHandler{service: svc}
}

// CreateCharacterRequest represents the request body for creating a character
type CreateCharacterRequest struct {
	Name         string `json:"name" binding:"required"`
	NameJP       string `json:"name_jp"`
	Gender       string `json:"gender"`
	Age          string `json:"age"`
	Personality  string `json:"personality"`
	Backstory    string `json:"backstory"`
	VisualPrompt string `json:"visual_prompt"`
	Tags         string `json:"tags"`
	Category     string `json:"category"`
}

// Create handles POST /api/v1/characters
func (h *CharacterHandler) Create(c *gin.Context) {
	var req CreateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context (would normally come from auth middleware)
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		userID = "default-user" // For demo purposes
	}

	result, err := h.service.Create(c.Request.Context(), service.CreateCharacterRequest{
		UserID:       userID,
		Name:         req.Name,
		NameJP:       req.NameJP,
		Gender:       req.Gender,
		Age:          req.Age,
		Personality:  req.Personality,
		Backstory:    req.Backstory,
		VisualPrompt: req.VisualPrompt,
		Tags:         req.Tags,
		Category:     req.Category,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// List handles GET /api/v1/characters
func (h *CharacterHandler) List(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		userID = "default-user"
	}

	characters, err := h.service.List(c.Request.Context(), userID, 50, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, characters)
}

// Get handles GET /api/v1/characters/:id
func (h *CharacterHandler) Get(c *gin.Context) {
	id := c.Param("id")

	character, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Character not found"})
		return
	}

	c.JSON(http.StatusOK, character)
}

// Update handles PUT /api/v1/characters/:id
func (h *CharacterHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req CreateCharacterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.Update(c.Request.Context(), id, service.CreateCharacterRequest{
		Name:         req.Name,
		NameJP:       req.NameJP,
		Gender:       req.Gender,
		Age:          req.Age,
		Personality:  req.Personality,
		Backstory:    req.Backstory,
		VisualPrompt: req.VisualPrompt,
		Tags:         req.Tags,
		Category:     req.Category,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// Delete handles DELETE /api/v1/characters/:id
func (h *CharacterHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// RegisterRoutes registers character routes
func (h *CharacterHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/characters", h.Create)
	r.GET("/characters", h.List)
	r.GET("/characters/:id", h.Get)
	r.PUT("/characters/:id", h.Update)
	r.DELETE("/characters/:id", h.Delete)
}