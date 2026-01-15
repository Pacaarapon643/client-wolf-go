package services

import (
	"context"
	"werewolf-backend/internal/dto"
	"werewolf-backend/internal/models"
	"werewolf-backend/internal/util"
)

func (s *Service) Register(ctx context.Context, req dto.RegisterRequest) error {

	password, err := util.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &models.User{
		UserName: req.UserName,
		Email:    req.Email,
		Password: password,
	}

	return s.r.Register(ctx, user)
}
