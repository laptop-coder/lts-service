package handlers

import (
	. "backend/database"
	. "backend/logger"
	"backend/types"
	. "backend/utils"
	"encoding/json"
	"net/http"
)

func GetThingsListMy(w http.ResponseWriter, r *http.Request) {
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

	// Get username from the JWT access
	publicKey, _, err := GetPublicKey()
	if err != nil {
		msg := "Error getting public key: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	accessToken, err := GetJWTAccess(r)
	if err != nil {
		msg := err.Error()
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}
	username, err := GetUsername(accessToken, publicKey)
	if err != nil {
		msg := "Can't get username from JWT access: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Get data from the database
	switch thingsType {
	case "lost":
		rows, err := DB.Query(
			"SELECT * FROM lost_thing WHERE advertisement_owner=? ORDER BY name;",
			username,
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
				&lostThing.PublicationDatetime,
				&lostThing.Name,
				&lostThing.UserMessage,
				&lostThing.Verified,
				&lostThing.Found,
				&lostThing.AdvertisementOwner,
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

	case "found":
		rows, err := DB.Query(
			"SELECT * FROM found_thing WHERE advertisement_owner=? ORDER BY name;",
			username,
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
				&foundThing.PublicationDatetime,
				&foundThing.Name,
				&foundThing.Location,
				&foundThing.UserMessage,
				&foundThing.Verified,
				&foundThing.Found,
				&foundThing.AdvertisementOwner,
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

	default:
		msg := "Error. GET parameter \"things_type\" could be \"lost\" or \"found\""
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
}
