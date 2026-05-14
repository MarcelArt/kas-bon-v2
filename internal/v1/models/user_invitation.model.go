package models

import (
	"time"

	"gorm.io/gorm"
)

const userInvitationTableName = "user_invitations"

type UserInvitation struct {
	gorm.Model
	AcceptedAt *time.Time `json:"acceptedAt"`
	RejectedAt *time.Time `json:"rejectedAt"`

	UserID   uint `gorm:"not null" json:"userId"`
	DomainID uint `gorm:"not null" json:"domainId"`
	RoleID   uint `gorm:"not null" json:"roleId"`

	User   *User   `json:"user,omitzero"`
	Domain *Domain `json:"domain,omitzero"`
	Role   *Role   `json:"role,omitzero"`
}

type UserInvitationInput struct {
	Input
	AcceptedAt *time.Time `json:"acceptedAt"`
	RejectedAt *time.Time `json:"rejectedAt"`

	UserID   uint `gorm:"not null" json:"userId"`
	DomainID uint `gorm:"not null" json:"domainId"`
	RoleID   uint `gorm:"not null" json:"roleId"`
}

func (UserInvitationInput) TableName() string {
	return userInvitationTableName
}

type UserInvitationPage struct {
	ID         uint       `gorm:"primarykey"`
	CreatedAt  time.Time  `json:"createdAt"`
	AcceptedAt *time.Time `json:"acceptedAt"`
	RejectedAt *time.Time `json:"rejectedAt"`
	UserID     uint       `gorm:"not null" json:"userId"`
	DomainID   uint       `gorm:"not null" json:"domainId"`
	RoleID     uint       `gorm:"not null" json:"roleId"`
	User       string     `json:"user"`
	Domain     string     `json:"domain"`
	Role       string     `json:"role"`
}
