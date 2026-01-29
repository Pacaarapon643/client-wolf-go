package handler

import (
	"time"
	"werewolf-backend/internal/infrastructure/config"
	"werewolf-backend/internal/infrastructure/database"

	"github.com/gofiber/fiber/v2"
)

func HealthCheck(cfg *config.Config, db *database.Database) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check database health
		dbStatus := "connected"
		if err := db.Health(c.Context()); err != nil {
			dbStatus = "disconnected"
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":      "unhealthy",
				"database":    dbStatus,
				"environment": cfg.Server.Environment,
				"timestamp":   time.Now().Format(time.RFC3339),
			})
		}

		return c.JSON(fiber.Map{
			"status":      "healthy",
			"database":    dbStatus,
			"db_stats":    db.GetStats(),
			"environment": cfg.Server.Environment,
			"message":     "Quiz Game API is running",
			"timestamp":   time.Now().Format(time.RFC3339),
		})
	}
}
