package models

import "gorm.io/gorm"

const appTableName = "apps"

type App struct {
	gorm.Model
	Name        string `gorm:"not null;unique" json:"name"`
	Description string `json:"description"`
}

type AppInput struct {
	Input
	Name        string `gorm:"not null;unique" json:"name"`
	Description string `json:"description"`
}

func (AppInput) TableName() string {
	return appTableName
}
