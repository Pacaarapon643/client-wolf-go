package port

import (
	"context"
	"werewolf-backend/internal/models"
)

type RoomRepo interface {
	CreateRoom(ctx context.Context, obj *models.Room) error
	GetRoom(ctx context.Context, obj *[]models.Room) error
	CountRoom(ctx context.Context, count *int64) error
}

type RoomService interface {
	CreateRoom(ctx context.Context, obj *models.Room) error
	GetRoom(ctx context.Context, obj *[]models.Room) error
	CountRoom(ctx context.Context, count *int64) error
}
