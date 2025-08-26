package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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

type ValkeyConfig struct {
	Host string
	Port int
}

type Config struct {
	App     AppConfig
	DB      DBConfig
	Env     EnvConfig
	Logs    LogsConfig
	RSA     RSAConfig
	SSL     SSLConfig
	Storage StorageConfig
	Valkey  ValkeyConfig
}

func newConfig() *Config {
	var VALKEY_PORT, err = strconv.Atoi(getEnv("VALKEY_PORT"))
	if err != nil {
		panic(err)
	}
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
		RSA: RSAConfig{
			PathToPrivateKey: filepath.Join(
				getEnv("PATH_TO_ENV"),
				getEnv("RSA_PRIVATE_KEY"),
			),
			PathToPublicKey: filepath.Join(
				getEnv("PATH_TO_ENV"),
				getEnv("RSA_PUBLIC_KEY"),
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
		Valkey: ValkeyConfig{
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
