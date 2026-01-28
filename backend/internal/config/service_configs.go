package config

import (
	"backend/internal/service"
)

type ServiceConfigs struct {
	User service.UserServiceConfig
}

func NewServiceConfigs() ServiceConfigs {
	return ServiceConfigs{
		User: newUserServiceConfig(),
	}
}

func newUserServiceConfig() service.UserServiceConfig {
	return service.UserServiceConfig{
		BcryptCost:             Bcrypt.Cost,
		AvatarMaxSize:          Storage.Avatar.MaxSize,
		AvatarUploadPath:       Storage.Avatar.UploadPath,
		AvatarAllowedMIMETypes: Storage.Avatar.AllowedMIMETypes,
	}
}
