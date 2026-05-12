package handlers

import (
	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/gofiber/fiber/v3"
)

type TokenHandler struct {
	svc services.ITokenService
}

func NewTokenHandler(svc services.ITokenService) *TokenHandler {
	return &TokenHandler{svc: svc}
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

	ok, err := h.svc.CheckPermission(userID, payload)
	if err != nil {
		return c.Status(common.StatusCodeFromError(err)).JSON(common.NewJSONResponse(err, "failed checking permission"))
	}

	return c.JSON(common.NewJSONResponse(ok, "permitted"))
}
