package routes

import (
	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/handlers"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/middlewares"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
)

func SetupUserRoutes(v1 fiber.Router, authz *middlewares.CasbinMiddleware) {
	users := v1.Group("/users")

	h := handlers.NewUserHandler(repositories.NewUserRepo(configs.DB))

	users.Get("/", middlewares.Authn(), authz.HasPermission("users#read"), h.Read)
	users.Get("/:id", middlewares.Authn(), authz.HasPermission("users#read"), h.GetByID)

	users.Post("/", h.Create)
	users.Post("/login", h.Login)
	users.Post("/refresh", middlewares.Refresh(), h.Refresh)

	users.Put("/:id", middlewares.Authn(), authz.HasPermission("users#update"), h.Update)

	users.Delete("/:id", middlewares.Authn(), authz.HasPermission("users#delete"), h.Delete)
}
