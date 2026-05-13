package models

import "gorm.io/gorm"

const permissionTableName = "permissions"

type Permission struct {
	gorm.Model
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`

	AppID uint `gorm:"not null" json:"appId"`

	App *App `json:"app,omitzero"`
}

type PermissionInput struct {
	Input
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`

	AppID uint `gorm:"not null" json:"appId"`
}

func (PermissionInput) TableName() string {
	return permissionTableName
}
