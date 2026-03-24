package sqlite

import (
	"context"

	"github.com/aimerneige/auto-you-koma/internal/model"
	"github.com/aimerneige/auto-you-koma/internal/repository"

	"gorm.io/gorm"
)

type seriesRepo struct {
	db *gorm.DB
}

func NewSeriesRepository(db *gorm.DB) repository.SeriesRepository {
	return &seriesRepo{db: db}
}

func (r *seriesRepo) Create(ctx context.Context, series *model.Series) error {
	return r.db.WithContext(ctx).Create(series).Error
}

func (r *seriesRepo) GetByID(ctx context.Context, id string) (*model.Series, error) {
	var s model.Series
	if err := r.db.WithContext(ctx).First(&s, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *seriesRepo) ListByUser(ctx context.Context, userID string) ([]*model.Series, error) {
	var seriesList []*model.Series
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&seriesList).Error; err != nil {
		return nil, err
	}
	return seriesList, nil
}

func (r *seriesRepo) Update(ctx context.Context, series *model.Series) error {
	return r.db.WithContext(ctx).Save(series).Error
}

func (r *seriesRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Series{}, "id = ?", id).Error
}
