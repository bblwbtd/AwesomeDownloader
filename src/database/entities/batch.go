package entities

import "gorm.io/gorm"

type Batch struct {
	gorm.Model
	Name  string `json:"name" gorm:"index"`
	Tasks []Task `json:"tasks" gorm:"foreignKey:id"`
}
