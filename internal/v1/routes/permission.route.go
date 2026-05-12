package routes

import (
	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/handlers"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/middlewares"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

func SetupPermissionRoutes(v1 fiber.Router, authz *middlewares.CasbinMiddleware) {
	permissions := v1.Group("/permissions")

	h := handlers.NewPermissionHandler(services.NewPermissionService(repositories.NewPermissionRepo(configs.DB)))

	permissions.Get("/", middlewares.Authn(), authz.HasPermission("permissions#read"), h.Read)
	permissions.Get("/:id", middlewares.Authn(), authz.HasPermission("permissions#read"), h.GetByID)

	permissions.Post("/", middlewares.Authn(), authz.HasPermission("permissions#create"), h.Create)

	permissions.Put("/:id", middlewares.Authn(), authz.HasPermission("permissions#update"), h.Update)

	permissions.Delete("/:id", middlewares.Authn(), authz.HasPermission("permissions#delete"), h.Delete)
}
