package handlers

import (
	"net/url"

	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/casbin/casbin/v3"
	"github.com/gofiber/fiber/v3"
)

type AccessControlHandler struct {
	e *casbin.Enforcer
}

func NewAccessControlHandler(e *casbin.Enforcer) *AccessControlHandler {
	return &AccessControlHandler{e: e}
}

// @Summary		Get all roles
// @Description	Retrieve a list of all roles from the access control system
// @Tags			access-controls
// @Produce			json
// @Param			domain		path		string	true	"Domain identifier"
// @Success			200			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/access-controls/roles/{domain} [get]
func (h *AccessControlHandler) GetAllRoles(c fiber.Ctx) error {
	domain := c.Params("domain")
	roles, err := h.e.GetAllRolesByDomain(domain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed retrieving roles"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(roles, "Success retrieving roles"))
}

// @Summary		Get permissions for user
// @Description	Retrieve all permissions assigned to a specific user
// @Tags			access-controls
// @Produce			json
// @Param			app			path		string	true	"App identifier"
// @Param			domain		path		string	true	"Domain identifier"
// @Param			user		path		string	true	"User identifier"
// @Success			200			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/access-controls/permissions/{app}/{domain}/{user} [get]
func (h *AccessControlHandler) GetPermissionsForUser(c fiber.Ctx) error {
	user := c.Params("user")
	// app := c.Params("app")
	domain := c.Params("domain")

	user, _ = url.QueryUnescape(user)
	domain, _ = url.QueryUnescape(domain)

	permissions, err := h.e.GetImplicitPermissionsForUser(user, domain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed retrieving permissions"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(permissions, "Success retrieving permissions"))
}

// @Summary		Evaluate access control
// @Description	Evaluate whether a subject has permission to perform an action on an object
// @Tags			access-controls
// @Accept			json
// @Produce			json
// @Param			request		body		models.AccessControlEval	true	"Access control evaluation request"
// @Success			200			{object}	common.JSONResponse
// @Failure			400			{object}	common.JSONResponse
// @Failure			500			{object}	common.JSONResponse
// @Router			/v1/access-controls/eval [post]
func (h *AccessControlHandler) Eval(c fiber.Ctx) error {
	var req models.AccessControlEval
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "failed parsing request"))
	}

	ok := common.IsAuthorized(h.e, req.Sub, req.App, req.Dom, req.Obj, req.Act)
	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(ok, "permitted"))
}
