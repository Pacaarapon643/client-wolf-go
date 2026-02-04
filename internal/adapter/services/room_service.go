package services

import (
	"context"
	"werewolf-backend/internal/models"
)

func (s Service) CreateRoom(ctx context.Context, obj *models.Room) error {
	err := s.r.CreateRoom(ctx, obj)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) GetRoom(ctx context.Context, obj *[]models.Room) error {
	err := s.r.GetRoom(ctx, obj)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) CountRoom(ctx context.Context, count *int64) error {
	return s.r.CountRoom(ctx, count)
}
