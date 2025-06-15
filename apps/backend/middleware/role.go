package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		if claims["role"] != "admin" {
			return c.Status(403).JSON(fiber.Map{"error": "Admin only"})
		}
		return c.Next()
	}
}
