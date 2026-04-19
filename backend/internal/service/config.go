package service

import (
	"backend/pkg/env"
	"time"
)

// TODO: what is it? :) For what?
type Config struct {
	User UserServiceConfig
}

type UserServiceConfig struct {
	BcryptCost             int
	AvatarMaxSize          int64
	AvatarUploadPath       string
	AvatarAllowedMIMETypes []string
}

type PostServiceConfig struct {
	PhotoMaxSize          int64
	PhotoUploadPath       string
	PhotoAllowedMIMETypes []string
}

type AuthServiceConfig struct {
	JWTSecret          []byte
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	TokenIssuer        string
	CookieSecure       bool
}

type InviteServiceConfig struct {
	JWTSecret   []byte
	TokenExpiry time.Duration
	TokenIssuer string
	FrontendURL string
}

type EmailServiceConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	From        string
	FrontendURL string
	AppMode     env.AppMode
}
