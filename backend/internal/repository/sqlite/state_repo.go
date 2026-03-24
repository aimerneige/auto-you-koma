package sqlite

import (
	"context"

	"github.com/aimerneige/auto-you-koma/internal/model"
	"github.com/aimerneige/auto-you-koma/internal/repository"

	"gorm.io/gorm"
)

type stateRepo struct {
	db *gorm.DB
}

func NewCharacterStateRepository(db *gorm.DB) repository.CharacterStateRepository {
	return &stateRepo{db: db}
}

func (r *stateRepo) Save(ctx context.Context, state *model.CharacterState) error {
	// Upsert logic if state already exists for this SeriesID + CharacterID, but for simplicity we rely on manual ID checks
	return r.db.WithContext(ctx).Save(state).Error
}

func (r *stateRepo) Get(ctx context.Context, seriesID, characterID string) (*model.CharacterState, error) {
	var s model.CharacterState
	if err := r.db.WithContext(ctx).First(&s, "series_id = ? AND character_id = ?", seriesID, characterID).Error; err != nil {
		return nil, err
	}
	return &s, nil
}
