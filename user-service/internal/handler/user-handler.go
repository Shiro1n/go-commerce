package handler

import (
	"net/http"
	"strconv"

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

func (h *UserHandler) GetUserByUsername(c echo.Context) error {
	username := c.Param("username")
	user, err := h.UserService.GetUserByUsername(username)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetUserById(c echo.Context) error {
	userIdStr := c.Param("user-id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}
	user, err := h.UserService.GetUserById(userId)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Register(c echo.Context) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	// Register user
	user, err := h.UserService.RegisterUser(username, email, password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}
