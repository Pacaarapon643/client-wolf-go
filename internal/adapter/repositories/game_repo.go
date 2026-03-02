package repositories

import (
	"context"
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

func (r *Repository) Dead(ctx context.Context, gameId *string, index *int) error {
	if gameId == nil || index == nil {
		return &util.LocalizedError{
			Code:    fiber.StatusBadRequest,
			Message: "game_id and index is required",
		}
	}

	return r.db.WithContext(ctx).Table(models.Game{}.TableName()).Where("game_id = ? AND slot_index = ?", *gameId, *index).Update("is_dead", true).Error
}

func (r *Repository) CountAlive(ctx context.Context, gameId string) (int, int, error) {
	var werewolfCount int64
	var nonWerewolfCount int64

	err := r.db.WithContext(ctx).Table(models.Game{}.TableName()).
		Where("game_id = ? AND is_dead = ? AND role = ?", gameId, false, "werewolf").
		Count(&werewolfCount).Error
	if err != nil {
		return 0, 0, err
	}

	err = r.db.WithContext(ctx).Table(models.Game{}.TableName()).
		Where("game_id = ? AND is_dead = ? AND role != ?", gameId, false, "werewolf").
		Count(&nonWerewolfCount).Error
	if err != nil {
		return 0, 0, err
	}

	return int(werewolfCount), int(nonWerewolfCount), nil
}

func (r *Repository) GetPlayerBySlot(ctx context.Context, gameId string, slotIndex int) (*models.Game, error) {
	var game models.Game
	err := r.db.WithContext(ctx).Where("game_id = ? AND slot_index = ?", gameId, slotIndex).First(&game).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *Repository) GetPlayerRole(ctx context.Context, gameId string, userId string) (string, error) {
	var game models.Game
	err := r.db.WithContext(ctx).Where("game_id = ? AND user_id = ?", gameId, userId).First(&game).Error
	if err != nil {
		return "", err
	}
	return game.Role, nil
}
