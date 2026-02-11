package services

import (
	"werewolf-backend/internal/infrastructure/database"
	"werewolf-backend/internal/port"
)

type Service struct {
	r port.Repository
	rdb *database.RedisClient
}

func NewService(r port.Repository, rdb *database.RedisClient) *Service {
	return &Service{
		r:   r,
		rdb: rdb,
	}
}
