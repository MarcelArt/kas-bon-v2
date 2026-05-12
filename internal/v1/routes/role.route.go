package routes

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/handlers"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/middlewares"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

func SetupRoleRoutes(v1 fiber.Router, authz *middlewares.CasbinMiddleware, svc services.IRoleService) {
	roles := v1.Group("/roles")

	h := handlers.NewRoleHandler(svc)

	roles.Get("/", middlewares.Authn(), authz.HasPermission("roles#read"), h.Read)
	roles.Get("/:id", middlewares.Authn(), authz.HasPermission("roles#read"), h.GetByID)
	roles.Get("/:id/permissions", middlewares.Authn(), authz.HasPermission("roles#read"), h.GetPermissions)

	roles.Post("/", middlewares.Authn(), authz.HasPermission("roles#create"), h.Create)

	roles.Put("/:id", middlewares.Authn(), authz.HasPermission("roles#update"), h.Update)

	roles.Patch("/:id/permissions", middlewares.Authn(), authz.HasPermission("roles#update"), h.AssignPermissions)

	roles.Delete("/:id", middlewares.Authn(), authz.HasPermission("roles#delete"), h.Delete)
}
