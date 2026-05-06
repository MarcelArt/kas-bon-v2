package middlewares

import (
	"github.com/MarcelArt/kas-bon-v2/internal/common"
	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
)

func Refresh() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(configs.Env.JwtSecret)},
		Extractor:  extractors.FromHeader("X-Refresh-Token"),
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(common.NewJSONResponse(err, "unauthorized"))
		},
	})
}

func Authn() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(configs.Env.JwtSecret)},
		Extractor:  extractors.FromAuthHeader("Bearer"),
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).JSON(common.NewJSONResponse(err, "unauthorized"))
		},
	})
}
