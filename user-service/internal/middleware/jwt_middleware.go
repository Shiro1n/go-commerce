package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(secret),
	})
}
