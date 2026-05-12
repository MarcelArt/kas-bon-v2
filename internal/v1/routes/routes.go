package routes

import (
	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/middlewares"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/repositories"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
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

	userSvc := services.NewUserService(
		repositories.NewUserRepo(configs.DB),
		repositories.NewDomainRepo(configs.DB),
		repositories.NewRoleRepo(configs.DB),
		e,
	)
	roleSvc := services.NewRoleService(
		repositories.NewRoleRepo(configs.DB),
		repositories.NewAppRepo(configs.DB),
		repositories.NewDomainRepo(configs.DB),
		repositories.NewPermissionRepo(configs.DB),
		e,
	)
	acSvc := services.NewAccessControlService(e)
	tokenSvc := services.NewTokenService(
		e,
		repositories.NewUserRepo(configs.DB),
		repositories.NewAppRepo(configs.DB),
		repositories.NewDomainRepo(configs.DB),
	)

	SetupUserRoutes(v1, authz, userSvc)
	SetupRoleRoutes(v1, authz, roleSvc)
	SetupPermissionRoutes(v1, authz)
	SetupAppRoutes(v1, authz)
	SetupDomainRoutes(v1, authz)
	SetupAccessControlRoutes(v1, authz, acSvc)
	SetupTokenRoutes(v1, tokenSvc)
}
