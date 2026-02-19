package config

import (
	"backend/internal/service"
)

type ServiceConfigs struct {
	User service.UserServiceConfig
	Auth service.AuthServiceConfig
}

func NewServiceConfigs(sharedConfig SharedConfig) ServiceConfigs {
	return ServiceConfigs{
		User: newUserServiceConfig(sharedConfig),
		Auth: newAuthServiceConfig(sharedConfig),
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

func newAuthServiceConfig(sharedConfig SharedConfig) service.AuthServiceConfig {
	return service.AuthServiceConfig{
		JWTSecret:          sharedConfig.Security.JWTSecret,
		AccessTokenExpiry:  sharedConfig.Security.AccessTokenExpiry,
		RefreshTokenExpiry: sharedConfig.Security.RefreshTokenExpiry,
		TokenIssuer:        sharedConfig.Security.TokenIssuer,
		CookieSecure:       sharedConfig.Security.CookieSecure,
	}
}
