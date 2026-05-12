package routes

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/MarcelArt/kas-bon-v2/web/handlers"
	"github.com/MarcelArt/kas-bon-v2/web/middlewares"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func SetupWebRoutes(app fiber.Router, userSvc services.IUserService) {
	app.Use("/public", static.New("./web/public"))

	h := handlers.NewAuthHandler(userSvc)

	app.Get("/login", h.LoginPage)
	app.Get("/register", h.RegisterPage)

	auth := app.Group("/auth")
	auth.Post("/login", h.HandleLogin)
	auth.Post("/register", h.HandleRegister)

	protected := app.Group("/", middlewares.CookieAuth())
	protected.Get("/dashboard", func(c fiber.Ctx) error {
		return c.SendString("Dashboard - authenticated!")
	})
}
