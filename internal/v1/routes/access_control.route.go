package routes

import (
	"log"

	"github.com/MarcelArt/kas-bon-v2/internal/v1/handlers"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/middlewares"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

func SetupAccessControlRoutes(v1 fiber.Router, authz *middlewares.CasbinMiddleware, svc services.IAccessControlService) {
	log.Println("SetupAccessControlRoutes")
	g := v1.Group("/access-controls")

	h := handlers.NewAccessControlHandler(svc)

	g.Get("/roles/:domain", h.GetAllRoles)
	g.Get("/permissions/:app/:domain/:user", h.GetPermissionsForUser)

	g.Post("/eval", h.Eval)
}
