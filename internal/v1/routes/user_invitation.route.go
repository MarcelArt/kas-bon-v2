package routes

import (
	"github.com/MarcelArt/kas-bon-v2/internal/v1/handlers"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/middlewares"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

func SetupUserInvitationRoutes(v1 fiber.Router, authz *middlewares.CasbinMiddleware, userInvitationSvc services.IUserInvitationService) {
	userInvitations := v1.Group("/user-invitations")

	h := handlers.NewUserInvitationHandler(userInvitationSvc)

	userInvitations.Get("/", middlewares.Authn(), h.Read)
	userInvitations.Get("/:id", middlewares.Authn(), h.GetByID)

	userInvitations.Post("/", middlewares.Authn(), authz.HasPermission("user-invitations#create"), h.Create)

	userInvitations.Put("/:id", middlewares.Authn(), authz.HasPermission("user-invitations#update"), h.Update)

	userInvitations.Delete("/:id", middlewares.Authn(), authz.HasPermission("user-invitations#delete"), h.Delete)
}
