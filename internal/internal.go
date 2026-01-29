package internal

import (
	"net/http"
	"werewolf-backend/internal/adapter/http/handler"
	"werewolf-backend/internal/adapter/http/middleware"
	"werewolf-backend/internal/adapter/http/router"
	"werewolf-backend/internal/adapter/repositories"
	"werewolf-backend/internal/adapter/services"
	"werewolf-backend/internal/infrastructure/database"
	"werewolf-backend/internal/pkg/util"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, db *database.Database, jwtMiddleware *middleware.JWTMiddleware, jwtService *util.JWTService) {
	r := repositories.NewRepository(db.DB)
	s := services.NewService(r)
	h := handler.NewHandler(s, jwtService)
	SetupRoutes(app, h, jwtMiddleware)
}

func SetupRoutes(app *fiber.App, h *handler.Handler, jwtMiddleware *middleware.JWTMiddleware) {
	api := app.Group("/api")
	v1 := api.Group("/v1")

	router.AuthRoutes(v1, h, jwtMiddleware)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"result":  http.StatusText(fiber.StatusNotFound),
			"data":    nil,
			"message": "Route not found",
		})
	})
}
