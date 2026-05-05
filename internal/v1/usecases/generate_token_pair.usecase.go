package usecases

import (
	"fmt"

	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/models"
	"github.com/gofiber/fiber/v3"
)

type GenerateTokenPairUsecase struct {
	c          fiber.Ctx
	user       models.User
	isRemember bool
}

func InitGenerateTokenPairUsecase() *GenerateTokenPairUsecase {
	return &GenerateTokenPairUsecase{}
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

	claims := map[string]any{
		"sub":    user.Username,
		"userId": user.ID,
		"iss":    c.BaseURL(),
	}
	at, rt, err := common.GenerateJWTPair(claims, u.isRemember)
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
