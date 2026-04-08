package config

import (
	"backend/internal/service"
	"backend/pkg/env"
)

type ServiceConfigs struct {
	User   service.UserServiceConfig
	Post   service.PostServiceConfig
	Auth   service.AuthServiceConfig
	Invite service.InviteServiceConfig
	Email  service.EmailServiceConfig
}

func NewServiceConfigs(sharedConfig SharedConfig, appConfig AppConfig) ServiceConfigs {
	return ServiceConfigs{
		User:   newUserServiceConfig(sharedConfig),
		Post:   newPostServiceConfig(sharedConfig),
		Auth:   newAuthServiceConfig(sharedConfig),
		Invite: newInviteServiceConfig(sharedConfig, appConfig),
		Email:  newEmailServiceConfig(),
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

func newInviteServiceConfig(sharedConfig SharedConfig, appConfig AppConfig) service.InviteServiceConfig {
	return service.InviteServiceConfig{
		JWTSecret:   sharedConfig.Security.InviteJWTSecret,
		TokenExpiry: sharedConfig.Security.InviteTokenExpiry,
		TokenIssuer: sharedConfig.Security.InviteTokenIssuer,
		FrontendURL: appConfig.FrontendURL,
	}
}

func newEmailServiceConfig() service.EmailServiceConfig {
	return service.EmailServiceConfig{
		Host:     env.GetStringRequired("EMAIL_HOST"),
		Port:     env.GetIntRequired("EMAIL_PORT"),
		Username: env.GetStringRequired("EMAIL_USERNAME"),
		Password: env.GetStringRequired("EMAIL_PASSWORD"),
		From:     env.GetStringRequired("EMAIL_USERNAME"),
	}
}
