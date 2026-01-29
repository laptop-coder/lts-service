package config

import (
	"backend/pkg/env"
	"path/filepath"
)

type SharedConfig struct {
	Security SecurityConfig
	Storage  StorageConfig
}

type SecurityConfig struct {
	BcryptCost int
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
			BcryptCost: 15, // minimal is 4, maximum is 31, default is 10
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
				AllowedMIMETypes: []string{"image/jpeg", "image/png", "image/webp"},
			},
		},
	}
}
