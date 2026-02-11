package router

import (
	"werewolf-backend/internal/adapter/http/handler"
	"werewolf-backend/internal/adapter/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func GameRoutes(v1 fiber.Router, h *handler.Handler, jwtMiddleware *middleware.JWTMiddleware) {
	game := v1.Group("/games")
	game.Get("/test", h.Test123)
}
