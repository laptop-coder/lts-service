package config

import (
	"backend/internal/service"
)

type ServiceConfigs struct {
	User   service.UserServiceConfig
	Post   service.PostServiceConfig
	Auth   service.AuthServiceConfig
	Invite service.InviteServiceConfig
}

func NewServiceConfigs(sharedConfig SharedConfig) ServiceConfigs {
	return ServiceConfigs{
		User:   newUserServiceConfig(sharedConfig),
		Post:   newPostServiceConfig(sharedConfig),
		Auth:   newAuthServiceConfig(sharedConfig),
		Invite: newInviteServiceConfig(sharedConfig),
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

func newPostServiceConfig(sharedConfig SharedConfig) service.PostServiceConfig {
	return service.PostServiceConfig{
		PhotoMaxSize:          sharedConfig.Storage.PostPhoto.MaxSize,
		PhotoUploadPath:       sharedConfig.Storage.PostPhoto.UploadPath,
		PhotoAllowedMIMETypes: sharedConfig.Storage.PostPhoto.AllowedMIMETypes,
	}
}

func newAuthServiceConfig(sharedConfig SharedConfig) service.AuthServiceConfig {
	return service.AuthServiceConfig{
		JWTSecret:          sharedConfig.Security.AuthJWTSecret,
		AccessTokenExpiry:  sharedConfig.Security.AccessTokenExpiry,
		RefreshTokenExpiry: sharedConfig.Security.RefreshTokenExpiry,
		TokenIssuer:        sharedConfig.Security.AuthTokenIssuer,
		CookieSecure:       sharedConfig.Security.CookieSecure,
	}
}

func newInviteServiceConfig(sharedConfig SharedConfig) service.InviteServiceConfig {
	return service.InviteServiceConfig{
		JWTSecret:   sharedConfig.Security.InviteJWTSecret,
		TokenExpiry: sharedConfig.Security.InviteTokenExpiry,
		TokenIssuer: sharedConfig.Security.InviteTokenIssuer,
	}
}
