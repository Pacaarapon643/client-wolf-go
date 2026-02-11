package port

import (
	"context"
	"werewolf-backend/internal/adapter/http/dto"
	"werewolf-backend/internal/models"

	"github.com/google/uuid"
)

type RoomRepo interface {
	CreateRoom(ctx context.Context, obj *models.Room) error
	GetRoom(ctx context.Context, obj *[]models.Room) error
	CountRoom(ctx context.Context, count *int64) error
	FindRoom(ctx context.Context, obj *models.Room) error
	JoinRoom(ctx context.Context, id string) error
	JoinRoomMember(ctx context.Context, obj *models.RoomMember, maxRoom int) error
	GetRoomMember(ctx context.Context, obj *[]dto.RoomMemberResponse, roomId string) error
	LeaveRoom(ctx context.Context, roomId string, userId uuid.UUID) error
	LeaveRoomJoin(ctx context.Context, roomId string) error
	ReadyRoom(ctx context.Context, roomId string, userId uuid.UUID, action bool) error 
}

type RoomService interface {
	CreateRoom(ctx context.Context, obj *models.Room) error
	GetRoom(ctx context.Context, obj *[]models.Room) error
	CountRoom(ctx context.Context, count *int64) error
	FindRoom(ctx context.Context, obj *models.Room) error
	JoinRoom(ctx context.Context, id string) error
	JoinRoomMember(ctx context.Context, obj *models.RoomMember, maxRoom int) error
	GetRoomMember(ctx context.Context, obj *[]dto.RoomMemberResponse, query dto.QueryRoomMember) error
	LeaveRoom(ctx context.Context, roomId string, userId string) error
	LeaveRoomJoin(ctx context.Context, roomId string) error
	ReadyRoom(ctx context.Context, roomId string, userId string, actionText string) error 
}
