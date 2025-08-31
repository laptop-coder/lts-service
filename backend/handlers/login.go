// SEE https://pkg.go.dev/golang.org/x/crypto/bcrypt
package handlers

import (
	. "backend/database"
	. "backend/logger"
	"backend/types"
	. "backend/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func ModeratorLogin(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w)
	const bcryptCost = 15 // minimal is 4, maximum is 31, default is 10

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

	sqlQuery := fmt.Sprintf(
		"SELECT * FROM moderator WHERE username='%s';",
		username,
	)
	row := DB.QueryRow(sqlQuery)
	var moderatorAccountData types.ModeratorAccount
	err := row.Scan(
		&moderatorAccountData.ModeratorId,
		&moderatorAccountData.Username,
		&moderatorAccountData.Email,
		&moderatorAccountData.PasswordHash,
		&moderatorAccountData.CredentialsVersion,
	)
	switch {
	case err == sql.ErrNoRows:
		msg := "Moderator account with this username was not found"
		Logger.Warn(msg)
		http.Error(w, msg, http.StatusUnauthorized)
		return
	case err != nil:
		msg := "Error logging in: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	default:
		err := bcrypt.CompareHashAndPassword([]byte(moderatorAccountData.PasswordHash), []byte(password))
		if err != nil {
			msg := "Passwords don't match"
			Logger.Warn(msg)
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}

		// JWT

		privateKey, err := GetPrivateKey()
		if err != nil {
			msg := "Error getting private key: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		pair, err := CreateJWTPair(moderatorAccountData.Username, privateKey)
		if err != nil {
			msg := "Error creating new JWT pair: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		http.SetCookie(
			w,
			&http.Cookie{
				Name:        "jwt_access",
				Value:       *pair.AccessToken,
				Secure:      true,
				HttpOnly:    true,
				Partitioned: true,
				SameSite:    http.SameSiteNoneMode,
				Path:        "/",
				Expires:     time.Now().Add(5 * time.Minute),
			},
		)
		http.SetCookie(
			w,
			&http.Cookie{
				Name:        "jwt_refresh",
				Value:       *pair.RefreshToken,
				Secure:      true,
				HttpOnly:    true,
				Partitioned: true,
				SameSite:    http.SameSiteNoneMode,
				Path:        "/",
				Expires:     time.Now().Add(30 * 24 * time.Hour),
			},
		)

		jsonData, err := json.Marshal(pair)
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
