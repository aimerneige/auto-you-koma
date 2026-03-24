package service

import (
	"context"
	"fmt"

	"github.com/aimerneige/auto-you-koma/internal/llm"
	"github.com/aimerneige/auto-you-koma/internal/models"
	"github.com/aimerneige/auto-you-koma/internal/repository"
)

// CharacterService provides character business logic
type CharacterService struct {
	repo           *repository.SQLiteCharacterRepo
	imageGenerator llm.ImageGenerator
}

func NewCharacterService(repo *repository.SQLiteCharacterRepo, imgGen llm.ImageGenerator) *CharacterService {
	return &CharacterService{
		repo:           repo,
		imageGenerator: imgGen,
	}
}

type CreateCharacterRequest struct {
	UserID       string
	Name         string
	NameJP       string
	Gender       string
	Age          string
	Personality  string
	Backstory    string
	VisualPrompt string
	Tags         string
	Category     string
}

func (s *CharacterService) Create(ctx context.Context, req CreateCharacterRequest) (*models.Character, error) {
	character := &models.Character{
		UserID:       req.UserID,
		Name:         req.Name,
		NameJP:       req.NameJP,
		Gender:       req.Gender,
		Age:          req.Age,
		Personality:  req.Personality,
		Backstory:    req.Backstory,
		VisualPrompt: req.VisualPrompt,
		Tags:         req.Tags,
		Category:     req.Category,
	}
	err := s.repo.Create(ctx, character)
	return character, err
}

func (s *CharacterService) GetByID(ctx context.Context, id string) (*models.Character, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CharacterService) List(ctx context.Context, userID string, limit, offset int) ([]*models.Character, error) {
	return s.repo.List(ctx, userID, limit, offset)
}

func (s *CharacterService) Update(ctx context.Context, id string, req CreateCharacterRequest) (*models.Character, error) {
	character, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	character.Name = req.Name
	character.NameJP = req.NameJP
	character.Gender = req.Gender
	character.Age = req.Age
	character.Personality = req.Personality
	character.Backstory = req.Backstory
	character.VisualPrompt = req.VisualPrompt
	character.Tags = req.Tags
	character.Category = req.Category
	err = s.repo.Update(ctx, character)
	return character, err
}

func (s *CharacterService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *CharacterService) Search(ctx context.Context, userID string, query string) ([]*models.Character, error) {
	return s.repo.Search(ctx, userID, query)
}

// GenerateReferenceSheet generates a reference sheet for a character using AI image generation
func (s *CharacterService) GenerateReferenceSheet(ctx context.Context, characterID string) (string, error) {
	character, err := s.repo.GetByID(ctx, characterID)
	if err != nil {
		return "", err
	}

	// Build the prompt for reference sheet generation
	prompt := fmt.Sprintf(
		"Character reference sheet: %s, %s, %s. Style: anime, clean lineart, front view, side view, back view three poses, white background.",
		character.Name,
		character.VisualPrompt,
		character.Personality,
	)

	req := llm.ImageRequest{
		Prompt:         prompt,
		NegativePrompt: "blurry, low quality, distorted, deformed",
		Width:          1024,
		Height:         1024,
		Seed:           42,
	}

	resp, err := s.imageGenerator.GenerateImage(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.ImageURL, nil
}

// SetReferenceSheetURL saves the confirmed reference sheet URL to the character
func (s *CharacterService) SetReferenceSheetURL(ctx context.Context, characterID string, url string) error {
	character, err := s.repo.GetByID(ctx, characterID)
	if err != nil {
		return err
	}
	character.ReferenceSheetURL = url
	return s.repo.Update(ctx, character)
}