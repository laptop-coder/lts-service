package handlers

import (
	. "backend/database"
	. "backend/logger"
	. "backend/utils"
	"net/http"
)

func ThingChangeVerification(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w)
	// Think about it: GET or POST is better for this handler?
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
	action := r.FormValue("action")
	if action == "" {
		msg := "Error. POST parameter \"action\" is required"
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// Change notice verification
	switch action {
	case "approve":
		if _, err := DB.Exec(
			"UPDATE thing SET verified=1 WHERE id=?;",
			thingId,
		); err != nil {
			msg := "Thing approving (changing verification status) error: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	case "reject":
		if _, err := DB.Exec(
			"UPDATE thing SET verified=-1 WHERE id=?;",
			thingId,
		); err != nil {
			msg := "Thing rejecting (changing verification status) error: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	default:
		msg := "Error. POST parameter \"action\" must be \"approve\" or \"reject\""
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	msg := "Success. Thing has been approved/rejected by the moderator"
	Logger.Info(msg)
	w.Write([]byte(msg))
}
