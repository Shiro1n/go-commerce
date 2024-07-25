package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shiro1n/go-commerce/user-service/internal/config"
	"github.com/shiro1n/go-commerce/user-service/internal/database"
	"github.com/shiro1n/go-commerce/user-service/internal/handler"
	"github.com/shiro1n/go-commerce/user-service/internal/repository"
	"github.com/shiro1n/go-commerce/user-service/internal/service"
)

func main() {
	cfg := config.LoadConfig()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/users", userHandler.GetUsers)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
