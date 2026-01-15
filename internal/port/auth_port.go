package port

import (
	"context"
	"werewolf-backend/internal/dto"
	"werewolf-backend/internal/models"
)

type AuthRepo interface {
	Register(ctx context.Context, data *models.User) error
}

type AuthService interface {
	Register(ctx context.Context, data dto.RegisterRequest) error
}
