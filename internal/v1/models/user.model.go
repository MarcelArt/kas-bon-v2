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

type LoginInput struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	IsRemember bool   `json:"isRemember"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	User         User   `json:"user"`
}

func (UserInput) TableName() string {
	return userTableName
}
