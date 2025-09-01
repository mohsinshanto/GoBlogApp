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

	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  uint   `json:"userid" gorm:"column:user_id"` // explicitly map to DB column
}
