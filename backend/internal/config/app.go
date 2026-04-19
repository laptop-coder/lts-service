package config

import (
	"backend/pkg/env"
	"fmt"
)

type AppConfig struct {
	Port        int
	FrontendURL string
	AppMode     env.AppMode
}

func LoadAppConfig(appMode env.AppMode) AppConfig {
	return AppConfig{
		Port:        37190,
		FrontendURL: env.GetStringRequired("FRONTEND_URL"),
		AppMode:     appMode,
	}
}

func ParseAppMode(v string) env.AppMode {
	switch v {
	case string(env.AppModeDev):
		return env.AppModeDev
	case string(env.AppModeProd):
		return env.AppModeProd
	default:
		panic(fmt.Sprintf("unknown app mode: %s (expected dev or prod)", v))
	}
}
