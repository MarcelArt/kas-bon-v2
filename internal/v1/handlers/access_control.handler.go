package handlers

import (
	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type AccessControlHandler struct {
	e *casbin.Enforcer
}

func NewAccessControlHandler(db *gorm.DB) *AccessControlHandler {
	a, _ := gormadapter.NewAdapterByDB(db)

	e, _ := casbin.NewEnforcer("rbac_model.conf", a)

	return &AccessControlHandler{e: e}
}

// @Summary		Get all roles
// @Description	Retrieve a list of all roles from the access control system
// @Tags			access-controls
// @Produce			json
// @Success			200		{object}	common.JSONResponse
// @Failure			500		{object}	common.JSONResponse
// @Router			/v1/access-controls [get]
func (h *AccessControlHandler) Read(c fiber.Ctx) error {
	roles, err := h.e.GetAllRoles()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed retrieving roles"))
	}

	return c.Status(fiber.StatusOK).JSON(common.NewJSONResponse(roles, "Success retrieving roles"))
}
