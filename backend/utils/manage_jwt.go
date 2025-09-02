package utils

import (
	"backend/types"
	"crypto/rsa"
	"errors"
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

func parseJWT(token *string, publicKey *rsa.PublicKey) (*jwt.Token, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		method, ok := (*token).Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, errors.New("Unexpected JWT signing method: " + token.Header["alg"].(string))
		}
		if method.Alg() != "RS512" {
			return nil, errors.New("Unsupported algorithm: " + method.Alg())
		}
		return publicKey, nil
	}
	tokenParsed, err := jwt.Parse(*token, keyFunc)
	if err != nil {
		return nil, err
	}
	if !tokenParsed.Valid {
		return nil, errors.New("Invalid token")
	}
	return tokenParsed, nil
}

func VerifyJWTAccess(accessToken *string, publicKey *rsa.PublicKey) error {
	parsedToken, err := parseJWT(accessToken, publicKey)
	switch {
	case parsedToken.Valid:
		return nil
	case errors.Is(err, jwt.ErrTokenMalformed):
		return errors.New("that's not a JWT access token: " + err.Error())
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return errors.New("invalid signature of JWT access token: " + err.Error())
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return errors.New("JWT access token has expired or isn't active yet: " + err.Error())
	default:
		return errors.New("couldn't handle JWT access token: " + err.Error())
	}
}
