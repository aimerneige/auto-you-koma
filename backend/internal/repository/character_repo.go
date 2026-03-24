package repository

import (
	"context"

	"github.com/aimerneige/auto-you-koma/internal/models"
	"gorm.io/gorm"
)

// CharacterRepository defines character data access operations
type CharacterRepository interface {
	Create(ctx context.Context, character *models.Character) error
	GetByID(ctx context.Context, id string) (*models.Character, error)
	List(ctx context.Context, userID string, limit, offset int) ([]*models.Character, error)
	Update(ctx context.Context, character *models.Character) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, userID string, query string) ([]*models.Character, error)
}

// SQLiteCharacterRepo implements CharacterRepository using SQLite
type SQLiteCharacterRepo struct {
	db *gorm.DB
}

func NewCharacterRepository(db *gorm.DB) *SQLiteCharacterRepo {
	return &SQLiteCharacterRepo{db: db}
}

func (r *SQLiteCharacterRepo) Create(ctx context.Context, character *models.Character) error {
	return r.db.Create(character).Error
}

func (r *SQLiteCharacterRepo) GetByID(ctx context.Context, id string) (*models.Character, error) {
	var character models.Character
	err := r.db.Where("id = ?", id).First(&character).Error
	if err != nil {
		return nil, err
	}
	return &character, nil
}

func (r *SQLiteCharacterRepo) List(ctx context.Context, userID string, limit, offset int) ([]*models.Character, error) {
	var characters []*models.Character
	query := r.db.Where("user_id = ?", userID)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Find(&characters).Error
	return characters, err
}

func (r *SQLiteCharacterRepo) Update(ctx context.Context, character *models.Character) error {
	return r.db.Save(character).Error
}

func (r *SQLiteCharacterRepo) Delete(ctx context.Context, id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Character{}).Error
}

func (r *SQLiteCharacterRepo) Search(ctx context.Context, userID string, query string) ([]*models.Character, error) {
	var characters []*models.Character
	err := r.db.Where("user_id = ? AND (name LIKE ? OR name_jp LIKE ? OR tags LIKE ?)",
		userID, "%"+query+"%", "%"+query+"%", "%"+query+"%").Find(&characters).Error
	return characters, err
}