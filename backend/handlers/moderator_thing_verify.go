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
	thingType := r.FormValue("thingType")
	thingId := r.FormValue("thingId")
	if thingType == "" || thingId == "" {
		msg := "Error. POST parameters \"thingType\" and \"thingId\" are required"
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	switch thingType {
	case "lost":

		// Verify advertisement
		if _, err := DB.Exec(
			"UPDATE lost_thing SET verified=1 WHERE id=?;",
			thingId,
		); err != nil {
			msg := "Lost thing verification error: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		msg := "Success. Lost thing has been verified"
		Logger.Info(msg)
		w.Write([]byte(msg))

	case "found":

		// Verify advertisement
		if _, err := DB.Exec(
			"UPDATE found_thing SET verified=1 WHERE id=?;",
			thingId,
		); err != nil {
			msg := "Found thing verification error: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		msg := "Success. Found thing has been verified"
		Logger.Info(msg)
		w.Write([]byte(msg))

	default:
		msg := "Error. POST parameter \"thingType\" can be \"lost\" or \"found\""
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
}
