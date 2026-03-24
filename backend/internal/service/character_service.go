package service

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/aimerneige/auto-you-koma/internal/config"
	"github.com/aimerneige/auto-you-koma/internal/model"
	"github.com/aimerneige/auto-you-koma/internal/repository"
)

type CharacterService struct {
	repo repository.CharacterRepository
	cfg  config.StorageConfig
}

func NewCharacterService(repo repository.CharacterRepository, cfg config.StorageConfig) *CharacterService {
	return &CharacterService{repo: repo, cfg: cfg}
}

func (s *CharacterService) Create(ctx context.Context, char *model.Character) error {
	char.ID = uuid.New().String()
	return s.repo.Create(ctx, char)
}

func (s *CharacterService) GetByID(ctx context.Context, id string) (*model.Character, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CharacterService) List(ctx context.Context, filter repository.CharacterFilter) ([]*model.Character, error) {
	return s.repo.List(ctx, filter)
}

func (s *CharacterService) Update(ctx context.Context, char *model.Character) error {
	return s.repo.Update(ctx, char)
}

func (s *CharacterService) Delete(ctx context.Context, id string) error {
	// (Optional) Here we should ideally delete all physical files tied to the character images as well.
	return s.repo.Delete(ctx, id)
}

func (s *CharacterService) UploadImage(ctx context.Context, charID string, variantID *string, imageType string, file io.Reader, filename string, desc string, isPrimary bool) (*model.CharacterImage, error) {
	// Ensure directory exists
	targetDir := filepath.Join(s.cfg.BasePath, "characters", charID)
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return nil, err
	}

	ext := filepath.Ext(filename)
	newFileName := uuid.New().String() + ext
	filePath := filepath.Join(targetDir, newFileName)

	out, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return nil, err
	}

	// Relative path to store in DB
	relPath := "/assets/characters/" + charID + "/" + newFileName

	img := &model.CharacterImage{
		ID:          uuid.New().String(),
		CharacterID: charID,
		VariantID:   variantID,
		ImageType:   imageType,
		FilePath:    relPath,
		Description: desc,
		IsPrimary:   isPrimary,
	}

	if err := s.repo.AddImage(ctx, img); err != nil {
		return nil, err
	}
	return img, nil
}

func (s *CharacterService) AddVariant(ctx context.Context, v *model.CharacterVariant) error {
	v.ID = uuid.New().String()
	return s.repo.AddVariant(ctx, v)
}

func (s *CharacterService) Search(ctx context.Context, q string) ([]*model.Character, error) {
	return s.repo.Search(ctx, q)
}
