package url

import (
	"github.com/gofiber/fiber/v2"
)

// CleanURLHandler godoc
// @Summary Clean and redirect URL
// @Tags url
// @Accept json
// @Produce json
// @Param data body URLRequest true "URL cleanup input"
// @Success 200 {object} URLResponse
// @Failure 400
// @Router /url/clean [post]
func CleanURLHandler(c *fiber.Ctx) error {
	var req URLRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	cleaned, err := CleanURL(req.URL, req.Operation)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Processing failed"})
	}

	return c.JSON(URLResponse{ProcessedURL: cleaned})
}
