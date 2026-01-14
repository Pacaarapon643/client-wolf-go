package router

import (
	"werewolf-backend/internal/config"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, cfg *config.Config) {
	// Health check

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Route not found",
		})
	})
}
