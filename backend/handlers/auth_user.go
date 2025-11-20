// SEE https://pkg.go.dev/golang.org/x/crypto/bcrypt
package handlers

import (
	. "backend/config"
	. "backend/database"
	. "backend/logger"
	"backend/types"
	. "backend/utils"
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"time"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w)

	if r.Method != http.MethodPost {
		msg := "A POST request is required"
		Logger.Warn(msg)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		msg := "Error. Can't parse form: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	username, password :=
		r.FormValue("username"),
		r.FormValue("password")
	if username == "" || password == "" {
		msg := "error: the \"username\" and \"password\" parameters are required"
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if err := ValidatePassword(password); err != nil {
		msg := "Error validating password: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	row := DB.QueryRow(
		"SELECT * FROM user WHERE username=?;",
		username,
	)
	var userAccountData types.UserAccountAuthorizationData
	err := row.Scan(
		&userAccountData.Username,
		&userAccountData.PasswordHash,
	)
	switch {
	case err == sql.ErrNoRows:
		msg := "User account with this username was not found"
		Logger.Warn(msg)
		http.Error(w, msg, http.StatusUnauthorized)
		return
	case err != nil:
		msg := "Error logging in: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	default:
		err := bcrypt.CompareHashAndPassword([]byte(userAccountData.PasswordHash), []byte(password))
		if err != nil {
			msg := "Passwords don't match"
			Logger.Warn(msg)
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}

		// JWT

		privateKey, _, err := GetPrivateKey()
		if err != nil {
			msg := "Error getting private key: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		accessToken, err := CreateJWTAccess(&userAccountData.Username, privateKey, &Cfg.Role.User)
		if err != nil {
			msg := "Error creating new JWT access: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		// TODO: "Secure: true"
		http.SetCookie(
			w,
			&http.Cookie{
				Name:     "user_jwt_access",
				Value:    *accessToken,
				HttpOnly: true,
				Path:     "/",                                 // TODO: is it OK?
				Expires:  time.Now().Add(time.Hour * 24 * 30), // 30 days
				// Partitioned: true,
				// SameSite:    http.SameSiteNoneMode,
				// Domain:      "localhost",
			},
		)
		http.SetCookie(
			w,
			&http.Cookie{
				Name:     "user_authorized",
				Value:    "true",
				HttpOnly: false,
				Path:     "/",                                 // TODO: is it OK?
				Expires:  time.Now().Add(time.Hour * 24 * 30), // 30 days
				// Partitioned: true,
				// SameSite:    http.SameSiteNoneMode,
				// Domain:      "localhost",
			},
		)

		jsonData, err := json.Marshal(accessToken)
		if err != nil {
			msg := "JSON serialization error: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
		Logger.Info("Success. Logged in")
		return
	}
}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w)

	if r.Method != http.MethodPost {
		msg := "A POST request is required"
		Logger.Warn(msg)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		msg := "Error. Can't parse form: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	username, password :=
		r.FormValue("username"),
		r.FormValue("password")
	if username == "" || password == "" {
		msg := "Error: the \"username\" and \"password\" parameters are required"
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if err := ValidatePassword(password); err != nil {
		msg := "Error validating password: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), Cfg.Bcrypt.Cost)
	if err != nil {
		msg := "Error generating password hash: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	if _, err := DB.Exec(
		"INSERT INTO user (username, password_hash) VALUES (?, ?);",
		username,
		passwordHash,
	); err != nil {
		msg := "Error registering new user account: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	} else {
		msg := "Success. A new user account has been created"
		Logger.Info(msg)
		io.WriteString(w, msg)
		return
	}
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w)
	if r.Method != http.MethodPost {
		msg := "A POST request is required"
		Logger.Warn(msg)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		msg := "Error. Can't parse form: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// Delete cookies
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "user_jwt_access",
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			MaxAge:   -1,
		},
	)
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "user_authorized",
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			MaxAge:   -1,
		},
	)

	msg := "Success. Logged out from the user account"
	Logger.Info(msg)
	w.Write([]byte(msg))
}
