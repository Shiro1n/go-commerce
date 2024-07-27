package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shiro1n/go-commerce/auth-service/internal/config"
	"github.com/shiro1n/go-commerce/auth-service/internal/handler"
	"github.com/shiro1n/go-commerce/auth-service/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	authService := service.NewAuthService(cfg)
	authHandler := handler.NewAuthHandler(authService, cfg)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/login", authHandler.Login)
	e.POST("/register", authHandler.Register)
	e.POST("/refresh", authHandler.RefreshToken)
	e.PUT("/users/:id", authHandler.UpdateUser)

	e.Logger.Fatal(e.Start(":8080"))
}
