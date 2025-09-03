package utils

import (
	. "backend/database"
	"backend/types"
	"crypto/rsa"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func CreateJWTPair(username *string, privateKey *rsa.PrivateKey) (*types.JWTPair, error) {
	// TODO: refactor, the code is duplicated
	issuedAt := time.Now()
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"sub": *username,
		"iat": issuedAt.Unix(),
		"exp": (issuedAt.Add(5 * time.Minute).Unix()), // 5 minutes
	}).SignedString(privateKey)
	if err != nil {
		return nil, err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"sub":                 *username,
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

func RefreshJWTAccess(refreshToken *string) (*string, error) {
	publicKey, _, err := GetPublicKey()
	if err != nil {
		return nil, errors.New("Error getting public key: " + err.Error())
	}

	sub, err := VerifyJWTRefresh(refreshToken, publicKey)
	if err != nil {
		return nil, errors.New("Error getting \"sub\" claim from JWT refresh: " + err.Error())
	}

	privateKey, _, err := GetPrivateKey()
	if err != nil {
		return nil, errors.New("Error getting private key: " + err.Error())
	}
	// TODO: refactor, the code is duplicated
	issuedAt := time.Now()
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"sub": sub,
		"iat": issuedAt.Unix(),
		"exp": (issuedAt.Add(5 * time.Minute).Unix()), // 5 minutes
	}).SignedString(privateKey)
	if err != nil {
		return nil, err
	}
	return &accessToken, nil
}

func VerifyJWTRefresh(refreshToken *string, publicKey *rsa.PublicKey) (*string, error) {
	parsedToken, err := parseJWT(refreshToken, publicKey)
	switch {
	case parsedToken.Valid:
		// Verify the version of the credentials. The version updates if the
		// user changes credentials, e.g., password.

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errors.New("can't make type assertion (*jwt.Token to jwt.MapClaims)")
		}
		credentialsVersionFromToken := claims["credentials_version"].(float64)
		if !ok {
			return nil, errors.New("can't get JWT refresh \"credentials_version\" claim or can't make type assertion (any to float64)")
		}
		sub, err := parsedToken.Claims.GetSubject()
		if err != nil {
			return nil, errors.New("can't get JWT refresh \"sub\" claim")
		}

		sqlQuery := fmt.Sprintf(
			"SELECT credentials_version FROM moderator WHERE username='%s';",
			sub,
		)
		row := DB.QueryRow(sqlQuery)
		var credentialsVersionFromDatabase float64

		err = row.Scan(&credentialsVersionFromDatabase)
		switch {
		case err == sql.ErrNoRows:
			return nil, errors.New("error getting credentials version from the database: moderator account with this username was not found: " + err.Error())
		case err != nil:
			return nil, errors.New("error getting credentials version from the database: " + err.Error())
		default:
			if credentialsVersionFromToken != credentialsVersionFromDatabase {
				return nil, errors.New("the token is invalid: the credentials have been changed since the token was issued")
			}
		}

		return &sub, nil
	case errors.Is(err, jwt.ErrTokenMalformed):
		return nil, errors.New("that's not a JWT refresh token: " + err.Error())
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return nil, errors.New("invalid signature of JWT refresh token: " + err.Error())
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return nil, errors.New("JWT refresh token has expired or isn't active yet: " + err.Error())
	default:
		return nil, errors.New("couldn't handle JWT refresh token: " + err.Error())
	}
}
