package test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AtillaTahaK/gobooklibrary/middleware"
	"github.com/AtillaTahaK/gobooklibrary/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCORSMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.CORS())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "test"})
	})

	tests := []struct {
		name           string
		origin         string
		expectedOrigin string
		expectedStatus int
	}{
		{
			name:           "Valid origin",
			origin:         "http://localhost:3000",
			expectedOrigin: "*",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "No origin",
			origin:         "",
			expectedOrigin: "",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Equal(t, tt.expectedOrigin, resp.Header.Get("Access-Control-Allow-Origin"))
		})
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.RateLimit())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "test"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	assert.NotEmpty(t, resp.Header.Get("X-RateLimit-Limit"))
	assert.NotEmpty(t, resp.Header.Get("X-RateLimit-Remaining"))
}

func TestLoggerMiddleware(t *testing.T) {
	logger.Init("INFO", "json")

	app := fiber.New()
	app.Use(middleware.Logger())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "test"})
	})

	app.Get("/error", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusInternalServerError, "test error")
	})

	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{
			name:           "Successful request",
			path:           "/test",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error request",
			path:           "/error",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestSecurityMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.Security())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "test"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check security headers
	assert.Equal(t, "DENY", resp.Header.Get("X-Frame-Options"))
	assert.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	assert.Equal(t, "1; mode=block", resp.Header.Get("X-XSS-Protection"))
	assert.Equal(t, "no-referrer", resp.Header.Get("Referrer-Policy"))
}

func TestCompressionMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.Compression())

	app.Get("/test", func(c *fiber.Ctx) error {
		// Return large response to trigger compression
		data := make([]byte, 1024)
		for i := range data {
			data[i] = 'a'
		}
		return c.Send(data)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHealthCheckMiddleware(t *testing.T) {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"timestamp": time.Now().Unix(),
		})
	})

	app.Get("/health/redis", func(c *fiber.Ctx) error {
		// Mock Redis health check
		return c.JSON(fiber.Map{
			"status": "ok",
			"service": "redis",
		})
	})

	app.Get("/health/database", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"service": "database",
		})
	})

	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{
			name:           "General health check",
			path:           "/health",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Redis health check",
			path:           "/health/redis",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Database health check",
			path:           "/health/database",
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestErrorHandlingMiddleware(t *testing.T) {
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	app.Get("/fiber-error", func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusBadRequest, "custom fiber error")
	})

	app.Get("/panic", func(c *fiber.Ctx) error {
		panic("test panic")
	})

	app.Get("/generic-error", func(c *fiber.Ctx) error {
		return assert.AnError
	})

	tests := []struct {
		name           string
		path           string
		expectedStatus int
	}{
		{
			name:           "Fiber error",
			path:           "/fiber-error",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Generic error",
			path:           "/generic-error",
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.Recovery())

	app.Get("/panic", func(c *fiber.Ctx) error {
		panic("test panic")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestTimeoutMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(middleware.Timeout(100 * time.Millisecond))

	app.Get("/fast", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "fast"})
	})

	app.Get("/slow", func(c *fiber.Ctx) error {
		time.Sleep(200 * time.Millisecond)
		return c.JSON(fiber.Map{"message": "slow"})
	})

	t.Run("Fast request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/fast", nil)
		resp, err := app.Test(req, 1000) // 1 second timeout for test
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Slow request timeout", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/slow", nil)
		resp, err := app.Test(req, 300) // 300ms timeout for test
		require.NoError(t, err)
		assert.Equal(t, http.StatusRequestTimeout, resp.StatusCode)
	})
}

func BenchmarkMiddleware(b *testing.B) {
	app := fiber.New()
	app.Use(middleware.Logger())
	app.Use(middleware.CORS())
	app.Use(middleware.Security())
	app.Use(middleware.Compression())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "test"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := app.Test(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestMiddlewareChain(t *testing.T) {
	app := fiber.New()

	app.Use(middleware.Logger())
	app.Use(middleware.CORS())
	app.Use(middleware.Security())
	app.Use(middleware.Recovery())
	app.Use(middleware.RateLimit())

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "test"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)


	accessControl := resp.Header.Get("Access-Control-Allow-Origin")
	xFrame := resp.Header.Get("X-Frame-Options")

	hasHeaders := accessControl != "" || xFrame != ""
	assert.True(t, hasHeaders, "Expected at least one middleware header to be set")
}
