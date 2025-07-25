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
	if r.Method == http.MethodGet {
		thingsType := r.URL.Query().Get("things_type")
		// Parameter is empty
		if thingsType == "" {
			Logger.Error("Error. Send request with \"?things_type=lost\" or \"?things_type=found\"")
			fmt.Fprintf(w, "Error. Send request with \"?things_type=lost\" or \"?things_type=found\"")
			return
		}
		// Parameter is incorrect
		if thingsType != "lost" && thingsType != "found" {
			Logger.Error("Error. Things type should be \"lost\" or \"found\"")
			fmt.Fprintf(w, "Error. Things type should be \"lost\" or \"found\"")
			return
		}
		// Get data from the database
		sqlQuery := fmt.Sprintf("SELECT * FROM %s_thing WHERE status=0 ORDER BY id DESC;", thingsType)
		if rows, err := DB.Query(sqlQuery); err != nil {
			Logger.Error("Error getting things list from the database: " + err.Error())
			fmt.Fprintf(w, "Error getting things list from the database: %s", err.Error())
			return
		} else {
			Logger.Info(fmt.Sprintf("Success. Received list of %s things", thingsType))
			// Serialize data and send it in response
			if thingsType == "lost" {
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
						Logger.Error("Error creating \"thing\" object: " + err.Error())
						fmt.Fprintf(w, "Error creating \"thing\" object: %s", err.Error())
						return
					}
					lostThingsList = append(lostThingsList, thing)
				}
				jsonData, err := json.Marshal(lostThingsList)
				if err != nil {
					Logger.Error("Json serialization error: " + err.Error())
					fmt.Fprintf(w, "Json serialization error: %s", err.Error())
					return
				}
				// TODO: refactor. Response in JSON format
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonData)
				return
			} else {
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
						Logger.Error("Error creating \"thing\" object: " + err.Error())
						fmt.Fprintf(w, "Error creating \"thing\" object: %s", err.Error())
						return
					}
					foundThingsList = append(foundThingsList, thing)
				}
				jsonData, err := json.Marshal(foundThingsList)
				if err != nil {
					Logger.Error("Json serialization error: " + err.Error())
					fmt.Fprintf(w, "Json serialization error: %s", err.Error())
					return
				}
				// TODO: refactor. Response in JSON format
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonData)
				return
			}

		}
	} else {
		fmt.Fprintf(w, "A GET request is required")
		Logger.Warn("A GET request is required")
		return
	}
}
