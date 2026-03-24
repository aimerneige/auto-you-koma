package sqlite

import (
	"context"

	"github.com/aimerneige/auto-you-koma/internal/model"
	"github.com/aimerneige/auto-you-koma/internal/repository"

	"gorm.io/gorm"
)

type scriptRepo struct {
	db *gorm.DB
}

func NewScriptRepository(db *gorm.DB) repository.ScriptRepository {
	return &scriptRepo{db: db}
}

func (r *scriptRepo) Create(ctx context.Context, script *model.Script) error {
	return r.db.WithContext(ctx).Create(script).Error
}

func (r *scriptRepo) GetByID(ctx context.Context, id string) (*model.Script, error) {
	var s model.Script
	if err := r.db.WithContext(ctx).First(&s, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *scriptRepo) ListByProject(ctx context.Context, projectID string) ([]*model.Script, error) {
	var scripts []*model.Script
	if err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Find(&scripts).Error; err != nil {
		return nil, err
	}
	return scripts, nil
}

func (r *scriptRepo) Update(ctx context.Context, script *model.Script) error {
	return r.db.WithContext(ctx).Save(script).Error
}

func (r *scriptRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Script{}, "id = ?", id).Error
}
