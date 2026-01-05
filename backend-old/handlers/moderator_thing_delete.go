package handlers

import (
	. "backend/config"
	. "backend/database"
	. "backend/logger"
	. "backend/utils"
	"fmt"
	"net/http"
)

func ModeratorDeleteThing(w http.ResponseWriter, r *http.Request) {
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

	if _, err := DB.Exec("DELETE FROM thing WHERE id=?;", thingId); err != nil {
		msg := "Error deleting thing: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	pathToPhoto := fmt.Sprintf(
		"%s/%s.jpeg",
		Cfg.Storage.PathTo,
		thingId,
	)
	if err := DeleteThingPhotoFromStorageIfExists(pathToPhoto); err != nil {
		msg := "Error deleting thing photo from storage: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	msg := "Success. If a thing with the passed \"thingId\" existed, it has been deleted"
	Logger.Info(msg)
	w.Write([]byte(msg))
}
