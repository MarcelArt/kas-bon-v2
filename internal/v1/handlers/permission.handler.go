package handlers

import (
	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

type PermissionHandler struct {
	svc services.IPermissionService
}

func NewPermissionHandler(svc services.IPermissionService) *PermissionHandler {
	return &PermissionHandler{svc: svc}
}

// @Summary		Create a new permission
// @Description	Create a new permission with the provided JSON payload
// @Tags			permissions
// @Security		ApiKeyAuth
// @Accept			json
// @Produce			json
// @Param			X-App-Id	header		int					true	"App identifier"
// @Param			X-Domain-Id	header		int	false	"Domain identifier"
// @Param			request		body		models.PermissionInput	true	"Permission object"
// @Success			201			{object}	common.JSONResponse{items=int}
// @Failure			400			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/permissions [post]
func (h *PermissionHandler) Create(c fiber.Ctx) error {
	var permission models.PermissionInput
	if err := c.Bind().JSON(&permission); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	id, err := h.svc.Create(permission)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed creating permission"))
	}

	return c.Status(fiber.StatusCreated).JSON(common.NewJSONResponse(id, "Permission created"))
}

// @Summary		List permissions
// @Description	Get a paginated list of permissions
// @Tags			permissions
// @Security		ApiKeyAuth
// @Produce			json
// @Param			X-App-Id	header		int		false	"App identifier"
// @Param			X-Domain-Id	header		int	false	"Domain identifier"
// @Param			page		query		int		false	"Page"
// @Param			size		query		int		false	"Size"
// @Param			sort		query		string	false	"Sort"
// @Param			filters		query		string	false	"Filter"
// @Success			200			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/permissions [get]
func (h *PermissionHandler) Read(c fiber.Ctx) error {
	appID := fiber.GetReqHeader[uint](c, "X-App-Id")
	page, _ := h.svc.Read(c, appID)
	return c.Status(fiber.StatusOK).JSON(page)
}

// @Summary		Update a permission
// @Description	Update an existing permission by ID
// @Tags			permissions
// @Security		ApiKeyAuth
// @Accept			json
// @Produce			json
// @Param			X-App-Id	header		int					false	"App identifier"
// @Param			X-Domain-Id	header		int	false	"Domain identifier"
// @Param			id			path		string				true	"Permission ID"
// @Param			request		body		models.Permission	true	"Permission object"
// @Success			200			{object}	common.JSONResponse
// @Failure			400			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/permissions/{id} [put]
func (h *PermissionHandler) Update(c fiber.Ctx) error {
	id := c.Params("id")
	var permission models.Permission
	if err := c.Bind().JSON(&permission); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	if err := h.svc.Update(id, permission); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed updating permission"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Permission updated"))
}

// @Summary		Delete a permission
// @Description	Delete a permission by ID
// @Tags			permissions
// @Security		ApiKeyAuth
// @Produce			json
// @Param			X-App-Id	header		int		true	"App identifier"
// @Param			X-Domain-Id	header		int	false	"Domain identifier"
// @Param			id			path		string	true	"Permission ID"
// @Success			200			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/permissions/{id} [delete]
func (h *PermissionHandler) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed deleting permission"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Permission deleted"))
}

// @Summary		Get permission by ID
// @Description	Retrieve a single permission by their ID
// @Tags			permissions
// @Security		ApiKeyAuth
// @Produce			json
// @Param			X-App-Id	header		int		true	"App identifier"
// @Param			X-Domain-Id	header		int	false	"Domain identifier"
// @Param			id			path		string	true	"Permission ID"
// @Success			200			{object}	common.JSONResponse{items=models.Permission}
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/permissions/{id} [get]
func (h *PermissionHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	permission, err := h.svc.GetByID(id)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting permission"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(permission, "Permission found"))
}
