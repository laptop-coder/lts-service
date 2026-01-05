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

	thingId := r.URL.Query().Get("thing_id")

	if thingId == "" {
		msg := "Error. GET parameter \"thing_id\" is required"
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
	noticeEditor, err := GetUsername(accessToken, publicKey)
	if err != nil {
		msg := "Can't get username from JWT access: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Get data from the database
	// Check if notice belongs to registered user (compare username in database
	// and username in JWT)
	row := DB.QueryRow(
		"SELECT notice_owner FROM thing WHERE id=?;",
		thingId,
	)
	var noticeOwner string
	err = row.Scan(
		&noticeOwner,
	)
	if err != nil {
		msg := "Error getting notice owner: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	if *noticeEditor != noticeOwner {
		msg := "Access denied: it is not your notice"
		Logger.Error(msg)
		http.Error(w, msg, http.StatusForbidden)
		return
	}

	rows, err := DB.Query(
		"SELECT * FROM thing WHERE id=?;",
		thingId,
	)
	if err != nil {
		msg := "Error getting thing data from the database: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	Logger.Info("Success. Received thing data")
	// Serialize data and send it in response
	var thing types.Thing
	for rows.Next() {
		if err := rows.Scan(
			&thing.Id,
			&thing.Type,
			&thing.PublicationDatetime,
			&thing.Name,
			&thing.UserMessage,
			&thing.Verified,
			&thing.Found,
			&thing.NoticeOwner,
		); err != nil {
			msg := "Error (\"thing\" object): " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	}
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
