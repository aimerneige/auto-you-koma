package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

const (
	CharacterAssetsDir = "data/assets/characters"
)

// Storage handles file storage for images and assets
type Storage struct {
	basePath string
}

// NewStorage creates a new Storage instance
func NewStorage(basePath string) *Storage {
	if basePath == "" {
		basePath = "."
	}
	return &Storage{basePath: basePath}
}

// SaveCharacterImage saves an image file for a character
func (s *Storage) SaveCharacterImage(characterID string, filename string, data []byte) (string, error) {
	charDir := filepath.Join(s.basePath, CharacterAssetsDir, characterID)
	if err := os.MkdirAll(charDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create character directory: %w", err)
	}

	// Generate unique filename
	ext := filepath.Ext(filename)
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(charDir, newFilename)

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Return relative path for storage in database
	relativePath := filepath.Join(CharacterAssetsDir, characterID, newFilename)
	return relativePath, nil
}

// SaveImageFromReader saves an image from an io.Reader
func (s *Storage) SaveImageFromReader(characterID string, filename string, reader io.Reader) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read data: %w", err)
	}
	return s.SaveCharacterImage(characterID, filename, data)
}

// GetImagePath returns the full path for an image
func (s *Storage) GetImagePath(relativePath string) string {
	return filepath.Join(s.basePath, relativePath)
}

// DeleteImage deletes an image file
func (s *Storage) DeleteImage(relativePath string) error {
	filePath := filepath.Join(s.basePath, relativePath)
	return os.Remove(filePath)
}

// ImageExists checks if an image file exists
func (s *Storage) ImageExists(relativePath string) bool {
	filePath := filepath.Join(s.basePath, relativePath)
	_, err := os.Stat(filePath)
	return err == nil
}