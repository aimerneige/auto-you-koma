package model

import "time"

type CharacterState struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	SeriesID      string    `json:"series_id" gorm:"index"`
	CharacterID   string    `json:"character_id" gorm:"index"`
	Health        int       `json:"health"`
	Sanity        int       `json:"sanity"`
	Inventory     string    `json:"inventory"`      // JSON string array
	MemorySummary string    `json:"memory_summary"` // Continuity context paragraph
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
