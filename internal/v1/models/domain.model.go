package models

import "gorm.io/gorm"

const domainTableName = "domains"

type Domain struct {
	gorm.Model
	Name           string `gorm:"not null" json:"name"`
	Description    string `json:"description"`
	IsOrganization bool   `gorm:"not null" json:"isOrganization"`

	ParentID *uint `json:"parentId"`

	Parent *Domain `json:"parent,omitzero"`
}

type DomainInput struct {
	Input
	Name           string `gorm:"not null" json:"name"`
	Description    string `json:"description"`
	IsOrganization bool   `gorm:"not null" json:"isOrganization"`

	ParentID *uint `json:"parentId"`

	Parent *Domain `json:"parent,omitzero"`
}

func (DomainInput) TableName() string {
	return domainTableName
}
