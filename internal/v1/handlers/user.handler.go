package handlers

import (
	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/alexedwards/argon2id"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	repo repositories.IUserRepo
}

func NewUserHandler(repo repositories.IUserRepo) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

// @Summary		Create a new user
// @Description	Create a new user with the provided JSON payload
// @Tags			users
// @Accept			json
// @Produce			json
// @Param			request	body		models.UserInput	true	"User object"
// @Success			201		{object}	common.JSONResponse{items=int}
// @Failure			400		{object}	common.JSONResponse
// @Failure			500		{object}	common.JSONResponse
// @Router			/v1/users [post]
func (h *UserHandler) Create(c fiber.Ctx) error {
	var user models.UserInput
	if err := c.Bind().JSON(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	password, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed hashing password"))
	}
	user.Password = password

	id, err := h.repo.Create(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed creating user"))
	}

	return c.Status(fiber.StatusCreated).JSON(common.NewJSONResponse(id, "models.User created"))
}

// @Summary		List users
// @Description	Get a paginated list of users
// @Tags			users
// @Produce			json
// @Param			page	query		int		false	"Page"
// @Param			size	query		int		false	"Size"
// @Param			sort	query		string	false	"Sort"
// @Param			filters	query		string	false	"Filter"
// @Success			200		{object}	common.JSONResponse
// @Failure			500		{object}	common.JSONResponse
// @Router			/v1/users [get]
func (h *UserHandler) Read(c fiber.Ctx) error {
	page, _ := h.repo.Read(c)
	return c.Status(fiber.StatusOK).JSON(page)
}

// @Summary		Update a user
// @Description	Update an existing user by ID
// @Tags			users
// @Accept			json
// @Produce			json
// @Param			id		path		string		true	"User ID"
// @Param			request	body		models.User	true	"User object"
// @Success			200		{object}	common.JSONResponse
// @Failure			400		{object}	common.JSONResponse
// @Failure			500		{object}	common.JSONResponse
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
// @Produce			json
// @Param			id		path		string	true	"User ID"
// @Success			200		{object}	common.JSONResponse
// @Failure			500		{object}	common.JSONResponse
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
// @Produce			json
// @Param			id		path		string	true	"User ID"
// @Success			200		{object}	common.JSONResponse{items=models.User}
// @Failure			500		{object}	common.JSONResponse
// @Router			/v1/users/{id} [get]
func (h *UserHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.repo.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed getting user"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(user, "User found"))
}

// Login authenticates a user and returns JWT tokens
// @Summary Login user
// @Description Authenticate a user with username/email and password, returns access and refresh tokens
// @Tags users
// @Accept json
// @Produce json
// @Param request body models.LoginInput true "Login credentials"
// @Success 200 {object} common.JSONResponse{items=models.LoginResponse} "Authentication successful"
// @Failure 400 {object} common.JSONResponse "Invalid JSON format"
// @Failure 401 {object} common.JSONResponse "Invalid credentials"
// @Failure 500 {object} common.JSONResponse "Internal server error"
// @Router /v1/users/login [post]
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

	claims := map[string]any{
		"sub":    user.Username,
		"userId": user.ID,
		"iss":    c.BaseURL(),
	}
	at, rt, err := common.GenerateJWTPair(claims, login.IsRemember)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed generating tokens"))
	}

	res := models.LoginResponse{
		AccessToken:  at,
		RefreshToken: rt,
		User:         user,
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(res, "Authenticated"))
}
