package entities

import (
	"database/sql"
	"gorm.io/gorm"
)

type DownloadStatus = string

const (
	Pending     DownloadStatus = "Pending"
	Downloading DownloadStatus = "Downloading"
	Paused      DownloadStatus = "Paused"
	Canceled    DownloadStatus = "Canceled"
	Finished    DownloadStatus = "Finished"
	Error       DownloadStatus = "Error"
)

type DownloadTask struct {
	gorm.Model
	URL    string        `json:"url" gorm:"index"`
	Path   string        `json:"path"`
	Size   uint64        `json:"size"`
	Status string        `json:"status" gorm:"index"`
	Batch  sql.NullInt64 `json:"batch" gorm:"index"`
}
