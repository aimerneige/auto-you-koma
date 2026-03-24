package repository

import (
	"context"

	"github.com/aimerneige/auto-you-koma/internal/models"
	"gorm.io/gorm"
)

// ProjectRepository defines project data access operations
type ProjectRepository interface {
	Create(ctx context.Context, project *models.Project) error
	GetByID(ctx context.Context, id string) (*models.Project, error)
	List(ctx context.Context, userID string, limit, offset int) ([]*models.Project, error)
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, id string) error
}

// SQLiteProjectRepo implements ProjectRepository using SQLite
type SQLiteProjectRepo struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *SQLiteProjectRepo {
	return &SQLiteProjectRepo{db: db}
}

func (r *SQLiteProjectRepo) Create(ctx context.Context, project *models.Project) error {
	return r.db.Create(project).Error
}

func (r *SQLiteProjectRepo) GetByID(ctx context.Context, id string) (*models.Project, error) {
	var project models.Project
	err := r.db.Where("id = ?", id).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *SQLiteProjectRepo) List(ctx context.Context, userID string, limit, offset int) ([]*models.Project, error) {
	var projects []*models.Project
	query := r.db.Where("user_id = ?", userID)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	err := query.Find(&projects).Error
	return projects, err
}

func (r *SQLiteProjectRepo) Update(ctx context.Context, project *models.Project) error {
	return r.db.Save(project).Error
}

func (r *SQLiteProjectRepo) Delete(ctx context.Context, id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Project{}).Error
}