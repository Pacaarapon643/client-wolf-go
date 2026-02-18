package router

import (
	"werewolf-backend/internal/adapter/http/handler"
	"werewolf-backend/internal/adapter/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func GameRoutes(v1 fiber.Router, h *handler.Handler, jwtMiddleware *middleware.JWTMiddleware) {
	game := v1.Group("/games", jwtMiddleware.MiddlewareAuth())
	game.Get("/test", h.Test123)
	game.Get("/role", h.GetRoleGame)
	game.Get("/game", h.GetGame)
	game.Get("/ws-game", websocket.New(h.GameWebSocket))
}
