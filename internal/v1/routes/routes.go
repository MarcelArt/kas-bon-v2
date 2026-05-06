package routes

import (
	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/middlewares"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(api fiber.Router) {
	v1 := api.Group("/v1")

	authz := middlewares.NewCasbinMiddleware(configs.DB)
	v1.Use(authz.PolicyLoader)

	SetupUserRoutes(v1, authz)
	SetupRoleRoutes(v1, authz)
	SetupAppRoutes(v1, authz)
	SetupDomainRoutes(v1, authz)
	SetupAccessControlRoutes(v1, authz)
}
