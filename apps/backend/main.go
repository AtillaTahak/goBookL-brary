package main

import (
	    "log"
    "github.com/gofiber/fiber/v2"
    "github.com/swaggo/fiber-swagger"
    _ "github.com/AtillaTahaK/gobooklibrary/docs"
    "github.com/AtillaTahaK/gobooklibrary/book"
)

// @title           ByFood Book API
// @version         1.0
// @description     API for managing a book library and URL cleaner service
// @host            localhost:8080
// @BasePath        /
func main() {
    app := fiber.New()

    app.Get("/swagger/*", fiberSwagger.WrapHandler)

    app.Get("/books", book.GetBook)
    app.Get("/books/:id", book.GetBook)
    app.Post("/books", book.AddBookHandler)
    app.Put("/books/:id", book.UpdateBookHandler)
    app.Delete("/books/:id", book.DeleteBookHandler)

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Backend is running!")
    })
    log.Println("Server starting on :8080")
    log.Fatal(app.Listen(":8080"))
}
