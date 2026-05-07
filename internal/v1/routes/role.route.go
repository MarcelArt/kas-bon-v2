package routes

import (
	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/handlers"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/middlewares"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/casbin/casbin/v3"
	"github.com/gofiber/fiber/v3"
)

func SetupRoleRoutes(v1 fiber.Router, authz *middlewares.CasbinMiddleware, e *casbin.Enforcer) {
	roles := v1.Group("/roles")

	h := handlers.NewRoleHandler(repositories.NewRoleRepo(configs.DB), repositories.NewAppRepo(configs.DB), repositories.NewDomainRepo(configs.DB), repositories.NewPermissionRepo(configs.DB), e)

	roles.Get("/", middlewares.Authn(), authz.HasPermission("roles#read"), h.Read)
	roles.Get("/:id", middlewares.Authn(), authz.HasPermission("roles#read"), h.GetByID)
	roles.Get("/:id/permissions", middlewares.Authn(), authz.HasPermission("roles#read"), h.GetPermissions)

	roles.Post("/", middlewares.Authn(), authz.HasPermission("roles#create"), h.Create)

	roles.Put("/:id", middlewares.Authn(), authz.HasPermission("roles#update"), h.Update)

	roles.Patch("/:id/permissions", middlewares.Authn(), authz.HasPermission("roles#update"), h.AssignPermissions)

	roles.Delete("/:id", middlewares.Authn(), authz.HasPermission("roles#delete"), h.Delete)
}
