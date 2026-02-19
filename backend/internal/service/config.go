package service

import (
	"time"
)

type Config struct {
	User UserServiceConfig
}

type UserServiceConfig struct {
	BcryptCost             int
	AvatarMaxSize          int64
	AvatarUploadPath       string
	AvatarAllowedMIMETypes []string
}

type AuthServiceConfig struct {
	JWTSecret          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	TokenIssuer        string
	CookieSecure       bool
}
