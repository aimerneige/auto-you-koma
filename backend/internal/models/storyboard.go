package models

import (
	"time"

	"github.com/google/uuid"
)

type Storyboard struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	ScriptID  string    `gorm:"type:char(36);index;not null" json:"script_id"`
	Content   string    `gorm:"type:text" json:"content"` // JSONB
	Version   int       `gorm:"default:1" json:"version"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Storyboard) BeforeCreate() error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}