// Package config provides constants, settings and environment variables to use
// in code.
package config

import (
	"backend/internal/types"
	"fmt"
	"os"
	"path/filepath"
)

var (
	App = types.AppConfig{
		PortBackend:  "37190",
		PortFrontend: GetEnv("FRONTEND_PORT"),
	}
	Bcrypt = types.BcryptConfig{
		Cost: 15, // minimal is 4, maximum is 31, default is 10
	}
	Postgres = types.PostgresConfig{
		DBName:   GetEnv("POSTGRES_DB"),
		Host:     GetEnv("POSTGRES_HOST"),
		Password: GetEnv("POSTGRES_PASSWORD"),
		Port:     GetEnv("POSTGRES_PORT"),
		SSLMode:  GetEnv("POSTGRES_SSL_MODE"),
		TimeZone: GetEnv("POSTGRES_TIME_ZONE"),
		User:     GetEnv("POSTGRES_USER"),
	}
	Env = types.EnvConfig{
		PathTo: GetEnv("PATH_TO_ENV"),
	}
	RSA = types.RSAConfig{
		PathToPrivateKey: filepath.Join(
			GetEnv("PATH_TO_ENV"),
			GetEnv("RSA_PRIVATE_KEY"),
		),
		PathToPublicKey: filepath.Join(
			GetEnv("PATH_TO_ENV"),
			GetEnv("RSA_PUBLIC_KEY"),
		),
	}
	Storage = types.StorageConfig{
		PathTo: GetEnv("PATH_TO_STORAGE"),
	}
)

// GetEnv allows to read environment variable by the key. If it doesn't exist,
// the function will cause panic.
func GetEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	panic(fmt.Sprintf("The required environment variable \"%s\" is not set", key))
}
