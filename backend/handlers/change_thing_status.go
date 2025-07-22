package handlers

import (
	. "backend/database"
	. "backend/logger"
	"fmt"
	"net/http"
)

func ChangeThingStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			Logger.Error("Error. Can't parse form: " + err.Error())
			fmt.Fprintf(w, "Error. Can't parse form: %s", err.Error())
			return
		}
		// TODO: is it normal that thingId is string, not int?
		thingId, thingType :=
			r.FormValue("thing_id"),
			r.FormValue("thing_type")
		if thingType != "lost" && thingType != "found" {
			Logger.Error("Error. Thing type should be \"lost\" or \"found\"")
			fmt.Fprintf(w, "Error. Thing type should be \"lost\" or \"found\"")
			return
		}
		sqlQuery := fmt.Sprintf("UPDATE %s_thing SET status=1 WHERE id=%s;",
			thingType,
			thingId,
		)
		if _, err := DB.Exec(sqlQuery); err != nil {
			Logger.Error("Error updating thing status: " + err.Error())
			fmt.Fprintf(w, "Error updating thing status: %s", err.Error())
			return
		} else {
			Logger.Info("Success. If a thing with this id and type exists, its status has been updated")
			fmt.Fprintf(w, "Success. If a thing with this id and type exists, its status has been updated")
			return
		}
	} else {
		fmt.Fprintf(w, "A POST request is required")
		Logger.Warn("A POST request is required")
		return
	}
}
