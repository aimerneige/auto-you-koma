package model

import "time"

type Character struct {
	ID           string    `json:"id" gorm:"primaryKey"`
	UserID       string    `json:"user_id" gorm:"index"`
	Name         string    `json:"name"`
	NameJP       string    `json:"name_jp"`
	Gender       string    `json:"gender"`
	Age          string    `json:"age"`
	Personality  string    `json:"personality"` // JSON Array representation
	Backstory    string    `json:"backstory"`
	VisualPrompt string    `json:"visual_prompt"`
	Tags         string    `json:"tags"`
	Category     string    `json:"category" gorm:"index"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
