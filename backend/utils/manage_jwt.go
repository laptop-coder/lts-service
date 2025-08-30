package utils

import (
	"backend/types"
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func CreateJWTPair(username string, privateKey *rsa.PrivateKey) (*types.JWTPair, error) {
	issuedAt := time.Now()
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"sub": username,
		"iat": issuedAt.Unix(),
		"exp": (issuedAt.Add(5 * time.Minute).Unix()), // 5 minutes
	}).SignedString(privateKey)
	if err != nil {
		return nil, err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"sub":                 username,
		"iat":                 issuedAt.Unix(),
		"exp":                 (issuedAt.Add(30 * 24 * time.Hour).Unix()), // 30 days
		"credentials_version": 0,                                          // TODO: maybe set from the function parameter
	}).SignedString(privateKey)
	if err != nil {
		return nil, err
	}
	return &types.JWTPair{
		AccessToken:  &accessToken,
		RefreshToken: &refreshToken,
	}, nil
}
