package repositories

import (
	"context"
	"log"
	"werewolf-backend/internal/adapter/http/dto"
	"werewolf-backend/internal/models"
	"werewolf-backend/internal/pkg/util"

	"github.com/gofiber/fiber/v2"
)

func (r *Repository) CreateGame(ctx context.Context, game []models.Game) error {
	return r.db.WithContext(ctx).CreateInBatches(&game, 100).Error
}

func (r *Repository) GetGame(ctx context.Context, arr *[]dto.Game, game_id *string) error {
	if game_id == nil {
		return &util.LocalizedError{
			Code:    fiber.StatusBadRequest,
			Message: "game_id is required",
		}
	}
	return r.db.WithContext(ctx).Where("game_id = ?", *game_id).Order("slot_index ASC").Find(arr).Error
}

func (r *Repository) JoinGame(ctx context.Context, gameId *string, userId *string) error {
	log.Println("gameId", *gameId)
	log.Println("userId", *userId)
	if gameId == nil || userId == nil {
		return &util.LocalizedError{
			Code:    fiber.StatusBadRequest,
			Message: "game_id and user_id is required",
		}
	}
	return r.db.WithContext(ctx).Table(models.Game{}.TableName()).Where("game_id = ? AND user_id = ?", *gameId, *userId).Update("is_join", true).Error
}

func (r *Repository) LeaveGame(ctx context.Context, gameId *string, userId *string) error {
	if gameId == nil || userId == nil {
		return &util.LocalizedError{
			Code:    fiber.StatusBadRequest,
			Message: "game_id and user_id is required",
		}
	}
	return r.db.WithContext(ctx).Table(models.Game{}.TableName()).Where("game_id = ? AND user_id = ?", *gameId, *userId).Update("is_join", false).Error
}
