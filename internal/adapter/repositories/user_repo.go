package repositories

import (
	"context"
	"werewolf-backend/internal/models"
)

func (r Repository) CountUser(ctx context.Context, count *int64) error {

	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("is_online = ?", true).
		Count(count).
		Error
	if err != nil {
		return err
	}

	return nil
}
