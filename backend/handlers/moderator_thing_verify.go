package handlers

import (
	. "backend/database"
	. "backend/logger"
	. "backend/utils"
	"net/http"
)

func VerifyThing(w http.ResponseWriter, r *http.Request) {
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
	thingId := r.FormValue("thingId")
	if thingId == "" {
		msg := "Error. POST parameter \"thingId\" is required"
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// Verify notice
	if _, err := DB.Exec(
		"UPDATE thing SET verified=1 WHERE id=?;",
		thingId,
	); err != nil {
		msg := "Thing verification error: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	msg := "Success. Thing has been verified"
	Logger.Info(msg)
	w.Write([]byte(msg))
}
