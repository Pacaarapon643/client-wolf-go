package services

import (
	"context"
	"werewolf-backend/internal/adapter/http/dto"
	"werewolf-backend/internal/models"
	"werewolf-backend/internal/pkg/util"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (s Service) CreateRoom(ctx context.Context, obj *models.Room) error {
	err := s.r.CreateRoom(ctx, obj)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) GetRoom(ctx context.Context, obj *[]models.Room) error {
	err := s.r.GetRoom(ctx, obj)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) CountRoom(ctx context.Context, count *int64) error {
	return s.r.CountRoom(ctx, count)
}

func (s Service) FindRoom(ctx context.Context, obj *models.Room) error {
	return s.r.FindRoom(ctx, obj)
}

func (s Service) JoinRoom(ctx context.Context, id string) error {
	return s.r.JoinRoom(ctx, id)
}

func (s Service) JoinRoomMember(ctx context.Context, obj *models.RoomMember, maxRoom int) error {
	return s.r.JoinRoomMember(ctx, obj, maxRoom)
}

func (s Service) GetRoomMember(ctx context.Context, obj *[]dto.RoomMemberResponse, query dto.QueryRoomMember) error {
	err := s.r.GetRoomMember(ctx, obj, query.RoomId)
	if err != nil {
		return err
	}

	room := make(map[int]dto.RoomMemberResponse)
	for _, v := range *obj {
		room[v.SlotIndex] = v
	}

	finalObj := make([]dto.RoomMemberResponse, query.MaxRoom)

	for i := 0; i < query.MaxRoom; i++ {
		if v, ok := room[i]; ok {
			finalObj[i] = v
		} else {
			finalObj[i] = dto.RoomMemberResponse{
				SlotIndex: i,
				EmptySlot: true,
			}
		}
	}

	*obj = finalObj
	return nil
}

func (s Service) LeaveRoom(ctx context.Context, roomId string, userId string) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return &util.LocalizedError{
			Code:    fiber.StatusBadRequest,
			Err:     err.Error(),
			Message: "Invalid user ID format",
		}
	}
	return s.r.LeaveRoom(ctx, roomId, userUUID)
}

func (s Service) LeaveRoomJoin(ctx context.Context, roomId string) error {
	return s.r.LeaveRoomJoin(ctx, roomId)
}

func (s Service) ReadyRoom(ctx context.Context, roomId string, userId string, actionText string) error {

	var action bool
	if actionText == "ready" {
		action = true
	} else {
		action = false
	}

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return &util.LocalizedError{
			Code:    fiber.StatusBadRequest,
			Err:     err.Error(),
			Message: "Invalid user ID format",
		}
	}

	return s.r.ReadyRoom(ctx, roomId, userUUID, action)
}

func (s Service) StartGame(ctx context.Context, roomId string) (bool, error) {

	joinRoom, maxRoom, err := s.r.UpdateJoinRoom(ctx, roomId)
	if err != nil {
		return false, err
	}

	if joinRoom != maxRoom {
		return false, nil
	}

	return true, nil

}

func (s Service) UpdateStartGame(ctx context.Context, roomId string) error {
	return s.r.UpdateStartGame(ctx, roomId)
}
