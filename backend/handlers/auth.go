package handlers

import (
	. "backend/logger"
	. "backend/utils"
	"net/http"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w)
	if r.Method != http.MethodGet {
		// TODO: maybe rewrite to GET requests
		msg := "A GET request is required"
		Logger.Warn(msg)
		http.Error(w, msg, http.StatusMethodNotAllowed)
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
