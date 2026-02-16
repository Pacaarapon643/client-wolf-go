package services

import (
	"context"
	"encoding/json"
	"math/rand/v2"
	"time"
	"werewolf-backend/internal/adapter/http/dto"
	"werewolf-backend/internal/models"
	"werewolf-backend/internal/pkg/appconst"

	"github.com/google/uuid"
)

func (s Service) RoleDistribution(ctx context.Context, roomId string) (*string, error) {
	var obj []dto.RoomMemberResponse
	err := s.r.GetRoomMember(ctx, &obj, roomId)
	if err != nil {
		return nil, err
	}

	var roles []string
	if len(obj) == 6 {
		roles = []string{appconst.RoleWerewolf, appconst.RoleWerewolf, appconst.RoleSeer, appconst.RoleGuard, appconst.RoleVillager, appconst.RoleVillager}
	}

	var img []string
	if len(obj) == 6 {
		img = []string{"avatat-1.png", "avatat-2.png", "avatat-3.png", "avatat-4.png", "avatat-5.png", "avatat-6.png"}
	}

	rand.Shuffle(len(roles), func(i, j int) {
		roles[i], roles[j] = roles[j], roles[i]
	})

	rand.Shuffle(len(img), func(i, j int) {
		img[i], img[j] = img[j], img[i]
	})

	gameData := make(map[string]string)
	for i, role := range roles {
		gameData[obj[i].UserId] = role
	}

	gameId := uuid.New().String()
	var game []models.Game
	for _, v := range obj {
		game = append(game, models.Game{
			GameId:    gameId,
			RoomId:    roomId,
			UserId:    v.UserId,
			UserName:  v.UserName,
			SlotIndex: v.SlotIndex,
			Role:      gameData[v.UserId],
			Img:       img[v.SlotIndex],
		})

	}

	err = s.r.CreateGame(ctx, game)
	if err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(gameData)
	if err != nil {
		return nil, err
	}

	err = s.rdb.Set(ctx, roomId, jsonData, 3600*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	return &gameId, nil
}

func (s Service) GetRoleGame(ctx context.Context, roomId string, userId string) (any, error) {
	jsonData, err := s.rdb.Get(ctx, roomId).Result()
	if err != nil {
		return nil, err
	}

	var obj map[string]string
	if err := json.Unmarshal([]byte(jsonData), &obj); err != nil {
		return nil, err
	}

	return obj[userId], nil
}

func (s Service) GetGame(ctx context.Context, gameId string) (any, error) {
	var arrGame []dto.Game
	err := s.r.GetGame(ctx, &arrGame, &gameId)
	if err != nil {
		return nil, err
	}

	return arrGame, nil
}

func (s Service) JoinGame(ctx context.Context, gameId string, userId string) error {
	return s.r.JoinGame(ctx, &gameId, &userId)
}

func (s Service) LeaveGame(ctx context.Context, gameId string, userId string) error {
	return s.r.LeaveGame(ctx, &gameId, &userId)
}
