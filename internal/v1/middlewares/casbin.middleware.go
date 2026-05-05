package middlewares

import (
	"strings"

	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type CasbinMiddleware struct {
	e *casbin.Enforcer
}

func NewCasbinMiddleware(db *gorm.DB) *CasbinMiddleware {
	a, _ := gormadapter.NewAdapterByDB(db)

	e, _ := casbin.NewEnforcer("rbac_model.conf", a)

	return &CasbinMiddleware{e: e}
}

func (m *CasbinMiddleware) PolicyLoader(c fiber.Ctx) error {
	m.e.LoadPolicy()
	return c.Next()
}

func (m *CasbinMiddleware) HasPermission(permission string) func(c fiber.Ctx) error {
	permParts := strings.Split(permission, "#")
	res := permParts[0]
	act := permParts[1]
	sub := "kandar"

	return func(c fiber.Ctx) error {
		m.e.LoadPolicy()

		ok, _ := m.e.Enforce(sub, res, act)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(common.NewJSONResponse(ok, "unauthorized"))
		}

		return c.Next()
	}
}
