package router

import (
	"werewolf-backend/internal/config"
	"werewolf-backend/internal/database"
	"werewolf-backend/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, cfg *config.Config, db *database.Database) {
	// Health check
	app.Get("/health", handler.HealthCheck(cfg, db))

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Route not found",
		})
	})
}
