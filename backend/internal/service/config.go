package service

type Config struct {
	User UserServiceConfig
}

type UserServiceConfig struct {
	BcryptCost             int
	AvatarMaxSize          int64
	AvatarUploadPath       string
	AvatarAllowedMIMETypes []string
}
