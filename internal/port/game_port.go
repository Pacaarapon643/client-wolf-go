package port

import (
	"context"
	"werewolf-backend/internal/adapter/http/dto"
	"werewolf-backend/internal/models"
)

type GameService interface {
	RoleDistribution(ctx context.Context, roomId string) (*string, error)
	GetRoleGame(ctx context.Context, roomId string, userId string) (any, error)
	GetGame(ctx context.Context, gameId string, role string) (any, error)
	JoinGame(ctx context.Context, gameId string, userId string) error
	LeaveGame(ctx context.Context, gameId string, userId string) error
	Vote(ctx context.Context, gameId string, userId string, vote string, role string) (any, error)
	SummaryVote(ctx context.Context, gameId string, phase string) (any, error)
}

type GameRepo interface {
	CreateGame(ctx context.Context, game []models.Game) error
	GetGame(ctx context.Context, arr *[]dto.Game, game_id *string) error
	JoinGame(ctx context.Context, gameId *string, userId *string) error
	LeaveGame(ctx context.Context, gameId *string, userId *string) error
	Dead(ctx context.Context, gameId *string, index *int) error
}
