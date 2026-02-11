package services

import (
	"context"
	"encoding/json"
	"math/rand/v2"
	"time"
	"werewolf-backend/internal/adapter/http/dto"
	"werewolf-backend/internal/pkg/appconst"
)

func (s Service) RoleDistribution(ctx context.Context, roomId string) error {
	var obj []dto.RoomMemberResponse
	err := s.r.GetRoomMember(ctx, &obj, roomId)
	if err != nil {
		return err
	}

	var roles []string
	if len(obj) == 6 {
		roles = []string{appconst.RoleWerewolf, appconst.RoleWerewolf, appconst.RoleSeer, appconst.RoleGuard, appconst.RoleVillager, appconst.RoleVillager}
	}

	rand.Shuffle(len(roles), func(i, j int) {
		roles[i], roles[j] = roles[j], roles[i]
	})

	gameData := make(map[string]string)
	for i, role := range roles {
		gameData[obj[i].UserId] = role
	}

	jsonData, err := json.Marshal(gameData)
	if err != nil {
		return err
	}

	err = s.rdb.Set(ctx, roomId, jsonData, 30*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}
