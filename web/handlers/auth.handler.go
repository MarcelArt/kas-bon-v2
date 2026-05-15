package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/enums"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/MarcelArt/kas-bon-v2/web/middlewares"
	webModels "github.com/MarcelArt/kas-bon-v2/web/models"
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
	setAppCookie(c)
	c.Locals("userID", res.User.ID)
	c.Locals("username", res.User.Username)

	c.Set("HX-Redirect", "/select-org")
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

	userID, err := h.userSvc.Create(userInput)
	if err != nil {
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
	setAppCookie(c)
	c.Locals("userID", userID)
	c.Locals("username", form.Username)

	c.Set("HX-Redirect", "/select-org")
	return c.SendStatus(fiber.StatusOK)
}

func (h *AuthHandler) SelectOrgPage(c fiber.Ctx) error {
	userID := c.Locals("userID")
	orgs, err := h.userSvc.GetOrganizations(userID)
	if err != nil || len(orgs) == 0 {
		middlewares.ClearTokenCookies(c)
		middlewares.ClearContextCookies(c)
		return c.Redirect().To("/login")
	}

	if len(orgs) == 1 {
		setDomainCookie(c, orgs[0].ID)
		return c.Redirect().To("/apps")
	}

	viewOrgs := make([]webModels.OrgViewModel, len(orgs))
	for i, o := range orgs {
		viewOrgs[i] = webModels.OrgViewModel{
			ID:          o.ID,
			Name:        o.Name,
			Description: o.Description,
		}
	}

	data := webModels.OrgSelectPageData{
		PageData: webModels.PageData{
			Title: "Select Organization",
		},
		Organizations: viewOrgs,
	}

	return c.Render("select_org", data)
}

func (h *AuthHandler) HandleSelectOrg(c fiber.Ctx) error {
	domainID := c.FormValue("domain_id")
	if domainID == "" {
		return c.Redirect().To("/select-org")
	}

	id := parseUint(domainID)
	if id == 0 {
		return c.Redirect().To("/select-org")
	}

	setDomainCookie(c, id)

	return c.Redirect().To("/apps")
}

func isProduction() bool {
	return configs.Env.ServerENV == "prod"
}

func setTokenCookies(c fiber.Ctx, accessToken, refreshToken string, isRemember bool) {
	atMaxAge := 5 * 60
	rtMaxAge := int((24 * time.Hour).Seconds())
	if isRemember {
		rtMaxAge = int((30 * 24 * time.Hour).Seconds())
	}
	secure := isProduction()

	cookie := fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "Strict",
		MaxAge:   atMaxAge,
		Path:     "/",
	}
	c.Cookie(&cookie)

	cookie = fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   secure,
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

func setAppCookie(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "current_app_id",
		Value:    strconv.FormatUint(uint64(enums.AppID), 10),
		HTTPOnly: true,
		Secure:   isProduction(),
		SameSite: "Strict",
		Path:     "/",
	})
}

func setDomainCookie(c fiber.Ctx, domainID uint) {
	c.Cookie(&fiber.Cookie{
		Name:     "current_domain_id",
		Value:    strconv.FormatUint(uint64(domainID), 10),
		HTTPOnly: true,
		Secure:   isProduction(),
		SameSite: "Strict",
		Path:     "/",
	})
}
