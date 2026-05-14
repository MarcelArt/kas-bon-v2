package middlewares

import (
	"strconv"

	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/casbin/casbin/v3"
	"github.com/gofiber/fiber/v3"
)

type WebCasbinMiddleware struct {
	e     *casbin.Enforcer
	aRepo repositories.IAppRepo
	dRepo repositories.IDomainRepo
}

func NewWebCasbinMiddleware(e *casbin.Enforcer, aRepo repositories.IAppRepo, dRepo repositories.IDomainRepo) *WebCasbinMiddleware {
	return &WebCasbinMiddleware{
		e:     e,
		aRepo: aRepo,
		dRepo: dRepo,
	}
}

func (m *WebCasbinMiddleware) getContext(c fiber.Ctx) (string, string, string, error) {
	appID, _ := strconv.ParseUint(c.Cookies("current_app_id"), 10, 64)
	domID, _ := strconv.ParseUint(c.Cookies("current_domain_id"), 10, 64)

	app, err := m.aRepo.GetByID(uint(appID))
	if err != nil {
		return "", "", "", err
	}

	dom, err := m.dRepo.GetByID(uint(domID))
	if err != nil {
		return "", "", "", err
	}

	sub, _ := c.Locals("username").(string)
	return sub, app.Name, dom.Name, nil
}

func (m *WebCasbinMiddleware) HasPermission(permission string) func(c fiber.Ctx) error {
	res, act := common.ExtractPermissionResourceAndAction(permission)

	return func(c fiber.Ctx) error {
		sub, appName, domName, err := m.getContext(c)
		if err != nil {
			return c.Status(fiber.StatusForbidden).SendString("Forbidden: invalid app/domain context")
		}

		if sub == "" {
			return c.Redirect().To("/login")
		}

		m.e.LoadPolicy()

		ok := common.IsAuthorized(m.e, sub, appName, domName, res, act)
		if !ok {
			return c.Status(fiber.StatusForbidden).SendString("Forbidden")
		}

		return c.Next()
	}
}

func (m *WebCasbinMiddleware) CheckPermissions(permissions ...string) func(c fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		permMap := make(map[string]bool)

		sub, appName, domName, err := m.getContext(c)
		if err != nil || sub == "" {
			c.Locals("permissions", permMap)
			return c.Next()
		}

		m.e.LoadPolicy()

		for _, p := range permissions {
			res, act := common.ExtractPermissionResourceAndAction(p)
			permMap[p] = common.IsAuthorized(m.e, sub, appName, domName, res, act)
		}

		c.Locals("permissions", permMap)
		return c.Next()
	}
}
