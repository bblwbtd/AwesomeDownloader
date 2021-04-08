package entities

import "gorm.io/gorm"

type Batch struct {
	gorm.Model
	Name  string         `json:"name"`
	Tasks []DownloadTask `json:"tasks" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
