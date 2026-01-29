package config

import (
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type JWTConfig struct {
	AccessTokenExpiry  time.Duration `validate:"required"`
	RefreshTokenExpiry time.Duration `validate:"required"`
	PrivateKeyPath     string        `validate:"required"`
	PublicKeyPath      string        `validate:"required"`
}

func LoadJWTConfig() *JWTConfig {
	jtwConfig := JWTConfig{
		AccessTokenExpiry:  viper.GetDuration("JWT_ACCESS_TOKEN_EXPIRY"),
		RefreshTokenExpiry: viper.GetDuration("JWT_REFRESH_TOKEN_EXPIRY"),
		PrivateKeyPath:     viper.GetString("JWT_PRIVATE_KEY_PATH"),
		PublicKeyPath:      viper.GetString("JWT_PUBLIC_KEY_PATH"),
	}

	validate := validator.New()
	if err := validate.Struct(jtwConfig); err != nil {
		log.Fatal("Invalid config: ", err)
	}

	return &jtwConfig
}
