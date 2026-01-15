package router

import (
	"werewolf-backend/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(v1 fiber.Router, h *handler.Handler) {
	auth := v1.Group("/auths")
	auth.Post("/register", h.Register)
}
