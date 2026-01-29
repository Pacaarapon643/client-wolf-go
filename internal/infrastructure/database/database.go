package database

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func New(db *gorm.DB) *Database {
	return &Database{DB: db}
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func (d *Database) Health(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}

func (d *Database) GetStats() map[string]interface{} {
	sqlDB, _ := d.DB.DB()
	stats := sqlDB.Stats()

	return map[string]interface{}{
		"open_connections": stats.OpenConnections,
		"in_use":           stats.InUse,
		"idle":             stats.Idle,
	}
}
