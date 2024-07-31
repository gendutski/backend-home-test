package middleware

import (
	"errors"
	"hometest4/config"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type ValidateJwtFunc func(token string) bool

func GetJWT() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.GetJwtSecret()),
		ContextKey: config.JwtContextKey,
	})
}

func GetUserFromJWT() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// validate token
		return func(c echo.Context) error {
			// get token
			token, ok := c.Get(config.JwtContextKey).(*jwt.Token)
			if !ok {
				c.Error(errors.New("JWT token missing or invalid"))
			}
			// get claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				c.Error(errors.New("failed to cast claims as jwt.MapClaims"))
			}
			// get username
			_username, ok := claims[config.TokenUsernameField]
			if !ok {
				c.Error(errors.New("failed to get username from claims"))
			}
			username, ok := _username.(string)
			if !ok {
				c.Error(errors.New("failed to cast username as string"))
			}

			c.Set(config.UserContextKey, username)
			return next(c)
		}
	}
}
