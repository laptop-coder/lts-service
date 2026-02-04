package config

import (
// "backend/pkg/env"
)

type AppConfig struct {
	Port int
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		Port: 37190,
	}
}
