package sqlite

import (
	"context"

	"github.com/aimerneige/auto-you-koma/internal/model"
	"github.com/aimerneige/auto-you-koma/internal/repository"

	"gorm.io/gorm"
)

type generationRepo struct {
	db *gorm.DB
}

func NewGenerationRepository(db *gorm.DB) repository.GenerationRepository {
	return &generationRepo{db: db}
}

func (r *generationRepo) Create(ctx context.Context, generation *model.Generation) error {
	return r.db.WithContext(ctx).Create(generation).Error
}

func (r *generationRepo) GetByID(ctx context.Context, id string) (*model.Generation, error) {
	var g model.Generation
	if err := r.db.WithContext(ctx).First(&g, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *generationRepo) ListByScript(ctx context.Context, scriptID string) ([]*model.Generation, error) {
	var generations []*model.Generation
	if err := r.db.WithContext(ctx).Where("script_id = ?", scriptID).Order("created_at desc").Find(&generations).Error; err != nil {
		return nil, err
	}
	return generations, nil
}

func (r *generationRepo) Update(ctx context.Context, generation *model.Generation) error {
	return r.db.WithContext(ctx).Save(generation).Error
}

func (r *generationRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Generation{}, "id = ?", id).Error
}
