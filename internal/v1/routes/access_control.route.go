package routes

import (
	"log"

	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/handlers"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/middlewares"
	"github.com/gofiber/fiber/v3"
)

func SetupAccessControlRoutes(v1 fiber.Router, authz *middlewares.CasbinMiddleware) {
	log.Println("SetupAccessControlRoutes")
	g := v1.Group("/access-controls")

	h := handlers.NewAccessControlHandler(configs.DB)

	g.Get("/", h.Read)
}
