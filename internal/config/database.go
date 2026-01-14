package config

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	URL             string        `validate:"required"`
	MaxOpenConns    int           `validate:"required"`
	MaxIdleConns    int           `validate:"required"`
	ConnMaxLifetime time.Duration `validate:"required"`
}

func LoadDatabaseConfig() *DatabaseConfig {
	dbConfig := DatabaseConfig{
		URL:             viper.GetString("DATABASE_URL"),
		MaxOpenConns:    viper.GetInt("DATABASE_MAX_OPEN_CONNS"),
		MaxIdleConns:    viper.GetInt("DATABASE_MAX_IDLE_CONNS"),
		ConnMaxLifetime: viper.GetDuration("DATABASE_CONN_MAX_LIFETIME"),
	}

	validate := validator.New()
	if err := validate.Struct(dbConfig); err != nil {
		log.Fatal("Invalid config: ", err)
	}

	return &dbConfig
}
