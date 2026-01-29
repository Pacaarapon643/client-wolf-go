package services

import (
	"werewolf-backend/internal/port"
)

type Service struct {
	r port.Repository
}

func NewService(r port.Repository) *Service {
	return &Service{
		r: r,
	}
}
