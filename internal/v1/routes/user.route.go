package routes

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/handlers"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/middlewares"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

func SetupUserRoutes(v1 fiber.Router, authz *middlewares.CasbinMiddleware, svc services.IUserService) {
	users := v1.Group("/users")

	h := handlers.NewUserHandler(svc)

	users.Get("/", middlewares.Authn(), authz.HasPermission("users#read"), h.Read)
	users.Get("/:id", middlewares.Authn(), authz.HasPermission("users#read"), h.GetByID)
	users.Get("/:id/roles", middlewares.Authn(), authz.HasPermission("users#read"), h.GetRoles)
	users.Get("/:id/permissions", middlewares.Authn(), authz.HasPermission("users#read"), h.GetPermissions)
	users.Get("/:id/organizations", middlewares.Authn(), h.GetOrganizations)

	users.Post("/", h.Create)
	users.Post("/login", h.Login)
	users.Post("/refresh", middlewares.Refresh(), h.Refresh)

	users.Put("/:id", middlewares.Authn(), authz.HasPermission("users#update"), h.Update)

	users.Patch("/:id/roles", middlewares.Authn(), authz.HasPermission("users#update"), h.AssignRoles)

	users.Delete("/:id", middlewares.Authn(), authz.HasPermission("users#delete"), h.Delete)
}
