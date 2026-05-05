package configs

import (
	"fmt"

	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/gofiber/contrib/v3/casbin"
	"github.com/gofiber/fiber/v3"
)

var Authz *casbin.Middleware

func SetupCasbin() error {
	policyAdapter, err := gormadapter.NewAdapter("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to create policy adapter: %w", err)
	}
	authz := casbin.New(casbin.Config{
		ModelFilePath: "rbac_model.conf",
		PolicyAdapter: policyAdapter,
		Lookup: func(c fiber.Ctx) string {
			return "kandar"
		},
	})

	Authz = authz
	return nil
}
