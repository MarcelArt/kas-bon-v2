package middlewares

import (
	"time"

	"github.com/MarcelArt/kas-bon-v2/internal/configs"
	"github.com/MarcelArt/kas-bon-v2/internal/v1/services"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func CookieAuth(userSvc services.IUserService) fiber.Handler {
	return func(c fiber.Ctx) error {
		accessToken := c.Cookies("access_token")
		refreshToken := c.Cookies("refresh_token")

		if accessToken == "" && refreshToken == "" {
			return c.Redirect().To("/login")
		}

		if accessToken != "" {
			token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
				return []byte(configs.Env.JwtSecret), nil
			})

			if err == nil && token.Valid {
				claims, ok := token.Claims.(jwt.MapClaims)
				if !ok {
					ClearTokenCookies(c)
					return c.Redirect().To("/login")
				}

				c.Locals("userID", claims["userId"])
				c.Locals("username", claims["sub"])
				return c.Next()
			}
		}

		if refreshToken == "" {
			ClearTokenCookies(c)
			return c.Redirect().To("/login")
		}

		rt, err := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
			return []byte(configs.Env.JwtSecret), nil
		})

		if err != nil || !rt.Valid {
			ClearTokenCookies(c)
			return c.Redirect().To("/login")
		}

		claims, ok := rt.Claims.(jwt.MapClaims)
		if !ok {
			ClearTokenCookies(c)
			return c.Redirect().To("/login")
		}

		userID := claims["userId"]
		isRemember, _ := claims["isRemember"].(bool)

		res, err := userSvc.Refresh(userID, isRemember, c)
		if err != nil {
			ClearTokenCookies(c)
			return c.Redirect().To("/login")
		}

		setTokenCookies(c, res.AccessToken, res.RefreshToken, isRemember)

		c.Locals("userID", claims["userId"])
		c.Locals("username", claims["sub"])
		return c.Next()
	}
}

func GuestOnly() fiber.Handler {
	return func(c fiber.Ctx) error {
		accessToken := c.Cookies("access_token")
		refreshToken := c.Cookies("refresh_token")

		if accessToken != "" {
			token, err := jwt.Parse(accessToken, func(token *jwt.Token) (any, error) {
				return []byte(configs.Env.JwtSecret), nil
			})
			if err == nil && token.Valid {
				return c.Redirect().To("/apps")
			}
		}

		if refreshToken != "" {
			rt, err := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
				return []byte(configs.Env.JwtSecret), nil
			})
			if err == nil && rt.Valid {
				return c.Redirect().To("/apps")
			}
		}

		return c.Next()
	}
}

func ClearTokenCookies(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Path:     "/",
	})
}

func setTokenCookies(c fiber.Ctx, accessToken, refreshToken string, isRemember bool) {
	atMaxAge := 5 * 60
	rtMaxAge := int((24 * time.Hour).Seconds())
	if isRemember {
		rtMaxAge = int((30 * 24 * time.Hour).Seconds())
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		MaxAge:   atMaxAge,
		Path:     "/",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		MaxAge:   rtMaxAge,
		Path:     "/",
	})
}
