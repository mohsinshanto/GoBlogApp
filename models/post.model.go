package models

import "gorm.io/gorm"

type Blog struct {
	gorm.Model
	Title   string `json:"title"`   // removed binding
	Content string `json:"content"` // removed binding
	UserID  uint   `json:"userid"`
}
