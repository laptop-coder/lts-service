package utils

import (
	. "backend/config"
	. "backend/database"
	. "backend/logger"
	"crypto/rsa"
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

func CreateJWTAccess(username *string, privateKey *rsa.PrivateKey, role *string) (*string, error) {
	issuedAt := time.Now()
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"sub":  *username,
		"iat":  issuedAt.Unix(),
		"exp":  (issuedAt.Add(time.Hour * 24 * 30).Unix()), // 30 days
		"role": *role,
	}).SignedString(privateKey)
	if err != nil {
		return nil, err
	}
	return &accessToken, nil
}

func parseJWT(token *string, publicKey *rsa.PublicKey) (*jwt.Token, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		method, ok := (*token).Method.(*jwt.SigningMethodRSA)
		if !ok {
			return nil, errors.New("unexpected JWT signing method: " + token.Header["alg"].(string))
		}
		if method.Alg() != "RS512" {
			return nil, errors.New("unsupported algorithm: " + method.Alg())
		}
		return publicKey, nil
	}
	tokenParsed, err := jwt.Parse(*token, keyFunc)
	return tokenParsed, err
}

func VerifyJWTAccess(accessToken *string, publicKey *rsa.PublicKey, requiredRole *string) error {
	parsedToken, err := parseJWT(accessToken, publicKey)
	switch {
	case parsedToken.Valid:
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			role := claims["role"]
			if role != *requiredRole {
				return errors.New("access denied (insufficient rights)")
			}
		} else {
			return errors.New("can't get JWT access claims")
		}
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

func GetUsername(accessToken *string, publicKey *rsa.PublicKey) (*string, error) {
	parsedToken, err := parseJWT(accessToken, publicKey)
	sub, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return nil, errors.New("can't get JWT access \"sub\" claim")
	}
	return &sub, nil
}

func checkAccountExistence(username *string, role *string) (*bool, error) {
	exists := false
	var row *sql.Row
	switch *role {
	case Cfg.Role.User:
		row = DB.QueryRow("SELECT COUNT(*) FROM user WHERE username=?;", *username)
	case Cfg.Role.Moderator:
		row = DB.QueryRow("SELECT COUNT(*) FROM moderator WHERE username=?;", *username)
	}
	var count int
	if err := row.Scan(&count); err != nil {
		return &exists, err
	}
	if count == 1 {
		exists = true
	}
	return &exists, nil
}

func GetJWTAccess(r *http.Request) (*string, error) {
	var accessToken *string
	accessTokenCookie, err := r.Cookie("LTS_jwt_access")
	if err != nil {
		return nil, errors.New(
			"can't get JWT access from the cookie: " +
				err.Error() +
				". If you are not logged in to your account yet, please log in",
		)
	} else {
		accessToken = &accessTokenCookie.Value
	}
	return accessToken, nil
}

func AuthMiddleware(role *string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		publicKey, _, err := GetPublicKey()
		if err != nil {
			msg := "Error getting public key: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		accessToken, err := GetJWTAccess(r)
		if err != nil {
			msg := err.Error()
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}

		if err := VerifyJWTAccess(accessToken, publicKey, role); err != nil {
			msg := "JWT access verification error: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}

		username, err := GetUsername(accessToken, publicKey)
		if err != nil {
			msg := "Can't get username from JWT access: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		exists, err := checkAccountExistence(username, role)
		if err != nil {
			msg := "Error checking account existence: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		if *exists == false {
			msg := "Account with the username from the JWT access token does not exist"
			Logger.Error(msg)
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
