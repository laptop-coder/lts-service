package config

import (
	"backend/pkg/env"
	"fmt"
	"path/filepath"
	"time"
)

type SharedConfig struct {
	Security SecurityConfig
	Storage  StorageConfig
}

type SecurityConfig struct {
	BcryptCost         int
	AuthJWTSecret      []byte
	InviteJWTSecret    []byte
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	InviteTokenExpiry  time.Duration
	AuthTokenIssuer    string
	InviteTokenIssuer  string
	CookieSecure       bool
}

type ImageStorageConfig struct {
	UploadPath       string
	MaxSize          int64 // in bytes
	AllowedMIMETypes []string
}

type StorageConfig struct {
	Avatar    ImageStorageConfig
	PostPhoto ImageStorageConfig
}

func LoadSharedConfig() SharedConfig {
	return SharedConfig{
		Security: SecurityConfig{
			BcryptCost:         12, // minimal is 4, maximum is 31, default is 10
			AuthJWTSecret:      []byte(env.GetStringRequired("JWT_SECRET_AUTH")),
			InviteJWTSecret:    []byte(env.GetStringRequired("JWT_SECRET_INVITE")),
			AccessTokenExpiry:  time.Duration(time.Minute * 15),
			RefreshTokenExpiry: time.Duration(time.Hour * 24 * 30),
			InviteTokenExpiry:  time.Duration(time.Hour * 24 * 7),
			AuthTokenIssuer:    fmt.Sprintf("%s-auth", env.GetStringRequired("JWT_ISSUER_PREFIX")),
			InviteTokenIssuer:  fmt.Sprintf("%s-invite", env.GetStringRequired("JWT_ISSUER_PREFIX")),
			CookieSecure:       env.GetBoolRequired("COOKIE_SECURE"),
		},
		Storage: StorageConfig{
			Avatar: ImageStorageConfig{
				UploadPath:       filepath.Join(env.GetStringRequired("PATH_TO_STORAGE"), "avatars"),
				MaxSize:          5 * 1024 * 1024, // 5 MB
				AllowedMIMETypes: []string{"image/jpeg", "image/png", "image/webp", "image/gif"},
			},
			PostPhoto: ImageStorageConfig{
				UploadPath:       filepath.Join(env.GetStringRequired("PATH_TO_STORAGE"), "post_photos"),
				MaxSize:          10 * 1024 * 1024, // 10 MB
				AllowedMIMETypes: []string{"image/jpeg", "image/png", "image/webp", "image/gif"},
			},
		},
	}
}
