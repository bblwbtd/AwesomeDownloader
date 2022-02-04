package entities

import (
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	URL     string `json:"url" gorm:"index"`
	Path    string `json:"path"`
	Size    uint64 `json:"size"`
	Status  string `json:"status" gorm:"index"`
	Headers string `json:"headers"`
}
