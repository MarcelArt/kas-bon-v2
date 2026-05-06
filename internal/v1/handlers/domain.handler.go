package handlers

import (
	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/gofiber/fiber/v3"
)

type DomainHandler struct {
	repo repositories.IDomainRepo
}

func NewDomainHandler(repo repositories.IDomainRepo) *DomainHandler {
	return &DomainHandler{
		repo: repo,
	}
}

// @Summary		Create a new domain
// @Description	Create a new domain with the provided JSON payload
// @Tags			domains
// @Security		ApiKeyAuth
// @Accept			json
// @Produce			json
// @Param			X-App		header		string					true	"App identifier"
// @Param			X-Domain	header		string					true	"Domain identifier"
// @Param			request		body		models.DomainInput		true	"Domain object"
// @Success			201			{object}	common.JSONResponse{items=int}
// @Failure			400			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/domains [post]
func (h *DomainHandler) Create(c fiber.Ctx) error {
	var domain models.DomainInput
	if err := c.Bind().JSON(&domain); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	id, err := h.repo.Create(domain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed creating domain"))
	}

	return c.Status(fiber.StatusCreated).JSON(common.NewJSONResponse(id, "models.Domain created"))
}

// @Summary		List domains
// @Description	Get a paginated list of domains
// @Tags			domains
// @Security		ApiKeyAuth
// @Produce			json
// @Param			X-App		header		string	false	"App identifier"
// @Param			X-Domain	header		string	false	"Domain identifier"
// @Param			page		query		int		false	"Page"
// @Param			size		query		int		false	"Size"
// @Param			sort		query		string	false	"Sort"
// @Param			filters		query		string	false	"Filter"
// @Success			200			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/domains [get]
func (h *DomainHandler) Read(c fiber.Ctx) error {
	page, _ := h.repo.Read(c)
	return c.Status(fiber.StatusOK).JSON(page)
}

// @Summary		Update a domain
// @Description	Update an existing domain by ID
// @Tags			domains
// @Security		ApiKeyAuth
// @Accept			json
// @Produce			json
// @Param			X-App		header		string				false	"App identifier"
// @Param			X-Domain	header		string				false	"Domain identifier"
// @Param			id			path		string				true	"Domain ID"
// @Param			request		body		models.Domain		true	"Domain object"
// @Success			200			{object}	common.JSONResponse
// @Failure			400			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/domains/{id} [put]
func (h *DomainHandler) Update(c fiber.Ctx) error {
	id := c.Params("id")
	var domain models.Domain
	if err := c.Bind().JSON(&domain); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	if err := h.repo.Update(id, domain); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed updating domain"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Domain updated"))
}

// @Summary		Delete a domain
// @Description	Delete a domain by ID
// @Tags			domains
// @Security		ApiKeyAuth
// @Produce			json
// @Param			X-App		header		string	true	"App identifier"
// @Param			X-Domain	header		string	true	"Domain identifier"
// @Param			id			path		string	true	"Domain ID"
// @Success			200			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/domains/{id} [delete]
func (h *DomainHandler) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.repo.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed deleting domain"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(nil, "Domain deleted"))
}

// @Summary		Get domain by ID
// @Description	Retrieve a single domain by their ID
// @Tags			domains
// @Security		ApiKeyAuth
// @Produce			json
// @Param			X-App		header		string	true	"App identifier"
// @Param			X-Domain	header		string	true	"Domain identifier"
// @Param			id			path		string	true	"Domain ID"
// @Success			200			{object}	common.JSONResponse{items=models.Domain}
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/domains/{id} [get]
func (h *DomainHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	domain, err := h.repo.GetByID(id)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting domain"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(domain, "Domain found"))
}
