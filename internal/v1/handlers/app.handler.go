package handlers

import (
	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
)

type AppHandler struct {
	repo repositories.IAppRepo
}

func NewAppHandler(repo repositories.IAppRepo) *AppHandler {
	return &AppHandler{
		repo: repo,
	}
}

// @Summary		Create a new app
// @Description	Create a new app with the provided JSON payload
// @Tags			apps
// @Accept			json
// @Produce			json
// @Param			request	body		models.AppInput	true	"App object"
// @Success			201		{object}	common.JSONResponse{items=int}
// @Failure			400		{object}	common.JSONResponse
// @Failure			500		{object}	common.JSONResponse
// @Router			/v1/apps [post]
func (h *AppHandler) Create(c fiber.Ctx) error {
	var app models.AppInput
	if err := c.Bind().JSON(&app); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	id, err := h.repo.Create(app)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed creating app"))
	}

	return c.Status(fiber.StatusCreated).JSON(common.NewJSONResponse(id, "App created"))
}

// @Summary		List apps
// @Description	Get a paginated list of apps
// @Tags			apps
// @Produce			json
// @Param			page	query		int		false	"Page"
// @Param			size	query		int		false	"Size"
// @Param			sort	query		string	false	"Sort"
// @Param			filters	query		string	false	"Filter"
// @Success			200		{object}	common.JSONResponse
// @Failure			500		{object}	common.JSONResponse
// @Router			/v1/apps [get]
func (h *AppHandler) Read(c fiber.Ctx) error {
	page, _ := h.repo.Read(c)
	return c.Status(fiber.StatusOK).JSON(page)
}

// @Summary		Update an app
// @Description	Update an existing app by ID
// @Tags			apps
// @Accept			json
// @Produce			json
// @Param			id		path		string		true	"App ID"
// @Param			request	body		models.App	true	"App object"
// @Success			200		{object}	common.JSONResponse
// @Failure			400		{object}	common.JSONResponse
// @Failure			500		{object}	common.JSONResponse
// @Router			/v1/apps/{id} [put]
func (h *AppHandler) Update(c fiber.Ctx) error {
	id := c.Params("id")
	var app models.App
	if err := c.Bind().JSON(&app); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	if err := h.repo.Update(id, app); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed updating app"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "App updated"))
}

// @Summary		Delete an app
// @Description	Delete an app by ID
// @Tags			apps
// @Produce			json
// @Param			id		path		string	true	"App ID"
// @Success			200		{object}	common.JSONResponse
// @Failure			500		{object}	common.JSONResponse
// @Router			/v1/apps/{id} [delete]
func (h *AppHandler) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.repo.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed deleting app"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "App deleted"))
}

// @Summary		Get app by ID
// @Description	Retrieve a single app by their ID
// @Tags			apps
// @Produce			json
// @Param			id		path		string	true	"App ID"
// @Success			200		{object}	common.JSONResponse{items=models.App}
// @Failure			500		{object}	common.JSONResponse
// @Router			/v1/apps/{id} [get]
func (h *AppHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	app, err := h.repo.GetByID(id)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting app"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(app, "App found"))
}
