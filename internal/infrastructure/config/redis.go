package config

import (
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type RedisConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	Password string `validate:"required"`
	DB       int
	UseTLS   bool `validate:"required"`
}

func LoadRedisConfig() *RedisConfig {
	redisConfig := &RedisConfig{
		Host:     viper.GetString("REDIS_HOST"),
		Port:     viper.GetString("REDIS_PORT"),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB:       viper.GetInt("REDIS_DB"),
		UseTLS:   viper.GetBool("REDIS_USE_TLS"),
	}

	validate := validator.New()
	if err := validate.Struct(redisConfig); err != nil {
		log.Fatal("Invalid config: ", err)
	}

	return redisConfig

}
