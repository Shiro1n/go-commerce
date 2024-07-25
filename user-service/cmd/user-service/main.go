package main

import (
	"log"
	"net/http"

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

	http.HandleFunc("/users", userHandler.GetUsers)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
