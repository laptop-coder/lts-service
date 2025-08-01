package config

import (
	"fmt"
	"os"
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

type RSAConfig struct {
	PrivateKeyPassword string
	PathToPrivateKey   string
	PathToPublicKey    string
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
	RSA     RSAConfig
	SSL     SSLConfig
	Storage StorageConfig
}

func New() *Config {
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
			PathToBackend: fmt.Sprintf(
				"%s/%s",
				getEnv("PATH_TO_LOGS"),
				getEnv("BACKEND_LOG"),
			),
		},
		RSA: RSAConfig{
			PathToPrivateKey: fmt.Sprintf(
				"%s/%s",
				getEnv("PATH_TO_ENV"),
				getEnv("RSA_PRIVATE_KEY"),
			),
			PathToPublicKey: fmt.Sprintf(
				"%s/%s",
				getEnv("PATH_TO_ENV"),
				getEnv("RSA_PUBLIC_KEY"),
			),
			PrivateKeyPassword: getEnv("PRIVATE_KEY_ENCRYPTION_PASSWORD"),
		},
		SSL: SSLConfig{
			PathToCert: fmt.Sprintf(
				"%s/%s",
				getEnv("PATH_TO_ENV"),
				getEnv("SSL_CERT"),
			),
			PathToKey: fmt.Sprintf(
				"%s/%s",
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
