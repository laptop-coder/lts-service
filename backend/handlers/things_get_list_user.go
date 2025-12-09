package handlers

import (
	. "backend/database"
	. "backend/logger"
	"backend/types"
	. "backend/utils"
	"database/sql"
	"encoding/json"
	"net/http"
)

func GetThingsListUser(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w)
	if r.Method != http.MethodGet {
		msg := "A GET request is required"
		Logger.Warn(msg)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	thingsType := r.URL.Query().Get("things_type")
	noticesOwnership := r.URL.Query().Get("notices_ownership")

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
	// TODO: refactor
	var rows *sql.Rows
	switch thingsType {
	case "all":
		switch noticesOwnership {
		case "all":
			rows, err = DB.Query(
				"SELECT * FROM thing WHERE (NOT notice_owner=? and verified=1) or notice_owner=? ORDER BY publication_datetime DESC;",
				*username,
				*username,
			)
		case "my":
			rows, err = DB.Query(
				"SELECT * FROM thing WHERE notice_owner=? ORDER BY publication_datetime DESC;",
				*username,
			)
		case "not_my":
			rows, err = DB.Query(
				"SELECT * FROM thing WHERE NOT notice_owner=? AND verified=1 ORDER BY publication_datetime DESC;",
				*username,
			)
		default:
			msg := "Error. GET parameter \"notices_ownership\" must be \"my\", \"not_my\" or \"all\""
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
	case "lost", "found":
		switch noticesOwnership {
		case "all":
			rows, err = DB.Query(
				"SELECT * FROM thing WHERE type=? AND ((NOT notice_owner=? and verified=1) or notice_owner=?) ORDER BY publication_datetime DESC;",
				thingsType,
				*username,
				*username,
			)
		case "my":
			rows, err = DB.Query(
				"SELECT * FROM thing WHERE type=? AND notice_owner=? ORDER BY publication_datetime DESC;",
				thingsType,
				*username,
			)
		case "not_my":
			rows, err = DB.Query(
				"SELECT * FROM thing WHERE type=? AND NOT notice_owner=? AND verified=1 ORDER BY publication_datetime DESC;",
				thingsType,
				*username,
			)
		default:
			msg := "Error. GET parameter \"notices_ownership\" must be \"my\", \"not_my\" or \"all\""
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
	default:
		msg := "Error. GET parameter \"things_type\" must be \"lost\", \"found\" or \"all\""
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if err != nil {
		msg := "Error getting things list from the database: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	Logger.Info("Success. Received things list")
	// Serialize data and send it in response
	var thingsList []types.Thing
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
		thingsList = append(thingsList, thing)
	}
	jsonData, err := json.Marshal(thingsList)
	if err != nil {
		msg := "JSON serialization error: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}
