package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AtillaTahaK/gobooklibrary/book"
	_ "github.com/AtillaTahaK/gobooklibrary/docs"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title           Book Library API
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

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Println("Server starting on :8080")
		if err := app.Listen(":8080"); err != nil {
			log.Printf("Error starting server: %v", err)
		}
	}()

	<-c
	log.Println("Gracefully shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
