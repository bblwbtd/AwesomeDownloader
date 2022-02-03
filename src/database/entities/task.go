package entities

import (
	"database/sql"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	URL    string        `json:"url" gorm:"index"`
	Path   string        `json:"path"`
	Size   uint64        `json:"size"`
	Status string        `json:"status" gorm:"index"`
	Batch  sql.NullInt64 `json:"batch"`
}
