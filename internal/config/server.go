package config

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Address         string        `validate:"required"`
	Port            string        `validate:"required"`
	Environment     string        `validate:"required"`
	ReadTimeout     time.Duration `validate:"required"`
	WriteTimeout    time.Duration `validate:"required"`
	ShutdownTimeout time.Duration `validate:"required"`
}

func LoadServerConfig() *ServerConfig {
	serverConfig := ServerConfig{
		Address:         viper.GetString("ADDRESS"),
		Port:            viper.GetString("PORT"),
		Environment:     viper.GetString("ENVIRONMENT"),
		ReadTimeout:     viper.GetDuration("SERVER_READ_TIMEOUT"),
		WriteTimeout:    viper.GetDuration("SERVER_WRITE_TIMEOUT"),
		ShutdownTimeout: viper.GetDuration("SERVER_SHUTDOWN_TIMEOUT"),
	}

	validate := validator.New()
	if err := validate.Struct(serverConfig); err != nil {
		log.Fatal("Invalid config: ", err)
	}

	return &serverConfig
}
