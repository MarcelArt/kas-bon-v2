package handlers

import (
	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/enums"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/casbin/casbin/v3"
	"github.com/gofiber/fiber/v3"
)

type TokenHandler struct {
	e     *casbin.Enforcer
	uRepo repositories.IUserRepo
	aRepo repositories.IAppRepo
	dRepo repositories.IDomainRepo
}

func NewTokenHandler(e *casbin.Enforcer, uRepo repositories.IUserRepo, aRepo repositories.IAppRepo, dRepo repositories.IDomainRepo) *TokenHandler {
	return &TokenHandler{
		e:     e,
		uRepo: uRepo,
		aRepo: aRepo,
		dRepo: dRepo,
	}
}

// @Summary		Check token permission
// @Description	Check if the authenticated user has permission for a given resource and action
// @Tags			token
// @Security		ApiKeyAuth
// @Accept			json
// @Produce			json
// @Param			request		body		models.TokenEndpointRequest	true	"Token request"
// @Success			200			{object}	common.JSONResponse
// @Failure			400			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/token [post]
func (h *TokenHandler) Token(c fiber.Ctx) error {
	claims := common.FiberCtxToClaims(c)
	userID := claims["userId"]

	var payload models.TokenEndpointRequest
	if err := c.Bind().JSON(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing json"))
	}

	user, err := h.uRepo.GetByID(userID)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed getting user"))
	}

	dom, err := h.dRepo.GetByID(payload.DomainID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "invalid domain id header"))
	}

	app, err := h.aRepo.GetByID(payload.AppID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "invalid app id header"))
	}

	res, act := common.ExtractPermissionResourceAndAction(payload.Permission)

	if ok, _ := h.e.Enforce(user.Username, app.Name, dom.Name, enums.ResourceAll, enums.PermissionFull); ok {
		return c.JSON(common.NewJSONResponse(true, "permitted"))
	}

	if ok, _ := h.e.Enforce(user.Username, app.Name, dom.Name, res, act); ok {
		return c.JSON(common.NewJSONResponse(true, "permitted"))
	}

	return c.JSON(common.NewJSONResponse(false, "denied"))
}
