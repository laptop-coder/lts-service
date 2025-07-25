package handlers

import (
	. "backend/database"
	. "backend/logger"
	"fmt"
	"io"
	"net/http"
)

func ChangeThingStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			msg := "Error. Can't parse form: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		// TODO: is it normal that thingId is string, not int?
		thingId, thingType :=
			r.FormValue("thing_id"),
			r.FormValue("thing_type")
		if thingId == "" || thingType == "" {
			msg := "error: the \"thing_id\" and \"thing_type\" parameters are required"
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		if thingType != "lost" && thingType != "found" {
			msg := "Error. Thing type should be \"lost\" or \"found\""
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		sqlQuery := fmt.Sprintf("UPDATE %s_thing SET status=1 WHERE id=%s;",
			thingType,
			thingId,
		)
		if _, err := DB.Exec(sqlQuery); err != nil {
			msg := "Error updating thing status: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		} else {
			msg := "Success. If a thing with this id and type exists, its status has been updated"
			Logger.Info(msg)
			io.WriteString(w, msg)
			return
		}
	} else {
		msg := "A POST request is required"
		Logger.Warn(msg)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}
}
