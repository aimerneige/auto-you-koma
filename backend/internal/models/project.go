package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	UserID      string    `gorm:"type:char(36);index;not null" json:"user_id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Mode        string    `gorm:"type:varchar(50)" json:"mode"` // standalone / serialized
	Status      string    `gorm:"type:varchar(50)" json:"status"` // draft / scripted / previewed / rendering / done
	Synopsis    string    `gorm:"type:text" json:"synopsis"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Project) BeforeCreate() error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}