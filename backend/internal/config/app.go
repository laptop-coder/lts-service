package config

import (
	"backend/pkg/env"
)

type AppConfig struct {
	Port int
	FrontendURL string
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		Port: 37190,
		FrontendURL: env.GetStringRequired("FRONTEND_URL"),
	}
}
