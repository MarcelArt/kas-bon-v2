package common

import (
	"fmt"
	"time"

	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/enums"
	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTPair(claims map[string]any, permissions any, isRemember bool) (string, string, error) {
	today := time.Now()
	atExp := time.Minute * 5
	rtExp := enums.Day
	if isRemember {
		rtExp = enums.Month
	}

	atClaims := jwt.MapClaims{
		"iat":         today.Unix(),
		"exp":         today.Add(atExp).Unix(),
		"permissions": permissions,
	}

	rtClaims := jwt.MapClaims{
		"iat":        today.Unix(),
		"exp":        today.Add(rtExp).Unix(),
		"isRemember": isRemember,
	}

	for k, v := range claims {
		atClaims[k] = v
		rtClaims[k] = v
	}

	at, err := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims).SignedString([]byte(configs.Env.JwtSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed generating access token: %w", err)
	}

	rt, err := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims).SignedString([]byte(configs.Env.JwtSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed generating refresh token: %w", err)
	}

	return at, rt, nil
}

func FiberCtxToClaims(c fiber.Ctx) jwt.MapClaims {
	token := jwtware.FromContext(c)
	return token.Claims.(jwt.MapClaims)
}
