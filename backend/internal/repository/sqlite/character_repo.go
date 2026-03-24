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
	if err := r.db.WithContext(ctx).Preload("Images").Preload("Variants").First(&c, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *characterRepo) List(ctx context.Context, filter repository.CharacterFilter) ([]*model.Character, error) {
	query := r.db.WithContext(ctx).Preload("Images")
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
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("character_id = ?", id).Delete(&model.CharacterVariant{}).Error; err != nil {
			return err
		}
		if err := tx.Where("character_id = ?", id).Delete(&model.CharacterImage{}).Error; err != nil {
			return err
		}
		if err := tx.Where("character_id = ?", id).Delete(&model.CharacterGroupMember{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Character{}, "id = ?", id).Error
	})
}

func (r *characterRepo) Search(ctx context.Context, q string) ([]*model.Character, error) {
	var chars []*model.Character
	query := "%" + q + "%"
	if err := r.db.WithContext(ctx).Preload("Images").Where("name LIKE ? OR name_jp LIKE ? OR tags LIKE ?", query, query, query).Find(&chars).Error; err != nil {
		return nil, err
	}
	return chars, nil
}

func (r *characterRepo) AddVariant(ctx context.Context, variant *model.CharacterVariant) error {
	return r.db.WithContext(ctx).Create(variant).Error
}

func (r *characterRepo) UpdateVariant(ctx context.Context, variant *model.CharacterVariant) error {
	return r.db.WithContext(ctx).Save(variant).Error
}

func (r *characterRepo) DeleteVariant(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.CharacterVariant{}, "id = ?", id).Error
}

func (r *characterRepo) AddImage(ctx context.Context, image *model.CharacterImage) error {
	// If primary, unset others for this character
	if image.IsPrimary {
		r.db.WithContext(ctx).Model(&model.CharacterImage{}).Where("character_id = ?", image.CharacterID).Update("is_primary", false)
	}
	return r.db.WithContext(ctx).Create(image).Error
}

func (r *characterRepo) DeleteImage(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.CharacterImage{}, "id = ?", id).Error
}

// Group methods
func (r *characterRepo) CreateGroup(ctx context.Context, group *model.CharacterGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *characterRepo) ListGroups(ctx context.Context, userID string) ([]*model.CharacterGroup, error) {
	var groups []*model.CharacterGroup
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

func (r *characterRepo) AddCharacterToGroup(ctx context.Context, member *model.CharacterGroupMember) error {
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *characterRepo) RemoveCharacterFromGroup(ctx context.Context, groupID, characterID string) error {
	return r.db.WithContext(ctx).Where("group_id = ? AND character_id = ?", groupID, characterID).Delete(&model.CharacterGroupMember{}).Error
}
