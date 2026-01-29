package config

import (
	"backend/internal/service"
)

type ServiceConfigs struct {
	User service.UserServiceConfig
}

func NewServiceConfigs(sharedConfig SharedConfig) ServiceConfigs {
	return ServiceConfigs{
		User: newUserServiceConfig(sharedConfig),
	}
}

func newUserServiceConfig(sharedConfig SharedConfig) service.UserServiceConfig {
	return service.UserServiceConfig{
		BcryptCost:             sharedConfig.Security.BcryptCost,
		AvatarMaxSize:          sharedConfig.Storage.Avatar.MaxSize,
		AvatarUploadPath:       sharedConfig.Storage.Avatar.UploadPath,
		AvatarAllowedMIMETypes: sharedConfig.Storage.Avatar.AllowedMIMETypes,
	}
}
