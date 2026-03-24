package handler

import (
	"net/http"

	"github.com/aimerneige/auto-you-koma/internal/model"
	"github.com/aimerneige/auto-you-koma/internal/repository"
	"github.com/aimerneige/auto-you-koma/internal/service"
	"github.com/gin-gonic/gin"
)

type CharacterHandler struct {
	svc *service.CharacterService
}

func NewCharacterHandler(svc *service.CharacterService) *CharacterHandler {
	return &CharacterHandler{svc: svc}
}

func (h *CharacterHandler) List(c *gin.Context) {
	q := c.Query("q")
	if q != "" {
		chars, err := h.svc.Search(c.Request.Context(), q)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, chars)
		return
	}

	filter := repository.CharacterFilter{
		Category: c.Query("category"),
	}
	
	if tags := c.QueryArray("tags"); len(tags) > 0 {
		filter.Tags = tags
	}

	chars, err := h.svc.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chars)
}

func (h *CharacterHandler) Create(c *gin.Context) {
	var char model.Character
	if err := c.ShouldBindJSON(&char); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	char.UserID = c.GetString("user_id")

	if err := h.svc.Create(c.Request.Context(), &char); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, char)
}

func (h *CharacterHandler) Get(c *gin.Context) {
	id := c.Param("id")
	char, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "character not found"})
		return
	}
	c.JSON(http.StatusOK, char)
}

func (h *CharacterHandler) UploadImage(c *gin.Context) {
	id := c.Param("id")
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file required"})
		return
	}
	defer file.Close()

	imageType := c.PostForm("image_type")
	desc := c.PostForm("description")
	isPrimary := c.PostForm("is_primary") == "true"
	
	var variantID *string
	if vid := c.PostForm("variant_id"); vid != "" {
		variantID = &vid
	}

	img, err := h.svc.UploadImage(c.Request.Context(), id, variantID, imageType, file, header.Filename, desc, isPrimary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, img)
}

func (h *CharacterHandler) AddVariant(c *gin.Context) {
	id := c.Param("id")
	var variant model.CharacterVariant
	if err := c.ShouldBindJSON(&variant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	variant.CharacterID = id

	if err := h.svc.AddVariant(c.Request.Context(), &variant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, variant)
}
