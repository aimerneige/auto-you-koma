package models

import (
	"time"

	"github.com/google/uuid"
)

type Panel struct {
	ID             string    `gorm:"type:char(36);primaryKey" json:"id"`
	StoryboardID   string    `gorm:"type:char(36);index;not null" json:"storyboard_id"`
	PanelNumber    int       `gorm:"not null" json:"panel_number"`
	ImageURL       string    `gorm:"type:varchar(500)" json:"image_url"`
	Dialogue       string    `gorm:"type:text" json:"dialogue"`
	Narration      string    `gorm:"type:text" json:"narration"`
	PanelPrompt    string    `gorm:"type:text" json:"panel_prompt"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (p *Panel) BeforeCreate() error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}