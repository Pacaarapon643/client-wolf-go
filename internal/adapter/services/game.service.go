package services

import (
	"context"
	"encoding/json"
	"log"
	"math/rand/v2"
	"strconv"
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
	var arrObjVote []dto.Vote
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

		arrObjVote = append(arrObjVote, dto.Vote{
			UserID: v.UserId,
			Vote:   nil,
			Role:   gameData[v.UserId],
		})

	}

	err = s.r.CreateGame(ctx, game)
	if err != nil {
		return nil, err
	}

	jsonDataRole, err := json.Marshal(gameData)
	if err != nil {
		return nil, err
	}

	jsonDataVote, err := json.Marshal(arrObjVote)
	if err != nil {
		return nil, err
	}

	err = s.rdb.Set(ctx, roomId, jsonDataRole, 3600*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	err = s.rdb.Set(ctx, gameId, jsonDataVote, 3600*time.Minute).Err()
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

func (s Service) GetGame(ctx context.Context, gameId string, role string) (any, error) {
	var arrGame []dto.Game
	err := s.r.GetGame(ctx, &arrGame, &gameId)
	if err != nil {
		return nil, err
	}

	if role == "werewolf" {
		var arrGameWolf []dto.GameResponse
		for _, v := range arrGame {
			if v.Role == "werewolf" {
				arrGameWolf = append(arrGameWolf, dto.GameResponse{
					GameId:     v.GameId,
					UserId:     v.UserId,
					UserName:   v.UserName,
					SlotIndex:  v.SlotIndex,
					IsDead:     v.IsDead,
					Img:        v.Img,
					IsJoin:     v.IsJoin,
					IsWerewolf: true,
				})
			} else {
				arrGameWolf = append(arrGameWolf, dto.GameResponse{
					GameId:     v.GameId,
					UserId:     v.UserId,
					UserName:   v.UserName,
					SlotIndex:  v.SlotIndex,
					IsDead:     v.IsDead,
					Img:        v.Img,
					IsJoin:     v.IsJoin,
					IsWerewolf: false,
				})
			}
		}
		return arrGameWolf, nil
	} else {
		var arrGameWolf []dto.GameResponse
		for _, v := range arrGame {
			arrGameWolf = append(arrGameWolf, dto.GameResponse{
				GameId:     v.GameId,
				UserId:     v.UserId,
				UserName:   v.UserName,
				SlotIndex:  v.SlotIndex,
				IsDead:     v.IsDead,
				Img:        v.Img,
				IsJoin:     v.IsJoin,
				IsWerewolf: false,
			})

		}
		return arrGameWolf, nil
	}

}

func (s Service) JoinGame(ctx context.Context, gameId string, userId string) error {
	return s.r.JoinGame(ctx, &gameId, &userId)
}

func (s Service) LeaveGame(ctx context.Context, gameId string, userId string) error {
	return s.r.LeaveGame(ctx, &gameId, &userId)
}

func (s Service) Vote(ctx context.Context, gameId string, userId string, vote string, role string) (any, error) {
	jsonData, err := s.rdb.Get(ctx, gameId).Result()
	if err != nil {
		return nil, err
	}

	var arrObjVote []dto.Vote
	if err := json.Unmarshal([]byte(jsonData), &arrObjVote); err != nil {
		return nil, err
	}

	intVote, err := strconv.Atoi(vote)
	if err != nil {
		return nil, err
	}

	for i, v := range arrObjVote {
		if v.UserID == userId {
			arrObjVote[i].Vote = &intVote
		}
	}

	jsonDataVote, err := json.Marshal(arrObjVote)
	if err != nil {
		return nil, err
	}

	err = s.rdb.Set(ctx, gameId, jsonDataVote, 3600*time.Minute).Err()
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s Service) SummaryVote(ctx context.Context, gameId string, phase string) (any, error) {
	log.Println("SummaryVote")
	log.Println("gameId", gameId)
	log.Println("phase", phase)

	jsonData, err := s.rdb.Get(ctx, gameId).Result()
	if err != nil {
		return nil, err
	}

	var arrObjVote []dto.Vote
	if err := json.Unmarshal([]byte(jsonData), &arrObjVote); err != nil {
		return nil, err
	}

	var kill []int
	save := -1 // เปลี่ยนจาก 0 เป็น -1 เพื่อป้องกันการชนกับ Index 0 ของผู้เล่น
	voteMap := make(map[int]int)

	for _, v := range arrObjVote {
		if v.Vote == nil {
			continue
		} // ป้องกัน nil pointer

		if phase == "night" {
			if v.Role == "werewolf" {
				isExist := false
				for _, k := range kill {
					if k == *v.Vote {
						isExist = true
						break
					}
				}
				if !isExist {
					kill = append(kill, *v.Vote)
				}
			} else if v.Role == "guard" {
				save = *v.Vote
			}
		} else {
			voteMap[*v.Vote]++
		}
	}

	if phase == "night" {
		if len(kill) == 0 {
			return "รอด (ไม่มีใครถูกล่า)", nil
		} // เพิ่มการเช็คป้องกัน panic
		if len(kill) >= 2 {
			return "รอด (หมาป่าโหวตไม่เป็นเอกฉันท์)", nil
		}
		if kill[0] == save {
			return "รอด (Guard ช่วยไว้)", nil
		}
		err := s.r.Dead(ctx, &gameId, &kill[0])
		if err != nil {
			return nil, err
		}
		return "ตาย", nil
	} else {
		if len(voteMap) == 0 {
			return "รอด (ไม่มีการโหวต)", nil
		} // เช็คกรณีไม่มีคนโหวตเลย

		maxVote := 0
		for _, count := range voteMap {
			if count > maxVote {
				maxVote = count
			}
		}

		maxCountPlayers := 0
		for _, count := range voteMap {
			if count == maxVote {
				maxCountPlayers++
			}
		}

		if maxCountPlayers > 1 {
			return "เสมอ (ไม่มีใครตาย)", nil
		}
		return "ตาย", nil
	}
}
