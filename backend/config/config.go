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
			PortFrontend: GetEnv("FRONTEND_PORT"),
		},
		Bcrypt: types.BcryptConfig {
			Cost: 15, // minimal is 4, maximum is 31, default is 10
		},
		DB: types.DBConfig{
			PathTo: GetEnv("PATH_TO_DB"),
		},
		Env: types.EnvConfig{
			PathTo: GetEnv("PATH_TO_ENV"),
		},
		Logs: types.LogsConfig{
			PathToBackend: filepath.Join(
				GetEnv("PATH_TO_LOGS"),
				GetEnv("BACKEND_LOG"),
			),
		},
		RSA: types.RSAConfig{
			PathToPrivateKey: filepath.Join(
				GetEnv("PATH_TO_ENV"),
				GetEnv("RSA_PRIVATE_KEY"),
			),
			PathToPublicKey: filepath.Join(
				GetEnv("PATH_TO_ENV"),
				GetEnv("RSA_PUBLIC_KEY"),
			),
		},
		Role: types.RolesConfig{
			User:      "user",
			Moderator: "moderator",
		},
		Storage: types.StorageConfig{
			PathTo: GetEnv("PATH_TO_STORAGE"),
		},
	}
}

// Read environment variable by the key or panic if it doesn't exist
func GetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	panic(fmt.Sprintf("The required environment variable \"%s\" is not set", key))
}

var Cfg = newConfig()
