package handlers

import (
	. "backend/logger"
	. "backend/utils"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
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
			Name:     "jwt_access",
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			MaxAge:   -1,
		},
	)
	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "authorized",
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			MaxAge:   -1,
		},
	)

	msg := "Success. Logged out"
	Logger.Info(msg)
	w.Write([]byte(msg))
}
