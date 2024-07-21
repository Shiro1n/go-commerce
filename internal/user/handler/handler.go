package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
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
