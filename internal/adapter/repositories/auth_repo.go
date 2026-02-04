package repositories

import (
	"context"
	"fmt"
	"log"
	"werewolf-backend/internal/models"
	"werewolf-backend/internal/pkg/util"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *Repository) Register(ctx context.Context, user *models.User) error {

	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		// แปลง err เป็น pgconn.PgError
		if pgErr, ok := err.(*pgconn.PgError); ok {
			log.Println("pgErr", pgErr.Code)
			if pgErr.Code == "23505" { // 23505 คือรหัสข้อมูลซ้ำของ Postgres
				return fmt.Errorf("มี email นี้อยู่ในระบบแล้ว")
			}
		}
		return err
	}

	return nil
}

func (r *Repository) Login(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}

	err := r.db.WithContext(ctx).
		Raw(`
            UPDATE users 
            SET is_online = true, updated_at = NOW() 
            WHERE email = ? AND deleted_at IS NULL 
            RETURNING *
        `, email).
		Scan(user).Error

	if err != nil {
		return nil, &util.LocalizedError{
			Code:    fiber.StatusBadRequest,
			Err:     err.Error(),
			Message: "ไม่พบข้อมูลนี้ในระบบกรุณาสมัครสมาชิก",
		}
	}

	return user, nil
}
