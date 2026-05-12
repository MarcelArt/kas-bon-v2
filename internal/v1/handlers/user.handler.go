package handlers

import (
	"fmt"

	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	svc services.IUserService
}

func NewUserHandler(svc services.IUserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// @Summary		Create a new user
// @Description	Create a new user with the provided JSON payload
// @Tags			users
// @Accept			json
// @Produce			json
// @Param			request		body		models.UserInput	true	"User object"
// @Success			201			{object}	common.JSONResponse{items=int}
// @Failure			400			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/users [post]
func (h *UserHandler) Create(c fiber.Ctx) error {
	var user models.UserInput
	if err := c.Bind().JSON(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	id, err := h.svc.Create(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed creating user"))
	}

	return c.Status(fiber.StatusCreated).JSON(common.NewJSONResponse(id, "User created"))
}

// @Summary		List users
// @Description	Get a paginated list of users
// @Tags			users
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
// @Router			/v1/users [get]
func (h *UserHandler) Read(c fiber.Ctx) error {
	page, _ := h.svc.Read(c)
	return c.Status(fiber.StatusOK).JSON(page)
}

// @Summary		Update a user
// @Description	Update an existing user by ID
// @Tags			users
// @Security		ApiKeyAuth
// @Accept			json
// @Produce			json
// @Param			X-App-Id		header		int				false	"App identifier"
// @Param			X-Domain-Id	header		int				false	"Domain identifier"
// @Param			id			path		string				true	"User ID"
// @Param			request		body		models.User			true	"User object"
// @Success			200			{object}	common.JSONResponse
// @Failure			400			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/users/{id} [put]
func (h *UserHandler) Update(c fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if err := c.Bind().JSON(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	if err := h.svc.Update(id, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed updating user"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "User updated"))
}

// @Summary		Delete a user
// @Description	Delete a user by ID
// @Tags			users
// @Security		ApiKeyAuth
// @Produce			json
// @Param			X-App-Id		header		int	true	"App identifier"
// @Param			X-Domain-Id	header		int	true	"Domain identifier"
// @Param			id			path		string	true	"User ID"
// @Success			200			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/users/{id} [delete]
func (h *UserHandler) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.svc.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed deleting user"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "User deleted"))
}

// @Summary		Get user by ID
// @Description	Retrieve a single user by their ID
// @Tags			users
// @Security		ApiKeyAuth
// @Produce			json
// @Param			X-App-Id		header		int	true	"App identifier"
// @Param			X-Domain-Id	header		int	true	"Domain identifier"
// @Param			id			path		string	true	"User ID"
// @Success			200			{object}	common.JSONResponse{items=models.User}
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/users/{id} [get]
func (h *UserHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.svc.GetByID(id)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting user"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(user, "User found"))
}

// Login authenticates a user and returns JWT tokens
// @Summary		Login user
// @Description	Authenticate a user with username/email and password, returns access and refresh tokens
// @Tags			users
// @Accept			json
// @Produce			json
// @Param			request		body		models.LoginInput		true	"Login credentials"
// @Success			200			{object}	common.JSONResponse{items=models.LoginResponse}	"Authentication successful"
// @Failure			400			{object}	common.JSONResponse	"Invalid JSON format"
// @Failure			401			{object}	common.JSONResponse	"Invalid credentials"
// @Failure			500			{object}	common.JSONResponse	"Internal server error"
// @Router			/v1/users/login [post]
func (h *UserHandler) Login(c fiber.Ctx) error {
	var login models.LoginInput
	if err := c.Bind().JSON(&login); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	res, err := h.svc.Login(login, c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(common.NewJSONResponse(err, "username or password invalid"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(res, "Authenticated"))
}

// @Summary		Refresh token
// @Description	Generate a new token pair using the refresh token
// @Tags			users
// @Produce			json
// @Param			X-Refresh-Token	header		string	true	"Refresh token"
// @Success			200				{object}	common.JSONResponse
// @Failure			401				{object}	common.JSONResponse
// @Failure			500				{object}	common.JSONResponse
// @Router			/v1/users/refresh [post]
func (h *UserHandler) Refresh(c fiber.Ctx) error {
	claims := common.FiberCtxToClaims(c)
	id := claims["userId"]

	isRemember, ok := claims["isRemember"].(bool)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(fmt.Errorf("missing isRemember claim"), "invalid refresh token"))
	}

	res, err := h.svc.Refresh(id, isRemember, c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed generating tokens"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(res, "Authenticated"))
}

// @Summary		Get roles for user
// @Description	Retrieve the roles assigned to a user within a specific domain
// @Tags			users
// @Produce			json
// @Param			id				path		string	true	"User ID"
// @Param			X-App-Id		header		int	true	"App identifier"
// @Param			X-Domain-Id		header		uint	true	"Domain ID"
// @Success			200				{object}	common.JSONResponse
// @Failure			400				{object}	common.JSONResponse
// @Failure			404				{object}	common.JSONResponse
// @Security		ApiKeyAuth
// @Router			/v1/users/{id}/roles [get]
func (h *UserHandler) GetRoles(c fiber.Ctx) error {
	id := c.Params("id")
	domainID := fiber.GetReqHeader[uint](c, "X-Domain-Id")

	roles, err := h.svc.GetRoles(id, domainID)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting roles"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(roles, "Roles found"))
}

// @Summary		Get user permissions
// @Description	Retrieve implicit permissions for a user by ID
// @Tags			users
// @Security		ApiKeyAuth
// @Produce			json
// @Param			X-App-Id		header		int		true	"App identifier"
// @Param			X-Domain-Id	header		int		true	"Domain identifier"
// @Param			id				path		string	true	"Role ID"
// @Success			200				{object}	common.JSONResponse
// @Failure			400				{object}	common.JSONResponse
// @Failure			500				{object}	common.JSONResponse
// @Router			/v1/users/{id}/permissions [get]
func (h *UserHandler) GetPermissions(c fiber.Ctx) error {
	id := c.Params("id")
	domainID := fiber.GetReqHeader[uint](c, "X-Domain-Id")

	permissions, err := h.svc.GetPermissions(id, domainID)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed retrieving permissions"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(permissions, "Success retrieving permissions"))
}

// @Summary		Assign roles to user
// @Description	Assign roles to a user by ID
// @Tags			users
// @Security		ApiKeyAuth
// @Accept			json
// @Produce			json
// @Param			X-App-Id		header		int			true	"App identifier"
// @Param			X-Domain-Id	header		int			true	"Domain identifier"
// @Param			id				path		string		true	"User ID"
// @Param			request			body		[]uint		true	"Role IDs"
// @Success			200				{object}	common.JSONResponse
// @Failure			400				{object}	common.JSONResponse
// @Failure			500				{object}	common.JSONResponse
// @Router			/v1/users/{id}/roles [patch]
func (h *UserHandler) AssignRoles(c fiber.Ctx) error {
	id := c.Params("id")
	domainID := fiber.GetReqHeader[uint](c, "X-Domain-Id")

	var roleIDs []uint
	if err := c.Bind().JSON(&roleIDs); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	roles, err := h.svc.AssignRoles(id, domainID, roleIDs)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed assigning roles"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(roles, "Success assigning roles"))
}

// @Summary		Get user organizations
// @Description	Retrieve organizations/domains for a user by ID
// @Tags			users
// @Security		ApiKeyAuth
// @Produce			json
// @Param			id				path		string	true	"User ID"
// @Success			200				{object}	common.JSONResponse
// @Failure			500				{object}	common.JSONResponse
// @Router			/v1/users/{id}/organizations [get]
func (h *UserHandler) GetOrganizations(c fiber.Ctx) error {
	id := c.Params("id")
	domains, err := h.svc.GetOrganizations(id)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed retrieving domains"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(domains, "Success retrieving domains"))
}
