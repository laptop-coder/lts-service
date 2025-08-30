// SEE https://pkg.go.dev/golang.org/x/crypto/bcrypt
package handlers

import (
	. "backend/database"
	. "backend/logger"
	. "backend/utils"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
)

func ModeratorRegister(w http.ResponseWriter, r *http.Request) {
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
	username, email, password :=
		r.FormValue("username"),
		r.FormValue("email"),
		r.FormValue("password")
	if username == "" || password == "" {
		msg := "error: the \"username\", \"email\" and \"password\" parameters are required"
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

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		msg := "Error generating password hash: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	sqlQuery := fmt.Sprintf(
		"INSERT INTO moderator (username, email, password) VALUES ('%s', '%s', '%s');",
		username,
		email,
		passwordHash,
	)
	if _, err := DB.Exec(sqlQuery); err != nil {
		msg := "Error registering a new moderator account: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	} else {
		msg := "Success. A new moderator account has been created"
		Logger.Info(msg)
		io.WriteString(w, msg)
		return
	}
}
