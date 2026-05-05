package routes

import (
	"log"

	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/handlers"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/middlewares"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
)

func SetupUserRoutes(v1 fiber.Router, authz *middlewares.CasbinMiddleware) {
	log.Println("SetupUserRoutes")
	users := v1.Group("/users")

	h := handlers.NewUserHandler(repositories.NewUserRepo(configs.DB))

	users.Get("/",
		authz.HasPermission("users#read"),
		h.Read,
	)
	users.Get("/:id", h.GetByID)

	users.Post("/", h.Create)
	users.Post("/login", h.Login)

	users.Put("/:id", h.Update)

	users.Delete("/:id", h.Delete)
}
