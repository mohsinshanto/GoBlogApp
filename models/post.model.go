package models

import (
	"time"

	"gorm.io/gorm"
)

type Blog struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	Title   string `json:"title" gorm:"type:longtext;not null"`
	Content string `json:"content" gorm:"type:longtext;not null"`
	UserID  uint   `json:"user_id" gorm:"not null;index"`

	Published bool `json:"published" gorm:"not null;default:false"`
	Draft     bool `json:"draft" gorm:"not null;default:false"`
}
