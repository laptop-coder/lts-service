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

func GetThingsListWithoutAuth(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w)
	if r.Method != http.MethodGet {
		msg := "A GET request is required"
		Logger.Warn(msg)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	thingsType := r.URL.Query().Get("things_type")

	// Get data from the database
	var rows *sql.Rows
	var err error
	switch thingsType {
	case "all":
		rows, err = DB.Query(
			"SELECT * FROM thing WHERE verified=1 ORDER BY publication_datetime DESC;",
		)
	case "lost", "found":
		rows, err = DB.Query(
			"SELECT * FROM thing WHERE type=? AND verified=1 ORDER BY publication_datetime DESC;",
			thingsType,
		)
	default:
		msg := "Error. GET parameter \"things_type\" must be \"lost\", \"found\" or \"all\""
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if err != nil {
		msg := "Error getting things list (without auth) from the database: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	Logger.Info("Success. Received things list (without auth)")
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
