package handler

import "werewolf-backend/internal/port"

type Handler struct {
	s port.Service
}

func NewHandler(s port.Service) *Handler {
	return &Handler{
		s: s,
	}
}
