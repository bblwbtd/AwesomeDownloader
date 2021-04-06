package entities

import "gorm.io/gorm"

type DownloadTask struct {
    gorm.Model
    URL string `json:"url"`
    Path string `json:"path"`
    Size uint64 `json:"size"`
}

