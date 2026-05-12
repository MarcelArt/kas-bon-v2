package middlewares

import (
	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func CookieAuth() fiber.Handler {
	return func(c fiber.Ctx) error {
		accessToken := c.Cookies("access_token")
		if accessToken == "" {
			return c.Redirect().To("/login")
		}

		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
			return []byte(configs.Env.JwtSecret), nil
		})

		if err != nil || !token.Valid {
			refreshToken := c.Cookies("refresh_token")
			if refreshToken == "" {
				clearTokenCookies(c)
				return c.Redirect().To("/login")
			}

			rt, err := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
				return []byte(configs.Env.JwtSecret), nil
			})

			if err != nil || !rt.Valid {
				clearTokenCookies(c)
				return c.Redirect().To("/login")
			}

			claims, ok := rt.Claims.(jwt.MapClaims)
			if !ok {
				clearTokenCookies(c)
				return c.Redirect().To("/login")
			}

			c.Locals("userID", claims["userId"])
			c.Locals("username", claims["sub"])
			c.Locals("needsRefresh", true)
			return c.Next()
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			clearTokenCookies(c)
			return c.Redirect().To("/login")
		}

		c.Locals("userID", claims["userId"])
		c.Locals("username", claims["sub"])
		c.Locals("needsRefresh", false)
		return c.Next()
	}
}

func clearTokenCookies(c fiber.Ctx) {
	c.ClearCookie("access_token")
	c.ClearCookie("refresh_token")
}
