package internal

import (
	"net/http"
	"werewolf-backend/internal/database"
	"werewolf-backend/internal/handler"
	"werewolf-backend/internal/repositories"
	"werewolf-backend/internal/router"
	"werewolf-backend/internal/services"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, db *database.Database) {
	
	r := repositories.NewRepository(db.DB)
	s := services.NewService(r)
	h := handler.NewHandler(s)
	SetupRoutes(app, h)
}

func SetupRoutes(app *fiber.App, h *handler.Handler) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	router.AuthRoutes(v1, h)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"result":  http.StatusText(fiber.StatusNotFound),
			"data":    nil,
			"message": "Route not found",
		})
	})
}
