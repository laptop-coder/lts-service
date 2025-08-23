package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type AppConfig struct {
	DevMode string
}

type DBConfig struct {
	PathTo string
}

type EnvConfig struct {
	PathTo string
}

type LogsConfig struct {
	PathToBackend string
}

type ED25519Config struct {
	PathToPrivateKey string
	PathToPublicKey  string
}

type SSLConfig struct {
	PathToCert string
	PathToKey  string
}

type StorageConfig struct {
	PathTo string
}

type Config struct {
	App     AppConfig
	DB      DBConfig
	Env     EnvConfig
	Logs    LogsConfig
	ED25519 ED25519Config
	SSL     SSLConfig
	Storage StorageConfig
}

func newConfig() *Config {
	return &Config{
		App: AppConfig{
			DevMode: getEnv("LTS_SERVICE_DEV_MODE"),
		},
		DB: DBConfig{
			PathTo: getEnv("PATH_TO_DB"),
		},
		Env: EnvConfig{
			PathTo: getEnv("PATH_TO_ENV"),
		},
		Logs: LogsConfig{
			PathToBackend: filepath.Join(
				getEnv("PATH_TO_LOGS"),
				getEnv("BACKEND_LOG"),
			),
		},
		ED25519: ED25519Config{
			PathToPrivateKey: filepath.Join(
				getEnv("PATH_TO_ENV"),
				getEnv("ED25519_PRIVATE_KEY"),
			),
			PathToPublicKey: filepath.Join(
				getEnv("PATH_TO_ENV"),
				getEnv("ED25519_PUBLIC_KEY"),
			),
		},
		SSL: SSLConfig{
			PathToCert: filepath.Join(
				getEnv("PATH_TO_ENV"),
				getEnv("SSL_CERT"),
			),
			PathToKey: filepath.Join(
				getEnv("PATH_TO_ENV"),
				getEnv("SSL_KEY"),
			),
		},
		Storage: StorageConfig{
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
