package routes

import (
	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/MarcelArt/kas-bon-v2/web/handlers"
	"github.com/MarcelArt/kas-bon-v2/web/middlewares"
	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"
)

func SetupWebRoutes(app fiber.Router, userSvc services.IUserService) {
	app.Use("/public", static.New("./web/public"))

	authH := handlers.NewAuthHandler(userSvc)

	app.Use("/login", middlewares.GuestOnly())
	app.Use("/register", middlewares.GuestOnly())
	app.Get("/login", authH.LoginPage)
	app.Get("/register", authH.RegisterPage)

	auth := app.Group("/auth")
	auth.Post("/login", authH.HandleLogin)
	auth.Post("/register", authH.HandleRegister)
	auth.Post("/logout", func(c fiber.Ctx) error {
		middlewares.ClearTokenCookies(c)
		return c.Redirect().To("/login")
	})

	protected := app.Group("/", middlewares.CookieAuth(userSvc))

	protected.Get("/dashboard", func(c fiber.Ctx) error {
		return c.Redirect().To("/apps")
	})

	appSvc := services.NewAppService(repositories.NewAppRepo(configs.DB))
	appH := handlers.NewAppHandler(appSvc)
	protected.Get("/apps", appH.AppsPage)
	protected.Get("/apps/new", appH.CreateAppForm)
	protected.Post("/apps", appH.CreateApp)
	protected.Get("/apps/:id/edit", appH.EditAppForm)
	protected.Put("/apps/:id", appH.UpdateApp)
	protected.Delete("/apps/:id", appH.DeleteApp)

	domainSvc := services.NewDomainService(repositories.NewDomainRepo(configs.DB))
	domainH := handlers.NewDomainHandler(domainSvc)
	protected.Get("/domains", domainH.DomainsPage)
	protected.Get("/domains/new", domainH.CreateDomainForm)
	protected.Post("/domains", domainH.CreateDomain)
	protected.Get("/domains/:id/edit", domainH.EditDomainForm)
	protected.Put("/domains/:id", domainH.UpdateDomain)
	protected.Delete("/domains/:id", domainH.DeleteDomain)

	a, _ := gormadapter.NewAdapterByDB(configs.DB)
	e, _ := casbin.NewEnforcer("rbac_model.conf", a)
	roleSvc := services.NewRoleService(
		repositories.NewRoleRepo(configs.DB),
		repositories.NewAppRepo(configs.DB),
		repositories.NewDomainRepo(configs.DB),
		repositories.NewPermissionRepo(configs.DB),
		e,
	)

	domainDetailH := handlers.NewDomainDetailHandler(domainSvc, roleSvc)
	protected.Get("/domains/:id", domainDetailH.DomainDetailPage)
	protected.Get("/domains/:id/roles/new", domainDetailH.CreateRoleForm)
	protected.Post("/domains/:id/roles", domainDetailH.CreateRole)

	protected.Get("/roles/:id/edit", domainDetailH.EditRoleForm)
	protected.Put("/roles/:id", domainDetailH.UpdateRole)
	protected.Delete("/roles/:id", domainDetailH.DeleteRole)
}
