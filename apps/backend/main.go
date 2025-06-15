package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AtillaTahaK/gobooklibrary/auth"
	"github.com/AtillaTahaK/gobooklibrary/book"
	_ "github.com/AtillaTahaK/gobooklibrary/docs"
	"github.com/AtillaTahaK/gobooklibrary/middleware"
	"github.com/AtillaTahaK/gobooklibrary/pkg/cache"
	"github.com/AtillaTahaK/gobooklibrary/pkg/db"
	"github.com/AtillaTahaK/gobooklibrary/pkg/logger"
	"github.com/AtillaTahaK/gobooklibrary/pkg/metrics"
	"github.com/AtillaTahaK/gobooklibrary/url"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// Global instances
var (
	AppLogger  *logger.Logger
	RedisCache *cache.RedisCache
)


// @title           Book Library API
// @version         1.0
// @description     A comprehensive REST API for managing a book library with authentication, caching, logging, and metrics
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
    // Load environment variables
    if err := godotenv.Load(".env.local"); err != nil {
        if err := godotenv.Load(); err != nil {
            log.Println("No .env file found, using system environment variables")
        }
    }

    // Initialize logger
    AppLogger = logger.NewLogger()
    AppLogger.Info("🚀 Starting Book Library API...")

    // Initialize Redis cache (with fallback if Redis is not available)
    redisAddr := getEnv("REDIS_URL", "localhost:6379")
    redisPassword := getEnv("REDIS_PASSWORD", "")
    RedisCache = cache.NewRedisCache(redisAddr, redisPassword, 0)
    AppLogger.Info("✅ Redis cache initialized")

    // Set global instances for book package
    book.Cache = RedisCache
    book.Log = AppLogger
    auth.Log = AppLogger

    // Initialize database connection
    db.ConnectDB()
    AppLogger.Info("✅ Database connected")

    // Run auto migrations
    db.AutoMigrate(&auth.User{}, &book.Book{})
    AppLogger.Info("✅ Database migrations completed")

    AppLogger.Info("✅ Database seeded")

    // Create Fiber app
    app := fiber.New(fiber.Config{
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            code := fiber.StatusInternalServerError
            if e, ok := err.(*fiber.Error); ok {
                code = e.Code
            }

            // Log error
            AppLogger.LogError(err, map[string]interface{}{
                "method": c.Method(),
                "path":   c.Path(),
                "ip":     c.IP(),
                "status": code,
            })

            return c.Status(code).JSON(fiber.Map{
                "error": err.Error(),
                "timestamp": time.Now().UTC(),
            })
        },
    })

    // Add middleware
    app.Use(fiberLogger.New(fiberLogger.Config{
        Format: "${time} ${method} ${path} ${status} ${latency} ${ip}\n",
    }))

    app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
    }))

    // Metrics middleware
    app.Use(func(c *fiber.Ctx) error {
        start := time.Now()

        err := c.Next()

        duration := time.Since(start)
        status := c.Response().StatusCode()

        // Record metrics
        metrics.RecordHTTPRequest(
            c.Method(),
            c.Path(),
            fmt.Sprintf("%d", status),
            duration,
        )

        // Log request
        AppLogger.LogRequest(
            c.Method(),
            c.Path(),
            c.IP(),
            c.Get("User-Agent"),
            status,
            duration,
        )

        return err
    })

    // Prometheus metrics endpoint
    app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

    // Swagger documentation
    app.Get("/swagger/*", fiberSwagger.WrapHandler)

    // Health check with detailed status
    app.Get("/health", func(c *fiber.Ctx) error {
        // Check database connection
        sqlDB, err := db.DB.DB()
        if err != nil {
            return c.Status(503).JSON(fiber.Map{
                "status": "unhealthy",
                "database": "disconnected",
                "error": err.Error(),
            })
        }

        // Check Redis connection
        _, err = RedisCache.GetStats()
        redisStatus := "connected"
        if err != nil {
            redisStatus = "disconnected"
        }

        return c.JSON(fiber.Map{
            "status": "healthy",
            "message": "Book Library API is running!",
            "version": "1.0",
            "database": "PostgreSQL with GORM",
            "cache": "Redis",
            "redis_status": redisStatus,
            "connections": sqlDB.Stats(),
            "timestamp": time.Now().UTC(),
        })
    })

    app.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Book Library API",
            "version": "1.0",
            "documentation": "/swagger/",
            "health": "/health",
            "metrics": "/metrics",
        })
    })


    app.Post("/auth/register", auth.Register)
    app.Post("/auth/login", auth.Login)
    app.Post("/url/clean", url.CleanURLHandler)

    app.Get("/books", book.GetBooks)
    app.Get("/books/:id", book.GetBook)


    protected := app.Group("/", middleware.JWTProtected())
    protected.Post("/books", book.AddBookHandler)
    protected.Put("/books/:id", book.UpdateBookHandler)
    protected.Delete("/books/:id", book.DeleteBookHandler)

    admin := protected.Group("/", middleware.RequireAdmin())
    admin.Get("/admin/users", func(c *fiber.Ctx) error {
        var users []auth.User
        result := db.DB.Find(&users)
        if result.Error != nil {
            return c.Status(500).JSON(fiber.Map{
                "error": "Failed to fetch users",
            })
        }

        for i := range users {
            users[i].Password = ""
        }

        return c.JSON(fiber.Map{
            "users": users,
            "total": len(users),
        })
    })

    admin.Get("/admin/stats", func(c *fiber.Ctx) error {
        var bookCount int64
        var userCount int64

        db.DB.Model(&book.Book{}).Count(&bookCount)
        db.DB.Model(&auth.User{}).Count(&userCount)

        // Update metrics
        metrics.SetBooksTotal(float64(bookCount))
        metrics.SetUsersTotal(float64(userCount))

        return c.JSON(fiber.Map{
            "books_total": bookCount,
            "users_total": userCount,
            "timestamp": time.Now().UTC(),
        })
    })

    // Graceful shutdown
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)

    go func() {
        AppLogger.Info("🚀 Server starting on :8080")
        AppLogger.Info("📚 Swagger docs available at http://localhost:8080/swagger/")
        AppLogger.Info("📊 Metrics available at http://localhost:8080/metrics")
        AppLogger.Info("🔍 Health check available at http://localhost:8080/health")

        if err := app.Listen(":8080"); err != nil {
            AppLogger.LogError(err, map[string]interface{}{
                "component": "server",
                "action": "startup",
            })
        }
    }()

    <-c
    AppLogger.Info("🛑 Gracefully shutting down...")

    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    // Close Redis connection
    if RedisCache != nil {
        RedisCache.Close()
        AppLogger.Info("✅ Redis connection closed")
    }

    if err := app.ShutdownWithContext(ctx); err != nil {
        AppLogger.LogError(err, map[string]interface{}{
            "component": "server",
            "action": "shutdown",
        })
    }

    AppLogger.Info("✅ Server exited")
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
