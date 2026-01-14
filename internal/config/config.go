package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server *ServerConfig
	DB *DatabaseConfig
}

func Load() (*Config, error) {

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../../") // สำหรับรันจาก internal/
	viper.AutomaticEnv()

	// อ่าน config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	// สร้าง config object
	cfg := &Config{
		Server: LoadServerConfig(),
		DB: LoadDatabaseConfig(),
	}

	return cfg, nil
}

// IsDevelopment ตรวจว่าเป็น dev environment
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction ตรวจว่าเป็น production
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}
