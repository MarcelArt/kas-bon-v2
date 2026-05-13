package routes

import (
	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/MarcelArt/kas-bon-v2/web/handlers"
	"github.com/MarcelArt/kas-bon-v2/web/middlewares"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func SetupWebRoutes(app fiber.Router, userSvc services.IUserService) {
	app.Use("/public", static.New("./web/public"))

	authH := handlers.NewAuthHandler(userSvc)

	app.Get("/login", authH.LoginPage)
	app.Get("/register", authH.RegisterPage)

	auth := app.Group("/auth")
	auth.Post("/login", authH.HandleLogin)
	auth.Post("/register", authH.HandleRegister)
	auth.Post("/logout", func(c fiber.Ctx) error {
		c.ClearCookie("access_token")
		c.ClearCookie("refresh_token")
		return c.Redirect().To("/login")
	})

	protected := app.Group("/", middlewares.CookieAuth())

	appSvc := services.NewAppService(repositories.NewAppRepo(configs.DB))
	appH := handlers.NewAppHandler(appSvc)

	protected.Get("/dashboard", func(c fiber.Ctx) error {
		return c.Redirect().To("/apps")
	})
	protected.Get("/apps", appH.AppsPage)
}
