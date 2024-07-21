package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shiro1n/go-commerce/internal/user/model"
	"github.com/shiro1n/go-commerce/internal/user/repository"
)

// Hello example
// @Summary Show a hello message
// @Description Get a string message
// @ID get-hello
// @Produce json
// @Success 200 {string} string "Hello, User!"
// @Router / [get]
func Hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, User!")
}

// CreateUser creates a new user
func CreateUser(c echo.Context) error {
	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := repository.CreateUser(ctx, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, result)
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(c echo.Context) error {
	email := c.Param("email")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := repository.GetUserByEmail(ctx, email)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}
