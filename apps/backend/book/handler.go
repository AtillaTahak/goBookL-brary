package book

import (
	"fmt"
	"strconv"
	"time"

	"github.com/AtillaTahaK/gobooklibrary/pkg/cache"
	"github.com/AtillaTahaK/gobooklibrary/pkg/logger"
	"github.com/AtillaTahaK/gobooklibrary/pkg/metrics"
	"github.com/gofiber/fiber/v2"
)

var (
	Cache *cache.RedisCache
	Log   *logger.Logger
)

// GetBooks godoc
// @Summary      Get all books
// @Tags         books
// @Produce      json
// @Param        search query string false "Search books by title or author"
// @Success      200 {array} Book
// @Failure      500 {object} map[string]interface{}
// @Router       /books [get]
func GetBooks(c *fiber.Ctx) error {
	start := time.Now()
	search := c.Query("search")

	// Generate cache key
	cacheKey := "books:all"
	if search != "" {
		cacheKey = fmt.Sprintf("books:search:%s", search)
	}

	var books []Book
	var err error

	if Cache != nil {
		err = Cache.Get(cacheKey, &books)
		if err == nil {
			metrics.RecordCacheOperation("get", "hit")
			if Log != nil {
				Log.LogCache("get", cacheKey, true, time.Since(start))
			}
			return c.JSON(books)
		}
		metrics.RecordCacheOperation("get", "miss")
	}

	if search != "" {
		books, err = SearchBooks(search)
	} else {
		books, err = GetAllBooks()
	}

	if err != nil {
		if Log != nil {
			Log.LogError(err, map[string]interface{}{
				"operation": "get_books",
				"search":    search,
			})
		}
		metrics.RecordDatabaseQuery("select", "books", "error", time.Since(start))
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch books"})
	}

	if Cache != nil {
		Cache.Set(cacheKey, books, 5*time.Minute)
		metrics.RecordCacheOperation("set", "success")
	}

	if Log != nil {
		Log.LogDatabase("select", "books", time.Since(start), int64(len(books)))
	}
	metrics.RecordDatabaseQuery("select", "books", "success", time.Since(start))

	return c.JSON(books)
}

// GetBook godoc
// @Summary      Get a single book by ID
// @Tags         books
// @Produce      json
// @Param        id   path  int  true  "Book ID"
// @Success      200  {object} Book
// @Failure      400  {object} map[string]interface{}
// @Failure      404  {object} map[string]interface{}
// @Router       /books/{id} [get]
func GetBook(c *fiber.Ctx) error {
	start := time.Now()
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	cacheKey := fmt.Sprintf("book:%d", id)
	var book Book

	if Cache != nil {
		err = Cache.Get(cacheKey, &book)
		if err == nil {
			metrics.RecordCacheOperation("get", "hit")
			if Log != nil {
				Log.LogCache("get", cacheKey, true, time.Since(start))
			}
			return c.JSON(book)
		}
		metrics.RecordCacheOperation("get", "miss")
	}

	bookPtr, err := GetBookByID(uint(id))
	if err != nil {
		if Log != nil {
			Log.LogError(err, map[string]interface{}{
				"operation": "get_book",
				"book_id":   id,
			})
		}
		metrics.RecordDatabaseQuery("select", "books", "error", time.Since(start))
		return c.Status(404).JSON(fiber.Map{"error": "Book not found"})
	}

	book = *bookPtr

	if Cache != nil {
		Cache.Set(cacheKey, book, 10*time.Minute)
		metrics.RecordCacheOperation("set", "success")
	}

	if Log != nil {
		Log.LogDatabase("select", "books", time.Since(start), 1)
	}
	metrics.RecordDatabaseQuery("select", "books", "success", time.Since(start))

	return c.JSON(book)
}

