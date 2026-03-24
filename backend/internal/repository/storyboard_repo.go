package repository

import (
	"context"

	"github.com/aimerneige/auto-you-koma/internal/models"
	"gorm.io/gorm"
)

// StoryboardRepository defines storyboard data access operations
type StoryboardRepository interface {
	Create(ctx context.Context, storyboard *models.Storyboard) error
	GetByID(ctx context.Context, id string) (*models.Storyboard, error)
	GetByScriptID(ctx context.Context, scriptID string) (*models.Storyboard, error)
	Update(ctx context.Context, storyboard *models.Storyboard) error
	Delete(ctx context.Context, id string) error
}

// SQLiteStoryboardRepo implements StoryboardRepository using SQLite
type SQLiteStoryboardRepo struct {
	db *gorm.DB
}

func NewStoryboardRepository(db *gorm.DB) *SQLiteStoryboardRepo {
	return &SQLiteStoryboardRepo{db: db}
}

func (r *SQLiteStoryboardRepo) Create(ctx context.Context, storyboard *models.Storyboard) error {
	return r.db.Create(storyboard).Error
}

func (r *SQLiteStoryboardRepo) GetByID(ctx context.Context, id string) (*models.Storyboard, error) {
	var storyboard models.Storyboard
	err := r.db.Where("id = ?", id).First(&storyboard).Error
	if err != nil {
		return nil, err
	}
	return &storyboard, nil
}

func (r *SQLiteStoryboardRepo) GetByScriptID(ctx context.Context, scriptID string) (*models.Storyboard, error) {
	var storyboard models.Storyboard
	err := r.db.Where("script_id = ?", scriptID).First(&storyboard).Error
	if err != nil {
		return nil, err
	}
	return &storyboard, nil
}

func (r *SQLiteStoryboardRepo) Update(ctx context.Context, storyboard *models.Storyboard) error {
	return r.db.Save(storyboard).Error
}

func (r *SQLiteStoryboardRepo) Delete(ctx context.Context, id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Storyboard{}).Error
}