package services

import (
	"context"
	"werewolf-backend/internal/adapter/http/dto"
	"werewolf-backend/internal/models"
	"werewolf-backend/internal/pkg/util"

	"github.com/gofiber/fiber/v2"
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

func (s *Service) Login(ctx context.Context, req dto.LoginRequest) (*models.User, error) {
	user, err := s.r.Login(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if !util.CheckPasswordHash(req.Password, user.Password) {
		return nil, &util.LocalizedError{
			Code:    fiber.StatusBadRequest,
			Message: "รหัสผ่านไม่ถูกต้อง",
		}
	}

	return user, nil
}
