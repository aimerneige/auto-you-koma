package models

import (
	"time"

	"github.com/google/uuid"
)

type Script struct {
	ID            string    `gorm:"type:char(36);primaryKey" json:"id"`
	ProjectID     string    `gorm:"type:char(36);index;not null" json:"project_id"`
	EpisodeNumber int       `gorm:"default:1" json:"episode_number"`
	Content       string    `gorm:"type:text" json:"content"` // JSONB
	Version       int       `gorm:"default:1" json:"version"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (s *Script) BeforeCreate() error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}