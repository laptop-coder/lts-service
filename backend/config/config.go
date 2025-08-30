package config

import (
	"backend/types"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

func newConfig() *types.Config {
	var VALKEY_PORT, err = strconv.Atoi(getEnv("VALKEY_PORT"))
	if err != nil {
		panic(err)
	}
	return &types.Config{
		App: types.AppConfig{
			DevMode: getEnv("LTS_SERVICE_DEV_MODE"),
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
		SSL: types.SSLConfig{
			PathToCert: filepath.Join(
				getEnv("PATH_TO_ENV"),
				getEnv("SSL_CERT"),
			),
			PathToKey: filepath.Join(
				getEnv("PATH_TO_ENV"),
				getEnv("SSL_KEY"),
			),
		},
		Storage: types.StorageConfig{
			PathTo: getEnv("PATH_TO_STORAGE"),
		},
		Valkey: types.ValkeyConfig{
			Host: getEnv("VALKEY_HOST"),
			Port: VALKEY_PORT,
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
