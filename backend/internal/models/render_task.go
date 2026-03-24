package models

import (
	"time"

	"github.com/google/uuid"
)

type RenderTask struct {
	ID            string    `gorm:"type:char(36);primaryKey" json:"id"`
	ProjectID     string    `gorm:"type:char(36);index;not null" json:"project_id"`
	StoryboardID  string    `gorm:"type:char(36);index;not null" json:"storyboard_id"`
	ExportType    string    `gorm:"type:varchar(50)" json:"export_type"`    // native_text / clean_plate
	Layout        string    `gorm:"type:varchar(50)" json:"layout"`        // 2x2 / 1x4
	ImageWidth    int       `json:"image_width"`
	ImageHeight   int       `json:"image_height"`
	Status        string    `gorm:"type:varchar(50)" json:"status"`        // queued / rendering / qc_check / done / failed
	OutputPaths   string    `gorm:"type:text" json:"output_paths"`        // JSON array of image paths
	ErrorMessage  string    `gorm:"type:text" json:"error_message"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (r *RenderTask) BeforeCreate() error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}