package handlers

import (
	. "backend/database"
	. "backend/logger"
	"backend/types"
	. "backend/utils"
	"encoding/json"
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

	if thingsType == "" {
		msg := "Error. GET parameter \"things_type\" is required"
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// Get data from the database

	if thingsType == "lost" {
		rows, err := DB.Query(
			"SELECT * FROM lost_thing ORDER BY name;",
		)
		if err != nil {
			msg := "Error getting lost things list from the database: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		Logger.Info("Success. Received lost things list")
		// Serialize data and send it in response
		var lostThingsList []types.LostThing
		var lostThing types.LostThing
		for rows.Next() {
			if err := rows.Scan(
				&lostThing.Id,
				&lostThing.Name,
				&lostThing.UserEmail,
				&lostThing.UserMessage,
				&lostThing.Verified,
				&lostThing.Found,
			); err != nil {
				msg := "Error (\"lost thing\" object): " + err.Error()
				Logger.Error(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
			lostThingsList = append(lostThingsList, lostThing)
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

	} else if thingsType == "found" {
		rows, err := DB.Query(
			"SELECT * FROM found_thing ORDER BY name;",
		)
		if err != nil {
			msg := "Error getting found things list from the database: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		Logger.Info("Success. Received found things list")
		// Serialize data and send it in response
		var foundThingsList []types.FoundThing
		var foundThing types.FoundThing
		for rows.Next() {
			if err := rows.Scan(
				&foundThing.Id,
				&foundThing.Name,
				&foundThing.Location,
				&foundThing.UserMessage,
				&foundThing.Verified,
				&foundThing.Found,
			); err != nil {
				msg := "Error (\"found thing\" object): " + err.Error()
				Logger.Error(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
			foundThingsList = append(foundThingsList, foundThing)
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

	} else {
		msg := "Error. GET parameter \"things_type\" could be \"lost\" or \"found\""
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
}
