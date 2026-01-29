package router

import (
	"werewolf-backend/internal/adapter/http/handler"
	"werewolf-backend/internal/adapter/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(v1 fiber.Router, h *handler.Handler, jwtMiddleware *middleware.JWTMiddleware) {
	auth := v1.Group("/auths")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/test", jwtMiddleware.MiddlewareAuth(), h.Test)
}
