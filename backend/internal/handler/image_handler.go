package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/aimerneige/auto-you-koma/internal/storage"
)

// ImageHandler handles image upload and storage
type ImageHandler struct {
	storage *storage.Storage
}

func NewImageHandler(st *storage.Storage) *ImageHandler {
	return &ImageHandler{storage: st}
}

// UploadCharacterImage handles POST /api/v1/characters/:id/images
func (h *ImageHandler) UploadCharacterImage(c *gin.Context) {
	characterID := c.Param("id")

	// Get the file from the request
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}

	// Validate file type
	ext := filepath.Ext(file.Filename)
	validExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !validExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image file type"})
		return
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer src.Close()

	// Save the file
	relativePath, err := h.storage.SaveImageFromReader(characterID, file.Filename, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save image: %v", err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"path": relativePath,
		"url":  "/api/v1/images/" + relativePath,
	})
}

// ServeImage handles GET /api/v1/images/:path
func (h *ImageHandler) ServeImage(c *gin.Context) {
	pathParam := c.Param("path")
	fullPath := h.storage.GetImagePath(pathParam)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// Serve the file
	c.File(fullPath)
}

// RegisterRoutes registers image routes
func (h *ImageHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/characters/:id/images", h.UploadCharacterImage)
	r.GET("/images/*path", h.ServeImage)
}