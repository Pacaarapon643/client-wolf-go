package repositories

import (
	"context"
	"fmt"
	"log"
	"werewolf-backend/internal/models"

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
