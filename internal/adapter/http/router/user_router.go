package router

import (
	"werewolf-backend/internal/adapter/http/handler"
	"werewolf-backend/internal/adapter/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(v1 fiber.Router, h *handler.Handler, jwtMiddleware *middleware.JWTMiddleware) {
	user := v1.Group("/users")
	user.Get("/count", h.CountUser)
}
