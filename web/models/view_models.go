package models

import "time"

type PageData struct {
	Title      string
	ActivePage string
}

type AppViewModel struct {
	ID          uint
	Name        string
	Description string
	CreatedAt   time.Time
}

type AppsPageData struct {
	PageData
	Apps []AppViewModel
}

type LoginForm struct {
	Username   string `form:"username"`
	Password   string `form:"password"`
	IsRemember bool   `form:"isRemember"`
}

type RegisterForm struct {
	Username        string `form:"username"`
	Email           string `form:"email"`
	Password        string `form:"password"`
	ConfirmPassword string `form:"confirmPassword"`
}

type AlertData struct {
	Type    string
	Message string
}
