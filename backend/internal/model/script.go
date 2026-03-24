package model

import "time"

type Script struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	ProjectID  string    `json:"project_id" gorm:"index"` // FK to projects
	SeriesID   string    `json:"series_id" gorm:"index"`  // FK to series
	Title      string    `json:"title"`
	Content    string    `json:"content"`                 // the raw long-form story
	ParsedData string    `json:"parsed_data"`             // JSON representation of panels
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
