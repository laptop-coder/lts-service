package config

import (
	"backend/types"
	"fmt"
	"os"
	"path/filepath"
)

func newConfig() *types.Config {
	return &types.Config{
		App: types.AppConfig{
			PortBackend:  "37190",
			PortFrontend: getEnv("FRONTEND_PORT"),
			BcryptCost:   15, // minimal is 4, maximum is 31, default is 10
		},
		DB: types.DBConfig{
			PathTo: getEnv("PATH_TO_DB"),
		},
		Env: types.EnvConfig{
			PathTo: getEnv("PATH_TO_ENV"),
		},
		Logs: types.LogsConfig{
			PathToBackend: filepath.Join(
				getEnv("PATH_TO_LOGS"),
				getEnv("BACKEND_LOG"),
			),
		},
		RSA: types.RSAConfig{
			PathToPrivateKey: filepath.Join(
				getEnv("PATH_TO_ENV"),
				getEnv("RSA_PRIVATE_KEY"),
			),
			PathToPublicKey: filepath.Join(
				getEnv("PATH_TO_ENV"),
				getEnv("RSA_PUBLIC_KEY"),
			),
		},
		Storage: types.StorageConfig{
			PathTo: getEnv("PATH_TO_STORAGE"),
		},
	}
}

// Read environment variable by the key or panic if it doesn't exist
func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	panic(fmt.Sprintf("The required environment variable \"%s\" is not set", key))
}

var Cfg = newConfig()
