package utils

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTPair struct {
	AccessToken  *string
	RefreshToken *string
}

func CreateJWTPair(username string, privateKey *rsa.PrivateKey) (*JWTPair, error) {
	issuedAt := time.Now().Unix()
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"sub": username,
		"iat": issuedAt,
		"exp": (issuedAt + int64(5*60)), // 5 minutes
	}).SignedString(privateKey)
	if err != nil {
		return nil, err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"sub":                 username,
		"iat":                 issuedAt,
		"exp":                 (issuedAt + int64(30*24*60*60)), // 30 days
		"credentials_version": 0,                               // TODO: maybe set from the function parameter
	}).SignedString(privateKey)
	if err != nil {
		return nil, err
	}
	return &JWTPair{
		AccessToken:  &accessToken,
		RefreshToken: &refreshToken,
	}, nil
}
