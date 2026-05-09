package handlers

import (
	"fmt"

	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/usecases"
	"github.com/alexedwards/argon2id"
	"github.com/casbin/casbin/v3"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	repo  repositories.IUserRepo
	dRepo repositories.IDomainRepo
	rRepo repositories.IRoleRepo
	e     *casbin.Enforcer
}

func NewUserHandler(repo repositories.IUserRepo, dRepo repositories.IDomainRepo, rRepo repositories.IRoleRepo, e *casbin.Enforcer) *UserHandler {
	return &UserHandler{
		repo:  repo,
		dRepo: dRepo,
		rRepo: rRepo,
		e:     e,
	}
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

	tx := configs.DB.Begin()
	defer tx.Rollback()

	registerUser := usecases.InitRegisterUserUsecase(tx)
	registerUser.User = user

	id, err := registerUser.Execute()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed creating user"))
	}

	tx.Commit()
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
	page, _ := h.repo.Read(c)
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

	if err := h.repo.Update(id, user); err != nil {
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
	if err := h.repo.Delete(id); err != nil {
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
	user, err := h.repo.GetByID(id)
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

	user, err := h.repo.GetByUsernameOrEmail(login.Username)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(common.NewJSONResponse(err, "username or password invalid"))
	}

	ok, err := argon2id.ComparePasswordAndHash(login.Password, user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "unexpected error"))
	}

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(common.NewJSONResponse(err, "username or password invalid"))
	}

	res, err := usecases.InitGenerateTokenPairUsecase().SetCtx(c).SetUser(user).SetIsRemember(login.IsRemember).Execute()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed generating tokens"))
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

	user, err := h.repo.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(common.NewJSONResponse(err, "failed getting user"))
	}

	res, err := usecases.InitGenerateTokenPairUsecase().SetCtx(c).SetUser(user).SetIsRemember(isRemember).Execute()
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

	user, err := h.repo.GetByID(id)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting user"))
	}

	domain, err := h.dRepo.GetByID(domainID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "invalid domain id"))
	}

	roles := h.e.GetRolesForUserInDomain(user.Username, domain.Name)

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

	user, err := h.repo.GetByID(id)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting user"))
	}
	dom, err := h.dRepo.GetByID(domainID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "invalid domain id header"))
	}

	permissions, err := h.e.GetImplicitPermissionsForUser(user.Username, dom.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed retrieving permissions"))
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

	user, err := h.repo.GetByID(id)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting user"))
	}

	dom, err := h.dRepo.GetByID(domainID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "invalid domain id header"))
	}

	if _, err := h.e.RemoveFilteredGroupingPolicy(0, user.Username, "", dom.Name); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed removing policy"))
	}

	roles, err := h.rRepo.GetDistinctByIDs(roleIDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed retrieving roles"))
	}

	for _, role := range roles {
		if _, err := h.e.AddGroupingPolicy(user.Username, role, dom.Name); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed adding policy"))
		}
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
	user, err := h.repo.GetByID(id)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting user"))
	}

	doms, err := h.e.GetDomainsForUser(user.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed retrieving domains"))
	}

	domains, err := h.dRepo.GetOrganizationsByNames(doms)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed retrieving domains"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(domains, "Success retrieving domains"))
}
