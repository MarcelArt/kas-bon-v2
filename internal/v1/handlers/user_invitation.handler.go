package handlers

import (
	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

type UserInvitationHandler struct {
	svc services.IUserInvitationService
}

func NewUserInvitationHandler(svc services.IUserInvitationService) *UserInvitationHandler {
	return &UserInvitationHandler{svc: svc}
}

// @Summary		Create a new user invitation
// @Description	Create a new user invitation with the provided JSON payload
// @Tags			user-invitations
// @Security		ApiKeyAuth
// @Accept			json
// @Produce			json
// @Param			X-App-Id		header		int							true	"App identifier"
// @Param			X-Domain-Id	header		int							true	"Domain identifier"
// @Param			request		body		models.UserInvitationInput	true	"UserInvitation object"
// @Success			201			{object}	common.JSONResponse{items=int}
// @Failure			400			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/user-invitations [post]
func (h *UserInvitationHandler) Create(c fiber.Ctx) error {
	var invitation models.UserInvitationInput
	if err := c.Bind().JSON(&invitation); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	id, err := h.svc.Create(invitation)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed creating user invitation"))
	}

	return c.Status(fiber.StatusCreated).JSON(common.NewJSONResponse(id, "UserInvitation created"))
}

// @Summary		List user invitations
// @Description	Get a paginated list of user invitations
// @Tags			user-invitations
// @Security		ApiKeyAuth
// @Produce			json
// @Param			X-App-Id		header		int	false	"App identifier"
// @Param			X-Domain-Id	header		int	false	"Domain identifier"
// @Param			page		query		int		false	"Page"
// @Param			size		query		int		false	"Size"
// @Param			sort		query		string	false	"Sort"
// @Param			filters		query		string	false	"Filter"
// @Success			200			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/user-invitations [get]
func (h *UserInvitationHandler) Read(c fiber.Ctx) error {
	page, _ := h.svc.Read(c)
	return c.Status(fiber.StatusOK).JSON(page)
}

// @Summary		Update a user invitation
// @Description	Update an existing user invitation by ID
// @Tags			user-invitations
// @Security		ApiKeyAuth
// @Accept			json
// @Produce			json
// @Param			X-App-Id		header		int							false	"App identifier"
// @Param			X-Domain-Id	header		int							false	"Domain identifier"
// @Param			id			path		string						true	"UserInvitation ID"
// @Param			request		body		models.UserInvitation		true	"UserInvitation object"
// @Success			200			{object}	common.JSONResponse
// @Failure			400			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/user-invitations/{id} [put]
func (h *UserInvitationHandler) Update(c fiber.Ctx) error {
	id := c.Params("id")
	var invitation models.UserInvitationInput
	if err := c.Bind().JSON(&invitation); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	if err := h.svc.Update(id, invitation); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed updating user invitation"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "UserInvitation updated"))
}

// @Summary		Delete a user invitation
// @Description	Delete a user invitation by ID
// @Tags			user-invitations
// @Security		ApiKeyAuth
// @Produce			json
// @Param			X-App-Id		header		int	true	"App identifier"
// @Param			X-Domain-Id	header		int	true	"Domain identifier"
// @Param			id			path		string	true	"UserInvitation ID"
// @Success			200			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/user-invitations/{id} [delete]
func (h *UserInvitationHandler) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed deleting user invitation"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "UserInvitation deleted"))
}

// @Summary		Get user invitation by ID
// @Description	Retrieve a single user invitation by their ID
// @Tags			user-invitations
// @Security		ApiKeyAuth
// @Produce			json
// @Param			X-App-Id		header		int	true	"App identifier"
// @Param			X-Domain-Id	header		int	true	"Domain identifier"
// @Param			id			path		string	true	"UserInvitation ID"
// @Success			200			{object}	common.JSONResponse{items=models.UserInvitation}
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/user-invitations/{id} [get]
func (h *UserInvitationHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	invitation, err := h.svc.GetByID(id)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting user invitation"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(invitation, "UserInvitation found"))
}
