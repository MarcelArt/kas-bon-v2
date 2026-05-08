package routes

import (
	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/handlers"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/middlewares"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/casbin/casbin/v3"
	"github.com/gofiber/fiber/v3"
)

func SetupTokenRoutes(v1 fiber.Router, e *casbin.Enforcer) {
	h := handlers.NewTokenHandler(e, repositories.NewUserRepo(configs.DB), repositories.NewAppRepo(configs.DB), repositories.NewDomainRepo(configs.DB))

	g := v1.Group("/token")

	g.Post("/", middlewares.Authn(), h.Token)
}
