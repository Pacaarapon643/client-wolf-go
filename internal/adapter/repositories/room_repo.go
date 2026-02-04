package repositories

import (
	"context"
	"werewolf-backend/internal/models"
)

func (r Repository) CreateRoom(ctx context.Context, obj *models.Room) error {

	err := r.db.WithContext(ctx).Create(&obj).Error
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) GetRoom(ctx context.Context, obj *[]models.Room) error {

	err := r.db.WithContext(ctx).
		Where("is_end = ?", false).
		Find(&obj).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (r Repository) CountRoom(ctx context.Context, count *int64) error {

	err := r.db.WithContext(ctx).
		Model(&models.Room{}).
		Where("is_end = ?", false).
		Count(count).
		Error
	if err != nil {
		return err
	}

	return nil
}
