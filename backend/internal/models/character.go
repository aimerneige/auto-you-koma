package models

import (
	"time"

	"github.com/google/uuid"
)

type Character struct {
	ID               string    `gorm:"type:char(36);primaryKey" json:"id"`
	UserID           string    `gorm:"type:char(36);index;not null" json:"user_id"`
	Name             string    `gorm:"type:varchar(255);not null" json:"name"`
	NameJP           string    `gorm:"type:varchar(255)" json:"name_jp"`
	Gender           string    `gorm:"type:varchar(50)" json:"gender"`
	Age              string    `gorm:"type:varchar(50)" json:"age"`
	Personality      string    `gorm:"type:text" json:"personality"` // JSON array
	Backstory        string    `gorm:"type:text" json:"backstory"`
	VisualPrompt     string    `gorm:"type:text" json:"visual_prompt"`
	Tags             string    `gorm:"type:text" json:"tags"` // JSON array
	Category         string    `gorm:"type:varchar(100)" json:"category"`
	ReferenceSheetURL string   `gorm:"type:varchar(500)" json:"reference_sheet_url"` // Key field for reference image
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (c *Character) BeforeCreate() error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}