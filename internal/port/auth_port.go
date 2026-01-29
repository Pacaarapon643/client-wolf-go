package port

import (
	"context"

	"werewolf-backend/internal/adapter/http/dto"
	"werewolf-backend/internal/models"
)

type AuthRepo interface {
	Register(ctx context.Context, data *models.User) error
	Login(ctx context.Context, email string) (*models.User, error)
}

type AuthService interface {
	Register(ctx context.Context, data dto.RegisterRequest) error
	Login(ctx context.Context, req dto.LoginRequest) (*models.User, error)
}
