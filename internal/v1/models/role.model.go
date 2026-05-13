package models

import "gorm.io/gorm"

const roleTableName = "roles"

type Role struct {
	gorm.Model
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`

	DomainID uint `gorm:"not null" json:"domainId"`

	Domain *Domain `json:"domain,omitzero"`
}

type RoleInput struct {
	Input
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`

	DomainID uint `gorm:"not null" json:"domainId"`
}

func (RoleInput) TableName() string {
	return roleTableName
}
