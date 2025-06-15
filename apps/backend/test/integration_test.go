package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/AtillaTahaK/gobooklibrary/auth"
	"github.com/AtillaTahaK/gobooklibrary/book"
	"github.com/AtillaTahaK/gobooklibrary/middleware"
	"github.com/AtillaTahaK/gobooklibrary/pkg/cache"
	"github.com/AtillaTahaK/gobooklibrary/pkg/db"
	"github.com/AtillaTahaK/gobooklibrary/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/suite"
)

type BookAPITestSuite struct {
	suite.Suite
	app    *fiber.App
	cache  *cache.RedisCache
	logger *logger.Logger
	token  string
}

func (suite *BookAPITestSuite) SetupSuite() {
	// Setup test environment
	os.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/gobooklibrary_test?sslmode=disable")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("REDIS_URL", "localhost:6379")

	// Initialize logger
	suite.logger = logger.NewLogger()
	suite.logger.SetLevel(logger.DEBUG)

	// Initialize cache
	suite.cache = cache.NewRedisCache("localhost:6379", "", 2) // Use DB 2 for testing

	// Set global instances
	book.Cache = suite.cache
	book.Log = suite.logger
	auth.Log = suite.logger

	// Connect to test database
	db.ConnectDB()
	db.AutoMigrate(&auth.User{}, &book.Book{})

	// Setup Fiber app
	suite.app = fiber.New()

	// Setup routes
	suite.setupRoutes()

	// Create test user and get token
	suite.setupTestUser()
}

func (suite *BookAPITestSuite) TearDownSuite() {
	// Clean up test data
	if suite.cache != nil {
		suite.cache.FlushAll()
		suite.cache.Close()
	}

	// Clean up database
	db.DB.Exec("DELETE FROM books")
	db.DB.Exec("DELETE FROM users")
}

func (suite *BookAPITestSuite) SetupTest() {
	// Clean up books before each test
	db.DB.Exec("DELETE FROM books")

	// Clear cache
	if suite.cache != nil {
		suite.cache.FlushAll()
	}
}

func (suite *BookAPITestSuite) setupRoutes() {
	// Public routes
	suite.app.Post("/auth/register", auth.Register)
	suite.app.Post("/auth/login", auth.Login)
	suite.app.Get("/books", book.GetBooks)
	suite.app.Get("/books/:id", book.GetBook)

	// Protected routes
	protected := suite.app.Group("/", middleware.JWTProtected())
	protected.Post("/books", book.AddBookHandler)
	protected.Put("/books/:id", book.UpdateBookHandler)
	protected.Delete("/books/:id", book.DeleteBookHandler)
}

func (suite *BookAPITestSuite) setupTestUser() {
	// Create test user
	registerReq := auth.RegisterRequest{
		Username: "testuser",
		Password: "testpass123",
		Email:    "test@example.com",
	}

	registerBody, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewReader(registerBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := suite.app.Test(req)
	resp.Body.Close()

	// Login to get token
	loginReq := auth.LoginRequest{
		Username: "testuser",
		Password: "testpass123",
	}

	loginBody, _ := json.Marshal(loginReq)
	req = httptest.NewRequest("POST", "/auth/login", bytes.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ = suite.app.Test(req)
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var loginResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&loginResp)
		suite.token = loginResp["token"].(string)
	}
}

func (suite *BookAPITestSuite) TestGetBooks_Empty() {
	req := httptest.NewRequest("GET", "/books", nil)
	resp, err := suite.app.Test(req)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode)

	var books []book.Book
	json.NewDecoder(resp.Body).Decode(&books)
	suite.Equal(0, len(books))
}

