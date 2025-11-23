package handlers

import (
	. "backend/database"
	. "backend/logger"
	. "backend/utils"
	"net/http"
)

func EditThing(w http.ResponseWriter, r *http.Request) {
	SetupCORS(&w)
	if r.Method != http.MethodPost {
		msg := "A POST request is required"
		Logger.Warn(msg)
		http.Error(w, msg, http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		msg := "Error. Can't parse form: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	thingType := r.FormValue("thingType")
	if thingType == "" {
		msg := "Error. POST parameter \"thingType\" is required"
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	thingId := r.FormValue("thingId")
	if thingId == "" {
		msg := "Error. POST parameter \"thingId\" is required"
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

	switch thingType {
	case "lost":
		newThingName := r.FormValue("newThingName")
		newUserEmail := r.FormValue("newUserEmail")
		newUserMessage := r.FormValue("newUserMessage")
		newThingPhoto := r.FormValue("newThingPhoto")

		if newThingName == "" || newUserEmail == "" {
			msg := "Error. \"thingType\" is \"lost\", so POST parameters \"newThingName\" and \"newUserEmail\" are required"
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// Regular expressions checks
		// Thing name
		isSecure, err := CheckStringSecurity(newThingName)
		if err != nil {
			Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !(*isSecure) {
			msg := "Error. Found forbidden symbols in POST parameter \"newThingName\"."
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// User email
		isSecure, err = CheckStringSecurity(newUserEmail)
		if err != nil {
			Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !(*isSecure) {
			msg := "Error. Found forbidden symbols in POST parameter \"newUserEmail\"."
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// User message
		isSecure, err = CheckStringSecurity(newUserMessage)
		if err != nil {
			Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !(*isSecure) {
			msg := "Error. Found forbidden symbols in POST parameter \"newUserMessage\"."
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// Thing photo
		isSecure, err = CheckStringSecurity(newThingPhoto)
		if err != nil {
			Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !(*isSecure) {
			msg := "Error. Found forbidden symbols in POST parameter \"newThingPhoto\"."
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

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

		// Update advertisement
		if _, err := DB.Exec(
			"UPDATE lost_thing SET name=?, user_email=?, user_message=? WHERE id=?;",
			newThingName,
			newUserEmail,
			newUserMessage,
			thingId,
		); err != nil {
			msg := "Error editing lost thing: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		if newThingPhoto != "" {
			if err := SaveThingPhotoToStorage(newThingPhoto, thingId); err != nil {
				msg := "Error updating thing photo in storage: " + err.Error()
				Logger.Error(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
		}

		msg := "Success. Edited lost thing"
		Logger.Info(msg)
		w.Write([]byte(msg))
	case "found":
		newThingName := r.FormValue("newThingName")
		newThingLocation := r.FormValue("newThingLocation")
		newUserMessage := r.FormValue("newUserMessage")
		newThingPhoto := r.FormValue("newThingPhoto")

		if newThingName == "" || newThingLocation == "" {
			msg := "Error. \"thingType\" is \"found\", so POST parameters \"newThingName\" and \"newThingLocation\" are required"
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// Regular expressions checks
		// Thing name
		isSecure, err := CheckStringSecurity(newThingName)
		if err != nil {
			Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !(*isSecure) {
			msg := "Error. Found forbidden symbols in POST parameter \"newThingName\"."
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// Thing location
		isSecure, err = CheckStringSecurity(newThingLocation)
		if err != nil {
			Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !(*isSecure) {
			msg := "Error. Found forbidden symbols in POST parameter \"newThingLocation\"."
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// User message
		isSecure, err = CheckStringSecurity(newUserMessage)
		if err != nil {
			Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !(*isSecure) {
			msg := "Error. Found forbidden symbols in POST parameter \"newUserMessage\"."
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// Thing photo
		isSecure, err = CheckStringSecurity(newThingPhoto)
		if err != nil {
			Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !(*isSecure) {
			msg := "Error. Found forbidden symbols in POST parameter \"newThingPhoto\"."
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

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

		// Update advertisement
		if _, err := DB.Exec(
			"UPDATE found_thing SET name=?, location=?, user_message=? WHERE id=?;",
			newThingName,
			newThingLocation,
			newUserMessage,
			thingId,
		); err != nil {
			msg := "Error editing found thing: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		if newThingPhoto != "" {
			if err := SaveThingPhotoToStorage(newThingPhoto, thingId); err != nil {
				msg := "Error updating thing photo in storage: " + err.Error()
				Logger.Error(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
		}

		msg := "Success. Edited found thing"
		Logger.Info(msg)
		w.Write([]byte(msg))
	default:
		msg := "Error. POST parameter \"thingType\" can be \"lost\" or \"found\""
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
}
