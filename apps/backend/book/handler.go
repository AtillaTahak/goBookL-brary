// book/handler.go
package book

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetBooks godoc
// @Summary      Get all books
// @Tags         books
// @Produce      json
// @Success      200 {array} Book
// @Router       /books [get]
func GetBooks(c *fiber.Ctx) error {
	return c.JSON(GetAllBooks())
}

// GetBook godoc
// @Summary      Get a single book by ID
// @Tags         books
// @Produce      json
// @Param        id   path  int  true  "Book ID"
// @Success      200  {object} Book
// @Failure      404
// @Router       /books/{id} [get]
func GetBook(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	b := GetBookByID(id)
	if b == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Book not found"})
	}
	return c.JSON(b)
}

// AddBook godoc
// @Summary      Create a new book
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book  body  Book  true  "Book to add"
// @Success      201  {object} Book
// @Failure      400
// @Router       /books [post]
func AddBookHandler(c *fiber.Ctx) error {
	var b Book
	if err := c.BodyParser(&b); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	return c.Status(201).JSON(AddBook(b))
}

// UpdateBook godoc
// @Summary      Update a book by ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id    path  int   true  "Book ID"
// @Param        book  body  Book  true  "Updated book"
// @Success      200   {object} Book
// @Failure      404
// @Router       /books/{id} [put]
func UpdateBookHandler(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var b Book
	if err := c.BodyParser(&b); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}
	updated := UpdateBook(id, b)
	if updated == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Book not found"})
	}
	return c.JSON(updated)
}

// DeleteBook godoc
// @Summary      Delete a book by ID
// @Tags         books
// @Param        id   path  int  true  "Book ID"
// @Success      204
// @Failure      404
// @Router       /books/{id} [delete]
func DeleteBookHandler(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	if !DeleteBook(id) {
		return c.Status(404).JSON(fiber.Map{"error": "Book not found"})
	}
	return c.SendStatus(204)
}
