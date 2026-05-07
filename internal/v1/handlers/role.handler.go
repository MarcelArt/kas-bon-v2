package handlers

import (
	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
)

type RoleHandler struct {
	repo repositories.IRoleRepo
}

func NewRoleHandler(repo repositories.IRoleRepo) *RoleHandler {
	return &RoleHandler{
		repo: repo,
	}
}

// @Summary		Create a new role
// @Description	Create a new role with the provided JSON payload
// @Tags			roles
// @Security		ApiKeyAuth
// @Accept			json
// @Produce			json
// @Param			X-App-Id		header		int				true	"App identifier"
// @Param			X-Domain-Id	header		int				true	"Domain identifier"
// @Param			request		body		models.RoleInput	true	"Role object"
// @Success			201			{object}	common.JSONResponse{items=int}
// @Failure			400			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/roles [post]
func (h *RoleHandler) Create(c fiber.Ctx) error {
	var role models.RoleInput
	if err := c.Bind().JSON(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	id, err := h.repo.Create(role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed creating role"))
	}

	return c.Status(fiber.StatusCreated).JSON(common.NewJSONResponse(id, "Role created"))
}

// @Summary		List roles
// @Description	Get a paginated list of roles
// @Tags			roles
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
// @Router			/v1/roles [get]
func (h *RoleHandler) Read(c fiber.Ctx) error {
	domID := fiber.GetReqHeader[uint](c, "X-Domain-Id")
	page, _ := h.repo.Read(c, domID)
	return c.Status(fiber.StatusOK).JSON(page)
}

// @Summary		Update a role
// @Description	Update an existing role by ID
// @Tags			roles
// @Security		ApiKeyAuth
// @Accept			json
// @Produce			json
// @Param			X-App-Id		header		int			false	"App identifier"
// @Param			X-Domain-Id	header		int			false	"Domain identifier"
// @Param			id			path		string			true	"Role ID"
// @Param			request		body		models.Role		true	"Role object"
// @Success			200			{object}	common.JSONResponse
// @Failure			400			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/roles/{id} [put]
func (h *RoleHandler) Update(c fiber.Ctx) error {
	id := c.Params("id")
	var role models.Role
	if err := c.Bind().JSON(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	if err := h.repo.Update(id, role); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed updating role"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Role updated"))
}

// @Summary		Delete a role
// @Description	Delete a role by ID
// @Tags			roles
// @Security		ApiKeyAuth
// @Produce			json
// @Param			X-App-Id		header		int	true	"App identifier"
// @Param			X-Domain-Id	header		int	true	"Domain identifier"
// @Param			id			path		string	true	"Role ID"
// @Success			200			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/roles/{id} [delete]
func (h *RoleHandler) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.repo.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed deleting role"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Role deleted"))
}

// @Summary		Get role by ID
// @Description	Retrieve a single role by their ID
// @Tags			roles
// @Security		ApiKeyAuth
// @Produce			json
// @Param			X-App-Id		header		int	true	"App identifier"
// @Param			X-Domain-Id	header		int	true	"Domain identifier"
// @Param			id			path		string	true	"Role ID"
// @Success			200			{object}	common.JSONResponse{items=models.Role}
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/roles/{id} [get]
func (h *RoleHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	role, err := h.repo.GetByID(id)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting role"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(role, "Role found"))
}
