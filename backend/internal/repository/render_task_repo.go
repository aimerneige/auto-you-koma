package repository

import (
	"context"

	"github.com/aimerneige/auto-you-koma/internal/models"
	"gorm.io/gorm"
)

// RenderTaskRepository defines render task data access operations
type RenderTaskRepository interface {
	Create(ctx context.Context, task *models.RenderTask) error
	GetByID(ctx context.Context, id string) (*models.RenderTask, error)
	GetByProjectID(ctx context.Context, projectID string) (*models.RenderTask, error)
	Update(ctx context.Context, task *models.RenderTask) error
	Delete(ctx context.Context, id string) error
}

// SQLiteRenderTaskRepo implements RenderTaskRepository using SQLite
type SQLiteRenderTaskRepo struct {
	db *gorm.DB
}

func NewRenderTaskRepository(db *gorm.DB) *SQLiteRenderTaskRepo {
	return &SQLiteRenderTaskRepo{db: db}
}

func (r *SQLiteRenderTaskRepo) Create(ctx context.Context, task *models.RenderTask) error {
	return r.db.Create(task).Error
}

func (r *SQLiteRenderTaskRepo) GetByID(ctx context.Context, id string) (*models.RenderTask, error) {
	var task models.RenderTask
	err := r.db.Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *SQLiteRenderTaskRepo) GetByProjectID(ctx context.Context, projectID string) (*models.RenderTask, error) {
	var task models.RenderTask
	err := r.db.Where("project_id = ?", projectID).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *SQLiteRenderTaskRepo) Update(ctx context.Context, task *models.RenderTask) error {
	return r.db.Save(task).Error
}

func (r *SQLiteRenderTaskRepo) Delete(ctx context.Context, id string) error {
	return r.db.Where("id = ?", id).Delete(&models.RenderTask{}).Error
}