// AddBook godoc
// @Summary      Create a new book
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book  body  Book  true  "Book to add"
// @Success      201  {object} Book
// @Failure      400  {object} map[string]interface{}
// @Failure      500  {object} map[string]interface{}
// @Router       /books [post]
func AddBookHandler(c *fiber.Ctx) error {
	start := time.Now()
	var book Book
	if err := c.BodyParser(&book); err != nil {
		if Log != nil {
			Log.LogError(err, map[string]interface{}{
				"operation": "add_book",
				"error": "invalid_request_body",
			})
		}
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := CreateBook(&book); err != nil {
		if Log != nil {
			Log.LogError(err, map[string]interface{}{
				"operation": "add_book",
				"title": book.Title,
			})
		}
		metrics.RecordDatabaseQuery("insert", "books", "error", time.Since(start))
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create book"})
	}

	if Cache != nil {
		Cache.Delete("books:all")
		metrics.RecordCacheOperation("delete", "success")
	}

	if Log != nil {
		Log.LogDatabase("insert", "books", time.Since(start), 1)
		Log.LogBookOperation("create", "", book.ID, book.Title)
	}
	metrics.RecordDatabaseQuery("insert", "books", "success", time.Since(start))

	return c.Status(201).JSON(book)
}

// UpdateBook godoc
// @Summary      Update a book by ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id    path  int   true  "Book ID"
// @Param        book  body  Book  true  "Updated book"
// @Success      200   {object} Book
// @Failure      400   {object} map[string]interface{}
// @Failure      404   {object} map[string]interface{}
// @Failure      500   {object} map[string]interface{}
// @Router       /books/{id} [put]
func UpdateBookHandler(c *fiber.Ctx) error {
	start := time.Now()
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	var book Book
	if err := c.BodyParser(&book); err != nil {
		if Log != nil {
			Log.LogError(err, map[string]interface{}{
				"operation": "update_book",
				"book_id": id,
				"error": "invalid_request_body",
			})
		}
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	updatedBook, err := UpdateBook(uint(id), &book)
	if err != nil {
		if Log != nil {
			Log.LogError(err, map[string]interface{}{
				"operation": "update_book",
				"book_id": id,
			})
		}
		metrics.RecordDatabaseQuery("update", "books", "error", time.Since(start))
		return c.Status(404).JSON(fiber.Map{"error": "Book not found"})
	}

	if Cache != nil {
		Cache.Delete("books:all")
		Cache.Delete(fmt.Sprintf("book:%d", id))
		metrics.RecordCacheOperation("delete", "success")
	}

	if Log != nil {
		Log.LogDatabase("update", "books", time.Since(start), 1)
		Log.LogBookOperation("update", "", uint(id), updatedBook.Title)
	}
	metrics.RecordDatabaseQuery("update", "books", "success", time.Since(start))

	return c.JSON(updatedBook)
}

// DeleteBook godoc
// @Summary      Delete a book by ID
// @Tags         books
// @Param        id   path  int  true  "Book ID"
// @Success      204
// @Failure      400  {object} map[string]interface{}
// @Failure      404  {object} map[string]interface{}
// @Router       /books/{id} [delete]
func DeleteBookHandler(c *fiber.Ctx) error {
	start := time.Now()
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	if err := DeleteBook(uint(id)); err != nil {
		if Log != nil {
			Log.LogError(err, map[string]interface{}{
				"operation": "delete_book",
				"book_id": id,
			})
		}
		metrics.RecordDatabaseQuery("delete", "books", "error", time.Since(start))
		return c.Status(404).JSON(fiber.Map{"error": "Book not found"})
	}

	if Cache != nil {
		Cache.Delete("books:all")
		Cache.Delete(fmt.Sprintf("book:%d", id))
		metrics.RecordCacheOperation("delete", "success")
	}

	if Log != nil {
		Log.LogDatabase("delete", "books", time.Since(start), 1)
		Log.LogBookOperation("delete", "", uint(id), "")
	}
	metrics.RecordDatabaseQuery("delete", "books", "success", time.Since(start))

	return c.SendStatus(204)
}
