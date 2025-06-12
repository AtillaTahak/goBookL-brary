package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/swaggo/fiber-swagger"
	_ "github.com/AtillaTahaK/gobooklibrary/docs" // Yerel modül adını kullanın

)

// @title           ByFood Book API
// @version         1.0
// @description     API for managing a book library and URL cleaner service
// @host            localhost:8080
// @BasePath        /
func main() {
	app := fiber.New()

	// Swagger endpoint
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Backend working")
	})

	app.Listen(":8080")
}
