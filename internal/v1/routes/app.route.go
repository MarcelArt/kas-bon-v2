package routes

import (
	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/handlers"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/middlewares"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

func SetupAppRoutes(v1 fiber.Router, authz *middlewares.CasbinMiddleware) {
	apps := v1.Group("/apps")

	h := handlers.NewAppHandler(services.NewAppService(repositories.NewAppRepo(configs.DB)))

	apps.Get("/", middlewares.Authn(), authz.HasPermission("apps#read"), h.Read)
	apps.Get("/:id", middlewares.Authn(), authz.HasPermission("apps#read"), h.GetByID)

	apps.Post("/", middlewares.Authn(), authz.HasPermission("apps#create"), h.Create)

	apps.Put("/:id", middlewares.Authn(), authz.HasPermission("apps#update"), h.Update)

	apps.Delete("/:id", middlewares.Authn(), authz.HasPermission("apps#delete"), h.Delete)
}
