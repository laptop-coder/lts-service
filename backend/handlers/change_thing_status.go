package handlers

import (
	. "backend/database"
	. "backend/logger"
	. "backend/utils"
	"fmt"
	"io"
	"net/http"
)

func ChangeThingStatus(w http.ResponseWriter, r *http.Request) {
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
	thingId, thingType :=
		r.FormValue("thingId"),
		r.FormValue("thingType")
	if thingId == "" || thingType == "" {
		msg := "error: the \"thingId\" and \"thingType\" parameters are required"
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	if thingType != "lost" && thingType != "found" {
		msg := "Error. Thing type should be \"lost\" or \"found\""
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	sqlQuery := fmt.Sprintf("UPDATE %s_thing SET status=1 WHERE %s_thing_id=%s;",
		thingType,
		thingType,
		thingId,
	)
	if _, err := DB.Exec(sqlQuery); err != nil {
		msg := "Error updating thing status: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	} else {
		msg := "Success. If a thing with this id and type exists, its status has been updated"
		Logger.Info(msg)
		io.WriteString(w, msg)
		return
	}
}
