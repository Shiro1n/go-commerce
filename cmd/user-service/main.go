package main

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shiro1n/go-commerce/internal/user/handler"
	"github.com/shiro1n/go-commerce/pkg/database"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Go Commerce API
// @version 1.0
// @description This is a sample server for a Go Commerce application.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /

func main() {
	// Load environment variables
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://mongo:27017"
	}

	// Connect to MongoDB
	database.ConnectMongoDB(mongoURI)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", handler.Hello)
	e.POST("/users", handler.CreateUser)
	e.GET("/users/:email", handler.GetUserByEmail)

	// Swagger
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
