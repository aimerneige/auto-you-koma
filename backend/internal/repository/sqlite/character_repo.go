package sqlite

import (
	"context"
	"github.com/aimerneige/auto-you-koma/internal/model"
	"github.com/aimerneige/auto-you-koma/internal/repository"

	"gorm.io/gorm"
)

type characterRepo struct {
	db *gorm.DB
}

func NewCharacterRepository(db *gorm.DB) repository.CharacterRepository {
	return &characterRepo{db: db}
}

func (r *characterRepo) Create(ctx context.Context, character *model.Character) error {
	return r.db.WithContext(ctx).Create(character).Error
}

func (r *characterRepo) GetByID(ctx context.Context, id string) (*model.Character, error) {
	var c model.Character
	if err := r.db.WithContext(ctx).First(&c, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *characterRepo) List(ctx context.Context, filter repository.CharacterFilter) ([]*model.Character, error) {
	query := r.db.WithContext(ctx)
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	
	for _, tag := range filter.Tags {
		query = query.Where("tags LIKE ?", "%"+tag+"%")
	}

	var chars []*model.Character
	if err := query.Find(&chars).Error; err != nil {
		return nil, err
	}
	return chars, nil
}

func (r *characterRepo) Update(ctx context.Context, character *model.Character) error {
	return r.db.WithContext(ctx).Save(character).Error
}

func (r *characterRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Character{}, "id = ?", id).Error
}

func (r *characterRepo) Search(ctx context.Context, q string) ([]*model.Character, error) {
	var chars []*model.Character
	query := "%" + q + "%"
	if err := r.db.WithContext(ctx).Where("name LIKE ? OR name_jp LIKE ? OR tags LIKE ?", query, query, query).Find(&chars).Error; err != nil {
		return nil, err
	}
	return chars, nil
}
