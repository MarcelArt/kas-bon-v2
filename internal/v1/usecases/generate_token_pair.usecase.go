package usecases

import (
	"fmt"

	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/casbin/casbin/v3"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gofiber/fiber/v3"
)

type GenerateTokenPairUsecase struct {
	c          fiber.Ctx
	user       models.User
	isRemember bool
	e          *casbin.Enforcer
}

func InitGenerateTokenPairUsecase() *GenerateTokenPairUsecase {
	a, _ := gormadapter.NewAdapterByDB(configs.DB)

	e, _ := casbin.NewEnforcer("rbac_model.conf", a)
	return &GenerateTokenPairUsecase{e: e}
}

func (u *GenerateTokenPairUsecase) SetCtx(c fiber.Ctx) *GenerateTokenPairUsecase {
	u.c = c
	return u
}

func (u *GenerateTokenPairUsecase) SetUser(user models.User) *GenerateTokenPairUsecase {
	u.user = user
	return u
}

func (u *GenerateTokenPairUsecase) SetIsRemember(isRemember bool) *GenerateTokenPairUsecase {
	u.isRemember = isRemember
	return u
}

func (u *GenerateTokenPairUsecase) Execute() (res models.LoginResponse, err error) {
	user := u.user
	c := u.c

	permissions, err := u.e.GetImplicitPermissionsForUser(user.Username)
	if err != nil {
		return res, fmt.Errorf("failed retrieving permissions: %w", err)
	}

	claims := map[string]any{
		"sub":    user.Username,
		"userId": user.ID,
		"iss":    c.BaseURL(),
	}
	at, rt, err := common.GenerateJWTPair(claims, permissions, u.isRemember)
	if err != nil {
		return res, fmt.Errorf("failed generating tokens: %w", err)
	}

	res = models.LoginResponse{
		AccessToken:  at,
		RefreshToken: rt,
		User:         user,
	}

	return res, nil
}
