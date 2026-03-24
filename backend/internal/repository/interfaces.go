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

	AddVariant(ctx context.Context, variant *model.CharacterVariant) error
	UpdateVariant(ctx context.Context, variant *model.CharacterVariant) error
	DeleteVariant(ctx context.Context, id string) error

	AddImage(ctx context.Context, image *model.CharacterImage) error
	DeleteImage(ctx context.Context, id string) error

	CreateGroup(ctx context.Context, group *model.CharacterGroup) error
	ListGroups(ctx context.Context, userID string) ([]*model.CharacterGroup, error)
	AddCharacterToGroup(ctx context.Context, member *model.CharacterGroupMember) error
	RemoveCharacterFromGroup(ctx context.Context, groupID, characterID string) error
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