func (suite *BookAPITestSuite) TestAddBook_Success() {
	if suite.token == "" {
		suite.T().Skip("No auth token available")
	}

	newBook := book.Book{
		Title:  "Test Book",
		Author: "Test Author",
		Year:   2023,
		Genre:  "Fiction",
	}

	bookBody, _ := json.Marshal(newBook)
	req := httptest.NewRequest("POST", "/books", bytes.NewReader(bookBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(201, resp.StatusCode)

	var createdBook book.Book
	json.NewDecoder(resp.Body).Decode(&createdBook)
	suite.Equal("Test Book", createdBook.Title)
	suite.Equal("Test Author", createdBook.Author)
	suite.NotZero(createdBook.ID)
}

func (suite *BookAPITestSuite) TestAddBook_Unauthorized() {
	newBook := book.Book{
		Title:  "Test Book",
		Author: "Test Author",
	}

	bookBody, _ := json.Marshal(newBook)
	req := httptest.NewRequest("POST", "/books", bytes.NewReader(bookBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(401, resp.StatusCode)
}

func (suite *BookAPITestSuite) TestGetBook_ById() {
	// First create a book
	testBook := suite.createTestBook()

	// Now get it by ID
	req := httptest.NewRequest("GET", fmt.Sprintf("/books/%d", testBook.ID), nil)
	resp, err := suite.app.Test(req)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode)

	var retrievedBook book.Book
	json.NewDecoder(resp.Body).Decode(&retrievedBook)
	suite.Equal(testBook.ID, retrievedBook.ID)
	suite.Equal(testBook.Title, retrievedBook.Title)
}

func (suite *BookAPITestSuite) TestGetBook_NotFound() {
	req := httptest.NewRequest("GET", "/books/99999", nil)
	resp, err := suite.app.Test(req)

	suite.NoError(err)
	suite.Equal(404, resp.StatusCode)
}

func (suite *BookAPITestSuite) TestUpdateBook_Success() {
	if suite.token == "" {
		suite.T().Skip("No auth token available")
	}

	// Create a book first
	testBook := suite.createTestBook()

	// Update it
	updatedBook := book.Book{
		Title:  "Updated Title",
		Author: "Updated Author",
		Year:   2024,
		Genre:  "Non-Fiction",
	}

	bookBody, _ := json.Marshal(updatedBook)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/books/%d", testBook.ID), bytes.NewReader(bookBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(200, resp.StatusCode)

	var result book.Book
	json.NewDecoder(resp.Body).Decode(&result)
	suite.Equal("Updated Title", result.Title)
	suite.Equal("Updated Author", result.Author)
}

func (suite *BookAPITestSuite) TestDeleteBook_Success() {
	if suite.token == "" {
		suite.T().Skip("No auth token available")
	}

	// Create a book first
	testBook := suite.createTestBook()

	// Delete it
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/books/%d", testBook.ID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)

	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(204, resp.StatusCode)

	// Verify it's gone
	req = httptest.NewRequest("GET", fmt.Sprintf("/books/%d", testBook.ID), nil)
	resp, err = suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(404, resp.StatusCode)
}

func (suite *BookAPITestSuite) TestSearchBooks() {
	// Create some test books
	books := []book.Book{
		{Title: "Go Programming", Author: "John Doe", Year: 2020},
		{Title: "JavaScript Guide", Author: "Jane Smith", Year: 2021},
		{Title: "Python Basics", Author: "John Doe", Year: 2022},
	}

	for _, b := range books {
		suite.createBookInDB(b)
	}

	// Search by title
	req := httptest.NewRequest("GET", "/books?search=Go", nil)
	resp, err := suite.app.Test(req)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode)

	var results []book.Book
	json.NewDecoder(resp.Body).Decode(&results)
	suite.Len(results, 1)
	suite.Equal("Go Programming", results[0].Title)
}

func (suite *BookAPITestSuite) TestCacheIntegration() {
	if suite.cache == nil {
		suite.T().Skip("Cache not available")
	}

	// Create a test book
	testBook := suite.createTestBook()

	// First request should miss cache and hit database
	req := httptest.NewRequest("GET", fmt.Sprintf("/books/%d", testBook.ID), nil)
	start := time.Now()
	resp, err := suite.app.Test(req)
	firstDuration := time.Since(start)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode)
	resp.Body.Close()

	// Second request should hit cache and be faster
	req = httptest.NewRequest("GET", fmt.Sprintf("/books/%d", testBook.ID), nil)
	start = time.Now()
	resp, err = suite.app.Test(req)
	secondDuration := time.Since(start)

	suite.NoError(err)
	suite.Equal(200, resp.StatusCode)
	resp.Body.Close()

	// Cache hit should generally be faster (though not guaranteed in tests)
	suite.T().Logf("First request: %v, Second request: %v", firstDuration, secondDuration)
}

func (suite *BookAPITestSuite) TestInvalidJSON() {
	if suite.token == "" {
		suite.T().Skip("No auth token available")
	}

	req := httptest.NewRequest("POST", "/books", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	resp, err := suite.app.Test(req)
	suite.NoError(err)
	suite.Equal(400, resp.StatusCode)
}

func (suite *BookAPITestSuite) TestInvalidBookID() {
	req := httptest.NewRequest("GET", "/books/invalid", nil)
	resp, err := suite.app.Test(req)

	suite.NoError(err)
	suite.Equal(400, resp.StatusCode)
}

// Helper methods
func (suite *BookAPITestSuite) createTestBook() book.Book {
		if suite.token == "" {
		// Create directly in database if no token
		return suite.createBookInDB(book.Book{
			Title:  "Test Book",
			Author: "Test Author",
			Year:   2023,
			Genre:  "Fiction",
		})
	}

	newBook := book.Book{
		Title:  "Test Book",
		Author: "Test Author",
		Year:   2023,
		Genre:  "Fiction",
	}

	bookBody, _ := json.Marshal(newBook)
	req := httptest.NewRequest("POST", "/books", bytes.NewReader(bookBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)

	resp, _ := suite.app.Test(req)
	defer resp.Body.Close()

	var createdBook book.Book
	json.NewDecoder(resp.Body).Decode(&createdBook)
	return createdBook
}

func (suite *BookAPITestSuite) createBookInDB(b book.Book) book.Book {
	db.DB.Create(&b)
	return b
}

// Benchmark tests
func BenchmarkGetBooks(b *testing.B) {
	// Setup
	suite := new(BookAPITestSuite)
	suite.SetupSuite()
	defer suite.TearDownSuite()

	// Create some test data
	for i := 0; i < 100; i++ {
		suite.createBookInDB(book.Book{
			Title:  fmt.Sprintf("Book %d", i),
			Author: fmt.Sprintf("Author %d", i),
			Year:   2020 + (i % 5),
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/books", nil)
		resp, _ := suite.app.Test(req)
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
}

func TestBookAPITestSuite(t *testing.T) {
	suite.Run(t, new(BookAPITestSuite))
}
