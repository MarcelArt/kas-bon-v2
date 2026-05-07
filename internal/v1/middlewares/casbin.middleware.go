package middlewares

import (
	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/casbin/casbin/v3"
	"github.com/gofiber/fiber/v3"
)

type CasbinMiddleware struct {
	e     *casbin.Enforcer
	aRepo repositories.IAppRepo
	dRepo repositories.IDomainRepo
}

func NewCasbinMiddleware(e *casbin.Enforcer, aRepo repositories.IAppRepo, dRepo repositories.IDomainRepo) *CasbinMiddleware {
	return &CasbinMiddleware{
		e:     e,
		aRepo: aRepo,
		dRepo: dRepo,
	}
}

func (m *CasbinMiddleware) PolicyLoader(c fiber.Ctx) error {
	if err := m.e.LoadPolicy(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.NewJSONResponse(err, "failed loading policy engine"))
	}
	return c.Next()
}

func (m *CasbinMiddleware) HasPermission(permission string) func(c fiber.Ctx) error {
	res, act := common.ExtractPermissionResourceAndAction(permission)

	return func(c fiber.Ctx) error {
		appID := fiber.GetReqHeader[uint](c, "X-App-Id")
		domID := fiber.GetReqHeader[uint](c, "X-Domain-Id")

		app, err := m.aRepo.GetByID(appID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "invalid app id"))
		}

		dom, err := m.dRepo.GetByID(domID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(common.NewJSONResponse(err, "invalid domain id"))
		}

		claims := common.FiberCtxToClaims(c)
		sub, err := claims.GetSubject()
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(common.NewJSONResponse(err, "invalid token"))
		}

		m.e.LoadPolicy()

		ok := common.IsAuthorized(m.e, sub, app.Name, dom.Name, res, act)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(common.NewJSONResponse(ok, "unauthorized"))
		}

		return c.Next()
	}
}
