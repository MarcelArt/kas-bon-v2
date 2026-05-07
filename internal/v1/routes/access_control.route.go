package routes

import (
	"log"

	"github.com/MarcelArt/kas-bon-v2/internal/v1/handlers"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/middlewares"
	"github.com/casbin/casbin/v3"
	"github.com/gofiber/fiber/v3"
)

func SetupAccessControlRoutes(v1 fiber.Router, authz *middlewares.CasbinMiddleware, e *casbin.Enforcer) {
	log.Println("SetupAccessControlRoutes")
	g := v1.Group("/access-controls")

	h := handlers.NewAccessControlHandler(e)

	g.Get("/roles/:domain", h.GetAllRoles)
	g.Get("/permissions/:app/:domain/:user", h.GetPermissionsForUser)

	g.Post("/eval", h.Eval)
}
