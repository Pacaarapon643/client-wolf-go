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

func (s Service) CancelVote(ctx context.Context, gameId string, userId string) error {
	jsonData, err := s.rdb.Get(ctx, gameId).Result()
	if err != nil {
		return err
	}

	var arrObjVote []dto.Vote
	if err := json.Unmarshal([]byte(jsonData), &arrObjVote); err != nil {
		return err
	}

	for i, v := range arrObjVote {
		if v.UserID == userId {
			arrObjVote[i].Vote = nil
		}
	}

	jsonDataVote, err := json.Marshal(arrObjVote)
	if err != nil {
		return err
	}

	return s.rdb.Set(ctx, gameId, jsonDataVote, 3600*time.Minute).Err()
}

func (s Service) SummaryVote(ctx context.Context, gameId string, phase string) (*dto.SummaryResult, error) {
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

	log.Printf("[DEBUG] SummaryVote: GameId=%s, Phase=%s, VoteEntries=%d", gameId, phase, len(arrObjVote))
	for _, v := range arrObjVote {
		voteVal := "nil"
		if v.Vote != nil {
			voteVal = strconv.Itoa(*v.Vote)
		}
		log.Printf("[DEBUG]   User=%s Role=%s Vote=%s", v.UserID, v.Role, voteVal)
	}

	var kill []int
	save := -1
	voteMap := make(map[int]int)

	for _, v := range arrObjVote {
		if v.Vote == nil {
			continue
		}

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

	log.Printf("[DEBUG] After compute: kill=%v, save=%d, voteMap=%v", kill, save, voteMap)

	if phase == "night" {
		if len(kill) == 0 {
			return &dto.SummaryResult{
				Result:  "saved",
				Message: "คืนนี้ไม่มีใครถูกล่า ทุกคนปลอดภัย",
			}, nil
		}
		if len(kill) >= 2 {
			return &dto.SummaryResult{
				Result:  "saved",
				Message: "หมาป่าโหวตไม่เป็นเอกฉันท์ ไม่มีใครตาย",
			}, nil
		}
		if kill[0] == save {
			// หาชื่อคนที่ถูกช่วย
			player, _ := s.r.GetPlayerBySlot(ctx, gameId, kill[0])
			name := "ผู้เล่น"
			if player != nil {
				name = player.UserName
			}
			return &dto.SummaryResult{
				Result:  "saved",
				Message: "การ์ดได้ปกป้อง " + name + " ไว้ได้!",
			}, nil
		}

		// หาชื่อคนตาย
		player, _ := s.r.GetPlayerBySlot(ctx, gameId, kill[0])
		name := "ผู้เล่น"
		if player != nil {
			name = player.UserName
		}

		err := s.r.Dead(ctx, &gameId, &kill[0])
		if err != nil {
			return nil, err
		}

		return &dto.SummaryResult{
			Result:   "dead",
			Message:  name + " ถูกหมาป่าขย้ำเมื่อคืนนี้!",
			DeadName: name,
			DeadSlot: kill[0],
		}, nil
	}

	// Day vote phase
	if len(voteMap) == 0 {
		return &dto.SummaryResult{
			Result:  "tie",
			Message: "ไม่มีการโหวต ไม่มีใครถูกประหาร",
		}, nil
	}

	// หาจำนวนคนที่ยังมีชีวิต เพื่อคำนวณเสียงข้างมาก (ครึ่งหนึ่ง)
	aliveWolf, aliveNonWolf, _ := s.r.CountAlive(ctx, gameId)
	aliveTotal := aliveWolf + aliveNonWolf
	requiredVotes := aliveTotal / 2
	if requiredVotes < 1 {
		requiredVotes = 1
	}

	maxVote := 0
	maxSlot := -1
	for slot, count := range voteMap {
		if count > maxVote {
			maxVote = count
			maxSlot = slot
		}
	}

	maxCountPlayers := 0
	for _, count := range voteMap {
		if count == maxVote {
			maxCountPlayers++
		}
	}

	if maxCountPlayers > 1 {
		return &dto.SummaryResult{
			Result:  "tie",
			Message: "โหวตเสมอ! ไม่มีใครถูกประหาร",
		}, nil
	}

	// เช็คว่าโหวตถึงเกณฑ์ขั้นต่ำ (ครึ่งหนึ่งของคนที่ยังมีชีวิต)
	if maxVote < requiredVotes {
		return &dto.SummaryResult{
			Result:  "tie",
			Message: "โหวตไม่ถึงเกณฑ์ขั้นต่ำ (ต้องการ " + strconv.Itoa(requiredVotes) + " เสียง) ไม่มีใครถูกประหาร",
		}, nil
	}

	// หาชื่อคนถูกประหาร
	player, _ := s.r.GetPlayerBySlot(ctx, gameId, maxSlot)
	name := "ผู้เล่น"
	if player != nil {
		name = player.UserName
	}

	err = s.r.Dead(ctx, &gameId, &maxSlot)
	if err != nil {
		return nil, err
	}

	return &dto.SummaryResult{
		Result:   "dead",
		Message:  name + " ถูกชาวบ้านโหวตประหาร!",
		DeadName: name,
		DeadSlot: maxSlot,
	}, nil
}

func (s Service) CheckWinCondition(ctx context.Context, gameId string) (*dto.WinResult, error) {
	werewolfCount, nonWerewolfCount, err := s.r.CountAlive(ctx, gameId)
	if err != nil {
		return nil, err
	}

	log.Printf("Win check: werewolf=%d, non-werewolf=%d", werewolfCount, nonWerewolfCount)

	if werewolfCount == 0 {
		return &dto.WinResult{
			Winner:  "villager",
			Message: "ชาวบ้านชนะ! หมาป่าถูกกำจัดหมดแล้ว!",
		}, nil
	}

	if werewolfCount > nonWerewolfCount {
		return &dto.WinResult{
			Winner:  "werewolf",
			Message: "หมาป่าชนะ! หมาป่ามีจำนวนมากกว่าชาวบ้าน!",
		}, nil
	}

	return nil, nil
}

func (s Service) ResetVotes(ctx context.Context, gameId string) error {
	jsonData, err := s.rdb.Get(ctx, gameId).Result()
	if err != nil {
		return err
	}

	var arrObjVote []dto.Vote
	if err := json.Unmarshal([]byte(jsonData), &arrObjVote); err != nil {
		return err
	}

	for i := range arrObjVote {
		arrObjVote[i].Vote = nil
	}

	jsonDataVote, err := json.Marshal(arrObjVote)
	if err != nil {
		return err
	}

	return s.rdb.Set(ctx, gameId, jsonDataVote, 3600*time.Minute).Err()
}

func (s Service) SeerCheck(ctx context.Context, gameId string, targetUserId string) (string, error) {
	role, err := s.r.GetPlayerRole(ctx, gameId, targetUserId)
	if err != nil {
		return "", err
	}

	if role == "werewolf" {
		return "🐺 คนนี้เป็นหมาป่า!", nil
	}
	return "✅ คนนี้ไม่ใช่หมาป่า", nil
}
