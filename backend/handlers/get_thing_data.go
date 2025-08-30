package handlers

import (
	. "backend/database"
	. "backend/logger"
	"backend/types"
	. "backend/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

func GetThingData(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w)
	if r.Method != http.MethodGet {
		msg := "A GET request is required"
		Logger.Warn(msg)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	thingId := r.URL.Query().Get("thing_id")
	thingType := r.URL.Query().Get("thing_type")
	// Parameters are empty
	if thingId == "" || thingType == "" {
		msg := "Error. Send request with \"thing_id\" and \"thing_type\" (can be \"lost\" or \"found\") parameters"
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	// Parameters are incorrect
	if !regexp.MustCompile(`^[1-9]\d*$`).MatchString(thingId) {
		// Regular expression: string is a number without leading zeros
		msg := "Error. Thing id should be a number without leading zeros"
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
	// Get data from the database
	sqlQuery := fmt.Sprintf("SELECT * FROM %s_thing WHERE %s_thing_id=%s;", thingType, thingType, thingId)
	row := DB.QueryRow(sqlQuery)
	switch thingType {
	case "lost":
		var thing types.LostThing
		err := row.Scan(
			&thing.LostThingId,
			&thing.PublicationDatetime,
			&thing.ThingName,
			&thing.UserEmail,
			&thing.CustomText,
			&thing.Verified,
			&thing.Status,
		)
		// TODO: refactor, the code is duplicated
		switch {
		case err == sql.ErrNoRows:
			msg := "Thing not found"
			Logger.Info(msg)
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(msg))
			return
			// TODO: maybe rewrite to this:
			// msg := "Thing not found"
			// Logger.Warn(msg)
			// http.Error(w, msg, http.StatusNotFound)
			// return
		case err != nil:
			msg := "Error creating \"thing\" object: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		default:
			Logger.Info("Success. Thing data received from the database. Thing object created")
			Logger.Info("Serializing data of " + thingType + " thing")
			jsonData, err := json.Marshal(thing)
			if err != nil {
				msg := "JSON serialization error: " + err.Error()
				Logger.Error(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
			w.Write(jsonData)
			return
		}
	case "found":
		var thing types.FoundThing
		err := row.Scan(
			&thing.FoundThingId,
			&thing.PublicationDatetime,
			&thing.ThingName,
			&thing.ThingLocation,
			&thing.CustomText,
			&thing.Verified,
			&thing.Status,
		)
		// TODO: refactor, the code is duplicated
		switch {
		case err == sql.ErrNoRows:
			msg := "Thing not found"
			Logger.Info(msg)
			w.Header().Set("Content-Type", "text/plain")
			// TODO: refactor?
			w.Write([]byte(msg))
			return
		case err != nil:
			msg := "Error creating \"thing\" object: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		default:
			Logger.Info("Success. Thing data received from the database. Thing object created")
			Logger.Info("Serializing data of " + thingType + " thing")
			jsonData, err := json.Marshal(thing)
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
