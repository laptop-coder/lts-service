package handlers

import (
	. "backend/database"
	. "backend/logger"
	"encoding/json"
	"fmt"
	"net/http"
)

type lostThing struct {
	Id              int
	PublicationDate string
	PublicationTime string
	ThingName       string
	Email           string
	CustomText      string
	Verified        int
	Status          int
}

type foundThing struct {
	Id              int
	PublicationDate string
	PublicationTime string
	ThingName       string
	ThingLocation   string
	CustomText      string
	Verified        int
	Status          int
}

func GetThingsList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		msg := "A GET request is required"
		Logger.Warn(msg)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	thingsType := r.URL.Query().Get("things_type")
	// Parameter is empty
	if thingsType == "" {
		msg := "Error. Send request with \"?things_type=lost\" or \"?things_type=found\""
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	// Parameter is incorrect
	if thingsType != "lost" && thingsType != "found" {
		msg := "Error. Things type should be \"lost\" or \"found\""
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	// Get data from the database
	sqlQuery := fmt.Sprintf("SELECT * FROM %s_thing WHERE status=0 ORDER BY id DESC;", thingsType)
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
			var lostThingsList []lostThing
			var thing lostThing
			for rows.Next() {
				if err := rows.Scan(
					&thing.Id,
					&thing.PublicationDate,
					&thing.PublicationTime,
					&thing.ThingName,
					&thing.Email,
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
				msg := "Json serialization error: " + err.Error()
				Logger.Error(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
			w.Write(jsonData)
			return
		case "found":
			var foundThingsList []foundThing
			var thing foundThing
			for rows.Next() {
				if err := rows.Scan(
					&thing.Id,
					&thing.PublicationDate,
					&thing.PublicationTime,
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
				msg := "Json serialization error: " + err.Error()
				Logger.Error(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
			w.Write(jsonData)
			return
		}

	}
}
