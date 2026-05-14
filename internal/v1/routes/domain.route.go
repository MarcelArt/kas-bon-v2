package routes

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/handlers"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/middlewares"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

func SetupDomainRoutes(v1 fiber.Router, authz *middlewares.CasbinMiddleware, domainSvc services.IDomainService) {
	domains := v1.Group("/domains")

	h := handlers.NewDomainHandler(domainSvc)

	domains.Get("/", middlewares.Authn(), authz.HasPermission("domains#read"), h.Read)
	domains.Get("/:id", middlewares.Authn(), authz.HasPermission("domains#read"), h.GetByID)
	domains.Get("/:id/users", middlewares.Authn(), authz.HasPermission("domains#read"), h.GetUsers)

	domains.Post("/", middlewares.Authn(), authz.HasPermission("domains#create"), h.Create)

	domains.Put("/:id", middlewares.Authn(), authz.HasPermission("domains#update"), h.Update)

	domains.Delete("/:id", middlewares.Authn(), authz.HasPermission("domains#delete"), h.Delete)
}
