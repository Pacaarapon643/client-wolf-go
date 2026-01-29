package database

import (
	"log"
	"time"
	"werewolf-backend/internal/infrastructure/config"
	"werewolf-backend/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgres(cfg *config.DatabaseConfig) (*Database, error) {
	var db *gorm.DB
	var err error

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt: false, // Disable prepared statements to fix "prepared statement name is already in use" error
	}

	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  cfg.URL,
		PreferSimpleProtocol: true, // Disable implicit prepared statement usage
	}), gormCfg)

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return New(db), nil

}

func (d *Database) AutoMigrate() error {
	err := d.DB.AutoMigrate(
		&models.User{},
	)

	if err != nil {
		return err
	}

	log.Println("✅ Database migrated successfully")
	return nil
}
