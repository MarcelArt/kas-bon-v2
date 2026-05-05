package models

import "gorm.io/gorm"

const userTableName = "users"

type User struct {
	gorm.Model
	Username string `gorm:"not null;unique" json:"username"`
	Email    string `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null" json:"-"`
}

type UserInput struct {
	Input
	Username string `gorm:"not null;unique" json:"username"`
	Email    string `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null" json:"password"`
}

func (UserInput) TableName() string {
	return userTableName
}
