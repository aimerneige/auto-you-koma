package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           string    `gorm:"type:char(36);primaryKey" json:"id"`
	Email        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null" json:"-"`
	DisplayName  string    `gorm:"type:varchar(255)" json:"display_name"`
	TOTPSecret   string    `gorm:"type:varchar(255)" json:"-"`
	TOTPEnabled  bool      `gorm:"default:false" json:"totp_enabled"`
	QuotaLimit   int       `gorm:"default:100" json:"quota_limit"`
	QuotaUsed    int       `gorm:"default:0" json:"quota_used"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate() error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}