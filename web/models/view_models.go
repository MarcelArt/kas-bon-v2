package models

type PageData struct {
	Title string
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
