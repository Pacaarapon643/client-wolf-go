package router

import (
	"werewolf-backend/internal/adapter/http/handler"
	"werewolf-backend/internal/adapter/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func RoomRoutes(v1 fiber.Router, h *handler.Handler, jwtMiddleware *middleware.JWTMiddleware) {
	room := v1.Group("/rooms", jwtMiddleware.MiddlewareAuth())
	room.Post("/create", h.CreateRoom)
	room.Get("/all", h.GetRoom)
	room.Get("/count", h.CountRoom)
	room.Get("/detail", h.FindRoom)
	room.Put("/join", h.JoinRoom)
	room.Post("/join-member", h.JoinRoomMember)
	room.Get("/member", h.GetRoomMember)
	room.Put("/leave", h.LeaveRoom)
	room.Get("/ws", websocket.New(h.RoomWebSocket))

}
