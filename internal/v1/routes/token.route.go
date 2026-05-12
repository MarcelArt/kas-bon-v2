package routes

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/handlers"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/middlewares"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

func SetupTokenRoutes(v1 fiber.Router, svc services.ITokenService) {
	h := handlers.NewTokenHandler(svc)

	g := v1.Group("/token")

	g.Post("/", middlewares.Authn(), h.Token)
}
