package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shiro1n/go-commerce/user-service/internal/service"
)

type UserHandler struct {
	UserService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}
