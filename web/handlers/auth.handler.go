package handlers

import (
	"fmt"
	"time"

	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	webModels "github.com/MarcelArt/kas-bon-v2/web/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	userSvc services.IUserService
}

func NewAuthHandler(userSvc services.IUserService) *AuthHandler {
	return &AuthHandler{userSvc: userSvc}
}

func (h *AuthHandler) LoginPage(c fiber.Ctx) error {
	return c.Render("login", webModels.PageData{Title: "Login"})
}

func (h *AuthHandler) RegisterPage(c fiber.Ctx) error {
	return c.Render("register", webModels.PageData{Title: "Register"})
}

func (h *AuthHandler) HandleLogin(c fiber.Ctx) error {
	var form webModels.LoginForm
	if err := c.Bind().Form(&form); err != nil {
		return renderAlert(c, "error", "Invalid form data")
	}

	loginInput := models.LoginInput{
		Username:   form.Username,
		Password:   form.Password,
		IsRemember: form.IsRemember,
	}

	res, err := h.userSvc.Login(loginInput, c)
	if err != nil {
		return renderAlert(c, "error", "Invalid username or password")
	}

	setTokenCookies(c, res.AccessToken, res.RefreshToken, form.IsRemember)

	c.Set("HX-Redirect", "/dashboard")
	return c.SendStatus(fiber.StatusOK)
}

func (h *AuthHandler) HandleRegister(c fiber.Ctx) error {
	var form webModels.RegisterForm
	if err := c.Bind().Form(&form); err != nil {
		return renderAlert(c, "error", "Invalid form data")
	}

	if form.Password != form.ConfirmPassword {
		return renderAlert(c, "error", "Passwords do not match")
	}

	userInput := models.UserInput{
		Username: form.Username,
		Email:    form.Email,
		Password: form.Password,
	}

	if _, err := h.userSvc.Create(userInput); err != nil {
		return renderAlert(c, "error", fmt.Sprintf("Failed to create account: %s", err.Error()))
	}

	loginInput := models.LoginInput{
		Username:   form.Username,
		Password:   form.Password,
		IsRemember: false,
	}

	res, err := h.userSvc.Login(loginInput, c)
	if err != nil {
		return renderAlert(c, "success", "Account created! Please sign in.")
	}

	setTokenCookies(c, res.AccessToken, res.RefreshToken, false)

	c.Set("HX-Redirect", "/dashboard")
	return c.SendStatus(fiber.StatusOK)
}

func setTokenCookies(c fiber.Ctx, accessToken, refreshToken string, isRemember bool) {
	atMaxAge := 5 * 60
	rtMaxAge := int((24 * time.Hour).Seconds())
	if isRemember {
		rtMaxAge = int((30 * 24 * time.Hour).Seconds())
	}

	cookie := fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		MaxAge:   atMaxAge,
		Path:     "/",
	}
	c.Cookie(&cookie)

	cookie = fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		MaxAge:   rtMaxAge,
		Path:     "/",
	}
	c.Cookie(&cookie)
}

func renderAlert(c fiber.Ctx, alertType, message string) error {
	html := fmt.Sprintf(`<div class="alert alert-%s">%s</div>`, alertType, message)
	return c.SendString(html)
}
