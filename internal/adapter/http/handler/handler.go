package handler

import (
	"werewolf-backend/internal/pkg/util"
	"werewolf-backend/internal/port"
)

type Handler struct {
	s          port.Service
	jwtService *util.JWTService
}

func NewHandler(s port.Service, jwtService *util.JWTService) *Handler {
	return &Handler{
		s:          s,
		jwtService: jwtService,
	}
}
