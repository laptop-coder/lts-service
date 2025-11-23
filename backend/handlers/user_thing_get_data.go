package handlers

import (
	. "backend/database"
	. "backend/logger"
	"backend/types"
	. "backend/utils"
	"encoding/json"
	"net/http"
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

	thingType := r.URL.Query().Get("thing_type")
	thingId := r.URL.Query().Get("thing_id")

	if thingType == "" || thingId == "" {
		msg := "Error. GET parameters \"thing_type\" and \"thing_id\" are required"
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
	advertisementEditor, err := GetUsername(accessToken, publicKey)
	if err != nil {
		msg := "Can't get username from JWT access: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Get data from the database
	switch thingType {
	case "lost":
		// Check if advertisement belongs to registered user (compare username
		// in database and username in JWT)
		row := DB.QueryRow(
			"SELECT advertisement_owner FROM lost_thing WHERE id=?;",
			thingId,
		)
		var advertisementOwner string
		err = row.Scan(
			&advertisementOwner,
		)
		if err != nil {
			msg := "Error getting advertisement owner: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		if *advertisementEditor != advertisementOwner {
			msg := "Access denied: it is not your advertisement"
			Logger.Error(msg)
			http.Error(w, msg, http.StatusForbidden)
			return
		}

		rows, err := DB.Query(
			"SELECT * FROM lost_thing WHERE id=?;",
			thingId,
		)
		if err != nil {
			msg := "Error getting lost thing data from the database: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		Logger.Info("Success. Received lost thing data")
		// Serialize data and send it in response
		var lostThing types.LostThing
		for rows.Next() {
			if err := rows.Scan(
				&lostThing.Id,
				&lostThing.PublicationDatetime,
				&lostThing.Name,
				&lostThing.UserEmail,
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
		}
		jsonData, err := json.Marshal(lostThing)
		if err != nil {
			msg := "JSON serialization error: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		w.Write(jsonData)
		return

	case "found":
		// Check if advertisement belongs to registered user (compare username
		// in database and username in JWT)
		row := DB.QueryRow(
			"SELECT advertisement_owner FROM found_thing WHERE id=?;",
			thingId,
		)
		var advertisementOwner string
		err = row.Scan(
			&advertisementOwner,
		)
		if err != nil {
			msg := "Error getting advertisement owner: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		if *advertisementEditor != advertisementOwner {
			msg := "Access denied: it is not your advertisement"
			Logger.Error(msg)
			http.Error(w, msg, http.StatusForbidden)
			return
		}

		rows, err := DB.Query(
			"SELECT * FROM found_thing WHERE id=?;",
			thingId,
		)
		if err != nil {
			msg := "Error getting found thing data from the database: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		Logger.Info("Success. Received found thing data")
		// Serialize data and send it in response
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
		}
		jsonData, err := json.Marshal(foundThing)
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
