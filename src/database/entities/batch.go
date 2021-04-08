package entities

import "gorm.io/gorm"

type Batch struct {
	gorm.Model
	Name string `json:"name"`
}
