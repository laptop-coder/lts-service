package handlers

import (
	. "backend/database"
	. "backend/logger"
	"backend/types"
	. "backend/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetThingsList(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w)
	if r.Method != http.MethodGet {
		msg := "A GET request is required"
		Logger.Warn(msg)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	thingsType := r.URL.Query().Get("things_type")
	thingsStatus := r.URL.Query().Get("things_status")
	thingsVerified := r.URL.Query().Get("things_verified")

	if thingsType == "" {
		msg := "Error. Send request with \"?things_type=lost\" or \"?things_type=found\""
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if thingsType != "lost" && thingsType != "found" {
		msg := "Error. Things type should be \"lost\" or \"found\""
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	if thingsStatus != "" && thingsStatus != "0" && thingsStatus != "1" {
		msg := "Error. Things status should be \"0\", \"1\" or empty"
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	if thingsVerified != "" && thingsVerified != "-1" && thingsVerified != "0" && thingsVerified != "1" {
		msg := "Error. Things verification status should be \"-1\", \"0\", \"1\" or empty"
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// Get data from the database

	sqlQueryFilter := ""
	switch {
	case thingsStatus != "" && thingsVerified != "":
		sqlQueryFilter = fmt.Sprintf("WHERE status=%s AND verified=%s", thingsStatus, thingsVerified)
	case thingsStatus != "" && thingsVerified == "":
		sqlQueryFilter = fmt.Sprintf("WHERE status=%s", thingsStatus)
	case thingsStatus == "" && thingsVerified != "":
		sqlQueryFilter = fmt.Sprintf("WHERE verified=%s", thingsVerified)
	}

	sqlQuery := fmt.Sprintf("SELECT * FROM %s_thing %s ORDER BY %s_thing_id DESC;", thingsType, sqlQueryFilter, thingsType)
	if rows, err := DB.Query(sqlQuery); err != nil {
		msg := "Error getting things list from the database: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	} else {
		Logger.Info(fmt.Sprintf("Success. Received list of %s things", thingsType))
		// Serialize data and send it in response
		switch thingsType {
		case "lost":
			var lostThingsList []types.LostThing
			var thing types.LostThing
			for rows.Next() {
				if err := rows.Scan(
					&thing.LostThingId,
					&thing.PublicationDatetime,
					&thing.ThingName,
					&thing.UserEmail,
					&thing.CustomText,
					&thing.Verified,
					&thing.Status,
				); err != nil {
					msg := "Error creating \"thing\" object: " + err.Error()
					Logger.Error(msg)
					http.Error(w, msg, http.StatusInternalServerError)
					return
				}
				lostThingsList = append(lostThingsList, thing)
			}
			jsonData, err := json.Marshal(lostThingsList)
			if err != nil {
				msg := "JSON serialization error: " + err.Error()
				Logger.Error(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
			w.Write(jsonData)
			return
		case "found":
			var foundThingsList []types.FoundThing
			var thing types.FoundThing
			for rows.Next() {
				if err := rows.Scan(
					&thing.FoundThingId,
					&thing.PublicationDatetime,
					&thing.ThingName,
					&thing.ThingLocation,
					&thing.CustomText,
					&thing.Verified,
					&thing.Status,
				); err != nil {
					msg := "Error creating \"thing\" object: " + err.Error()
					Logger.Error(msg)
					http.Error(w, msg, http.StatusInternalServerError)
					return
				}
				foundThingsList = append(foundThingsList, thing)
			}
			jsonData, err := json.Marshal(foundThingsList)
			if err != nil {
				msg := "JSON serialization error: " + err.Error()
				Logger.Error(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
			w.Write(jsonData)
			return
		}

	}
}
