package repository

import (
	"context"
	"github.com/aimerneige/auto-you-koma/internal/model"
)

type CharacterFilter struct {
	Category string
	Tags     []string
}

type CharacterRepository interface {
	Create(ctx context.Context, character *model.Character) error
	GetByID(ctx context.Context, id string) (*model.Character, error)
	List(ctx context.Context, filter CharacterFilter) ([]*model.Character, error)
	Update(ctx context.Context, character *model.Character) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, query string) ([]*model.Character, error)
}

type ProjectRepository interface {
	Create(ctx context.Context, project *model.Project) error
	GetByID(ctx context.Context, id string) (*model.Project, error)
	List(ctx context.Context, userID string) ([]*model.Project, error)
	Update(ctx context.Context, project *model.Project) error
	Delete(ctx context.Context, id string) error
}

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}
