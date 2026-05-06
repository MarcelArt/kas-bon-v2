package models

import "gorm.io/gorm"

const roleTableName = "roles"

type Role struct {
	gorm.Model
	Name        string `gorm:"not null;unique" json:"name"`
	Description string `json:"description"`
}

type RoleInput struct {
	Input
	Name        string `gorm:"not null;unique" json:"name"`
	Description string `json:"description"`
}

func (RoleInput) TableName() string {
	return roleTableName
}
