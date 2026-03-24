package repository

import (
	"context"

	"github.com/aimerneige/auto-you-koma/internal/models"
	"gorm.io/gorm"
)

// ScriptRepository defines script data access operations
type ScriptRepository interface {
	Create(ctx context.Context, script *models.Script) error
	GetByID(ctx context.Context, id string) (*models.Script, error)
	GetByProjectID(ctx context.Context, projectID string) (*models.Script, error)
	Update(ctx context.Context, script *models.Script) error
	Delete(ctx context.Context, id string) error
}

// SQLiteScriptRepo implements ScriptRepository using SQLite
type SQLiteScriptRepo struct {
	db *gorm.DB
}

func NewScriptRepository(db *gorm.DB) *SQLiteScriptRepo {
	return &SQLiteScriptRepo{db: db}
}

func (r *SQLiteScriptRepo) Create(ctx context.Context, script *models.Script) error {
	return r.db.Create(script).Error
}

func (r *SQLiteScriptRepo) GetByID(ctx context.Context, id string) (*models.Script, error) {
	var script models.Script
	err := r.db.Where("id = ?", id).First(&script).Error
	if err != nil {
		return nil, err
	}
	return &script, nil
}

func (r *SQLiteScriptRepo) GetByProjectID(ctx context.Context, projectID string) (*models.Script, error) {
	var script models.Script
	err := r.db.Where("project_id = ?", projectID).First(&script).Error
	if err != nil {
		return nil, err
	}
	return &script, nil
}

func (r *SQLiteScriptRepo) Update(ctx context.Context, script *models.Script) error {
	return r.db.Save(script).Error
}

func (r *SQLiteScriptRepo) Delete(ctx context.Context, id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Script{}).Error
}