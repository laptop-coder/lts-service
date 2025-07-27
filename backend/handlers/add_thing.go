// TODO: is it error? When you send base64 data without data:image, function to
// save photo to the storage is not running, but the database query has already
// been completed, so the thing photo was not saved, but the information of the
// thing was added to the database
//
// POSSIBLY OUTDATED
package handlers

import (
	. "backend/database"
	. "backend/logger"
	. "backend/utils"
	"fmt"
	"net/http"
)

func AddThing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		msg := "A POST request is required"
		Logger.Warn(msg)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}

	thingType := r.URL.Query().Get("thing_type")
	// Parameter is empty
	if thingType == "" {
		msg := "Error. Send POST request with \"?thing_type=\" (can be \"lost\" or \"found\") parameter"
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	// Parameter is incorrect
	if thingType != "lost" && thingType != "found" {
		msg := "Error. Thing type should be \"lost\" or \"found\""
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		msg := "Error. Can't parse form: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	var sqlQuery string
	var thingPhoto string

	switch thingType {
	case "lost":
		thingName := r.FormValue("thing_name")
		email := r.FormValue("email")
		customText := r.FormValue("custom_text")
		thingPhoto = r.FormValue("thing_photo")
		if thingName == "" || email == "" || customText == "" {
			msg := "Error. Thing type is \"lost\", so the POST parameters thing_name, email and custom_text are required"
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		sqlQuery = fmt.Sprintf(`
            INSERT INTO lost_thing (
                publication_date,
                publication_time,
                thing_name,
                email,
                custom_text,
                verified,
                status
            )
            VALUES (
				date('now'), substr(time('now'), 1, 5), '%s', '%s', '%s', 0, 0
            );
		`, thingName, email, customText)
	case "found":
		thingName := r.FormValue("thing_name")
		thingLocation := r.FormValue("thing_location")
		customText := r.FormValue("custom_text")
		thingPhoto = r.FormValue("thing_photo")
		if thingName == "" || thingLocation == "" || customText == "" {
			msg := "Error. Thing type is \"found\", so the POST parameters thing_name, thing_location and custom_text are required"
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		sqlQuery = fmt.Sprintf(`
            INSERT INTO found_thing (
                publication_date,
                publication_time,
                thing_name,
                thing_location,
                custom_text,
                verified,
                status
            )
            VALUES (
                date('now'), substr(time('now'), 1, 5), '%s', '%s', '%s', 0, 0
            );
		`, thingName, thingLocation, customText)
	}

	sqlResult, err := DB.Exec(sqlQuery)
	if err != nil {
		msg := "Error adding new " + thingType + " thing: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	if thingPhoto != "" {
		thingId, err := sqlResult.LastInsertId()
		if err != nil {
			msg := "Error getting id of the added " + thingType + " thing: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		if err := SaveThingPhotoToStorage(thingPhoto, thingId, thingType); err != nil {
			msg := "Error saving thing photo to the storage: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	}

	msg := "Success. Added a new " + thingType + " thing"
	Logger.Info(msg)
	w.Write([]byte(msg))
}
