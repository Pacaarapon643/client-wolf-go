package router

import (
	"werewolf-backend/internal/adapter/http/handler"
	"werewolf-backend/internal/adapter/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func RoomRoutes(v1 fiber.Router, h *handler.Handler, jwtMiddleware *middleware.JWTMiddleware) {
	room := v1.Group("/rooms")
	room.Post("/create", h.CreateRoom)
	room.Get("/all", h.GetRoom)
	room.Get("/count", h.CountRoom)
}
