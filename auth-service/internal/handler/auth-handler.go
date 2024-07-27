package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/shiro1n/go-commerce/auth-service/internal/config"
	"github.com/shiro1n/go-commerce/auth-service/internal/model"
	"github.com/shiro1n/go-commerce/auth-service/internal/service"
)

type AuthHandler struct {
	AuthService service.AuthService
	Config      config.Config
}

func NewAuthHandler(authService service.AuthService, cfg config.Config) *AuthHandler {
	return &AuthHandler{AuthService: authService, Config: cfg}
}

func (h *AuthHandler) Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := h.AuthService.AuthenticateUser(username, password)
	if err != nil {
		return echo.ErrUnauthorized
	}

	// Store user in Redis
	err = h.AuthService.StoreUserInRedis(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to store user in Redis"})
	}

	tokens, err := h.AuthService.CreateTokens(user.ID)
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Expires:  time.Unix(tokens.RtExpires, 0),
		HttpOnly: true,
		Secure:   true,
	})

	return c.JSON(http.StatusOK, map[string]string{"access_token": tokens.AccessToken})
}

func (h *AuthHandler) Register(c echo.Context) error {
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	user, err := h.AuthService.RegisterUser(username, email, password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Store user in Redis
	err = h.AuthService.StoreUserInRedis(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to store user in Redis"})
	}

	tokens, err := h.AuthService.CreateTokens(user.ID)
	if err != nil {
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Expires:  time.Unix(tokens.RtExpires, 0),
		HttpOnly: true,
		Secure:   true,
	})

	return c.JSON(http.StatusCreated, map[string]string{"access_token": tokens.AccessToken})
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		return echo.ErrUnauthorized
	}

	tokens, err := h.AuthService.RefreshTokens(cookie.Value)
	if err != nil {
		return echo.ErrUnauthorized
	}

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Expires:  time.Unix(tokens.RtExpires, 0),
		HttpOnly: true,
		Secure:   true,
	})

	return c.JSON(http.StatusOK, map[string]string{"access_token": tokens.AccessToken})
}

func (h *AuthHandler) UpdateUser(c echo.Context) error {
	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	// Update user in primary database and Redis
	if err := h.AuthService.UpdateUser(&user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}
