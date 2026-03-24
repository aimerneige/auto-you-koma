package model

import "time"

type GenerationStatus string

const (
	GenerationPending    GenerationStatus = "pending"
	GenerationProcessing GenerationStatus = "processing"
	GenerationDone       GenerationStatus = "done"
	GenerationFailed     GenerationStatus = "failed"
)

type Generation struct {
	ID             string           `json:"id" gorm:"primaryKey"`
	ProjectID      string           `json:"project_id" gorm:"index"`
	ScriptID       string           `json:"script_id" gorm:"index"`
	Status         GenerationStatus `json:"status"`
	Layout         string           `json:"layout"` // "2x2" or "1x4"
	ResultImageURL string           `json:"result_image_url"`
	Error          string           `json:"error"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
}
