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

func SetupWebRoutes(app fiber.Router, userSvc services.IUserService, e *casbin.Enforcer) {
	a, _ := gormadapter.NewAdapterByDB(configs.DB)
	ce, _ := casbin.NewEnforcer("rbac_model.conf", a)

	app.Use("/public", static.New("./web/public"))

	authz := middlewares.NewWebCasbinMiddleware(ce, repositories.NewAppRepo(configs.DB), repositories.NewDomainRepo(configs.DB))

	allWebPermissions := []string{
		"apps#read", "apps#create", "apps#update", "apps#delete",
		"domains#read", "domains#create", "domains#update", "domains#delete",
		"roles#read", "roles#create", "roles#update", "roles#delete",
		"permissions#read", "permissions#create", "permissions#update", "permissions#delete",
		"userInvitations#read", "userInvitations#create", "userInvitations#update", "userInvitations#delete",
	}

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
		middlewares.ClearContextCookies(c)
		return c.Redirect().To("/login")
	})

	authed := app.Group("/", middlewares.CookieAuth(userSvc))
	authed.Get("/select-org", authH.SelectOrgPage)
	authed.Post("/select-org", authH.HandleSelectOrg)

	protected := app.Group("/", middlewares.CookieAuth(userSvc), middlewares.RequireContext(), authz.CheckPermissions(allWebPermissions...))

	protected.Get("/dashboard", func(c fiber.Ctx) error {
		return c.Redirect().To("/apps")
	})

	appSvc := services.NewAppService(repositories.NewAppRepo(configs.DB))
	appH := handlers.NewAppHandler(appSvc)
	protected.Get("/apps", authz.HasPermission("apps#read"), appH.AppsPage)
	protected.Get("/apps/new", authz.HasPermission("apps#create"), appH.CreateAppForm)
	protected.Post("/apps", authz.HasPermission("apps#create"), appH.CreateApp)
	protected.Get("/apps/:id/edit", authz.HasPermission("apps#update"), appH.EditAppForm)
	protected.Put("/apps/:id", authz.HasPermission("apps#update"), appH.UpdateApp)
	protected.Delete("/apps/:id", authz.HasPermission("apps#delete"), appH.DeleteApp)

	permSvc := services.NewPermissionService(repositories.NewPermissionRepo(configs.DB))
	appDetailH := handlers.NewAppDetailHandler(appSvc, permSvc)
	protected.Get("/apps/:id", authz.HasPermission("permissions#read"), appDetailH.AppDetailPage)
	protected.Get("/apps/:id/permissions/new", authz.HasPermission("permissions#create"), appDetailH.CreatePermissionForm)
	protected.Post("/apps/:id/permissions", authz.HasPermission("permissions#create"), appDetailH.CreatePermission)
	protected.Get("/permissions/:id/edit", authz.HasPermission("permissions#update"), appDetailH.EditPermissionForm)
	protected.Put("/permissions/:id", authz.HasPermission("permissions#update"), appDetailH.UpdatePermission)
	protected.Delete("/permissions/:id", authz.HasPermission("permissions#delete"), appDetailH.DeletePermission)

	domainSvc := services.NewDomainService(
		repositories.NewDomainRepo(configs.DB),
		repositories.NewUserRepo(configs.DB),
		e,
	)
	domainH := handlers.NewDomainHandler(domainSvc)
	protected.Get("/domains", authz.HasPermission("domains#read"), domainH.DomainsPage)
	protected.Get("/domains/new", authz.HasPermission("domains#create"), domainH.CreateDomainForm)
	protected.Post("/domains", authz.HasPermission("domains#create"), domainH.CreateDomain)
	protected.Get("/domains/:id/edit", authz.HasPermission("domains#update"), domainH.EditDomainForm)
	protected.Put("/domains/:id", authz.HasPermission("domains#update"), domainH.UpdateDomain)
	protected.Delete("/domains/:id", authz.HasPermission("domains#delete"), domainH.DeleteDomain)

	roleSvc := services.NewRoleService(
		repositories.NewRoleRepo(configs.DB),
		repositories.NewAppRepo(configs.DB),
		repositories.NewDomainRepo(configs.DB),
		repositories.NewPermissionRepo(configs.DB),
		e,
	)

	domainDetailH := handlers.NewDomainDetailHandler(domainSvc, roleSvc, userSvc, services.NewUserInvitationService(repositories.NewUserInvitationRepo(configs.DB)))
	protected.Get("/domains/:id", authz.HasPermission("roles#read"), domainDetailH.DomainDetailPage)
	protected.Get("/domains/:id/roles/new", authz.HasPermission("roles#create"), domainDetailH.CreateRoleForm)
	protected.Post("/domains/:id/roles", authz.HasPermission("roles#create"), domainDetailH.CreateRole)
	protected.Get("/domains/:id/invitations/new", authz.HasPermission("userInvitations#create"), domainDetailH.InviteUserForm)
	protected.Post("/domains/:id/invitations", authz.HasPermission("userInvitations#create"), domainDetailH.CreateInvitation)

	protected.Get("/roles/:id/edit", authz.HasPermission("roles#update"), domainDetailH.EditRoleForm)
	protected.Put("/roles/:id", authz.HasPermission("roles#update"), domainDetailH.UpdateRole)
	protected.Delete("/roles/:id", authz.HasPermission("roles#delete"), domainDetailH.DeleteRole)

	rolePermH := handlers.NewRolePermissionHandler(roleSvc, permSvc, appSvc)
	protected.Get("/roles/:id/permissions", authz.HasPermission("roles#read"), rolePermH.PermissionsPage)
	protected.Get("/roles/:id/permissions/list", authz.HasPermission("roles#read"), rolePermH.PermissionsList)
	protected.Post("/roles/:id/permissions", authz.HasPermission("roles#update"), rolePermH.AssignPermissions)

	userRoleH := handlers.NewUserRoleHandler(userSvc, roleSvc, domainSvc)
	protected.Get("/domains/:domainId/users/:userId/roles", authz.HasPermission("roles#read"), userRoleH.UserRolesPage)
	protected.Post("/domains/:domainId/users/:userId/roles", authz.HasPermission("roles#update"), userRoleH.AssignRoles)
}
