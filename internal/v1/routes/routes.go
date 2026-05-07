package routes

import (
	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/middlewares"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(api fiber.Router) {
	v1 := api.Group("/v1")

	a, _ := gormadapter.NewAdapterByDB(configs.DB)
	e, _ := casbin.NewEnforcer("rbac_model.conf", a)

	authz := middlewares.NewCasbinMiddleware(e, repositories.NewAppRepo(configs.DB), repositories.NewDomainRepo(configs.DB))
	v1.Use(authz.PolicyLoader)

	SetupUserRoutes(v1, authz)
	SetupRoleRoutes(v1, authz, e)
	SetupPermissionRoutes(v1, authz)
	SetupAppRoutes(v1, authz)
	SetupDomainRoutes(v1, authz)
	SetupAccessControlRoutes(v1, authz, e)
}
