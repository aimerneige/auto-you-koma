package model

import "time"

type Character struct {
	ID           string             `json:"id" gorm:"primaryKey"`
	UserID       string             `json:"user_id" gorm:"index"`
	Name         string             `json:"name"`
	NameJP       string             `json:"name_jp"`
	Gender       string             `json:"gender"`
	Age          string             `json:"age"`
	Personality  string             `json:"personality"` // JSON Array representation
	Backstory    string             `json:"backstory"`
	VisualPrompt string             `json:"visual_prompt"`
	Tags         string             `json:"tags"`
	Category     string             `json:"category" gorm:"index"`
	Images       []CharacterImage   `json:"images,omitempty" gorm:"foreignKey:CharacterID"`
	Variants     []CharacterVariant `json:"variants,omitempty" gorm:"foreignKey:CharacterID"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

type CharacterVariant struct {
	ID                   string    `json:"id" gorm:"primaryKey"`
	CharacterID          string    `json:"character_id" gorm:"index"`
	VariantName          string    `json:"variant_name"`
	PersonalityMod       string    `json:"personality_mod"`
	VisualPromptOverride string    `json:"visual_prompt_override"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type CharacterImage struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	CharacterID string    `json:"character_id" gorm:"index"`
	VariantID   *string   `json:"variant_id" gorm:"index"`
	ImageType   string    `json:"image_type"` // avatar, full_body, chibi, expression, reference
	FilePath    string    `json:"file_path"`
	Description string    `json:"description"`
	IsPrimary   bool      `json:"is_primary"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CharacterGroup struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	UserID      string    `json:"user_id" gorm:"index"`
	GroupName   string    `json:"group_name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CharacterGroupMember struct {
	GroupID     string `json:"group_id" gorm:"primaryKey"`
	CharacterID string `json:"character_id" gorm:"primaryKey"`
	SortOrder   int    `json:"sort_order"`
}
