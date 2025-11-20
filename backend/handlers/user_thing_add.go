package handlers

import (
	. "backend/database"
	. "backend/logger"
	. "backend/utils"
	"github.com/google/uuid"
	"net/http"
)

func AddThing(w http.ResponseWriter, r *http.Request) {
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
	} else if thingType == "lost" {
		thingName := r.FormValue("thingName")
		userEmail := r.FormValue("userEmail")
		userMessage := r.FormValue("userMessage")
		thingPhoto := r.FormValue("thingPhoto")

		if thingName == "" || userEmail = "" {
			msg := "Error. \"thingType\" is \"lost\", so POST parameters \"thingName\" and \"userEmail\" are required"
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// Regular expressions checks
		// Thing name
		isSecure, err := CheckStringSecurity(thingName)
		if err != nil {
			Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !(*isSecure) {
			msg := "Error. Found forbidden symbols in POST parameter \"thingName\"."
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// User email
		isSecure, err := CheckStringSecurity(userEmail)
		if err != nil {
			Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !(*isSecure) {
			msg := "Error. Found forbidden symbols in POST parameter \"userEmail\"."
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// User message
		isSecure, err := CheckStringSecurity(userMessage)
		if err != nil {
			Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !(*isSecure) {
			msg := "Error. Found forbidden symbols in POST parameter \"userMessage\"."
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// Thing photo
		isSecure, err := CheckStringSecurity(thingPhoto)
		if err != nil {
			Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !(*isSecure) {
			msg := "Error. Found forbidden symbols in POST parameter \"thingPhoto\"."
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// TODO: get advertisement owner from the JWT access here and check it's
		// security (maybe)

		thingId := uuid.New().String()
		if _, err := DB.Exec(
			"INSERT INTO lost_thing (id, name, user_email, user_message, verified, found, advertisement_owner) VALUES (?, ?, ?, ?);",
			thingId,
			thingName,
			userEmail,
			userMessage,
			0,
			0,
			advertisementOwner,
		); err != nil {
			msg := "Error adding new lost thing: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		if thingPhoto != "" {
			if err := SaveThingPhotoToStorage(thingPhoto, thingId); err != nil {
				msg := "Error saving thing photo to storage: " + err.Error()
				Logger.Error(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
		}

		msg := "Success. Added new lost thing"
		Logger.Info(msg)
		w.Write([]byte(msg))
	} else if thingType == "found" {
		thingName := r.FormValue("thingName")
		thingLocation := r.FormValue("thingLocation")
		userMessage := r.FormValue("userMessage")
		thingPhoto := r.FormValue("thingPhoto")

		if thingName == "" || thingLocation = "" {
			msg := "Error. \"thingType\" is \"found\", so POST parameters \"thingName\" and \"thingLocation\" are required"
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// Regular expressions checks
		// Thing name
		isSecure, err := CheckStringSecurity(thingName)
		if err != nil {
			Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !(*isSecure) {
			msg := "Error. Found forbidden symbols in POST parameter \"thingName\"."
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// Thing location
		isSecure, err := CheckStringSecurity(thingLocation)
		if err != nil {
			Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !(*isSecure) {
			msg := "Error. Found forbidden symbols in POST parameter \"thingLocation\"."
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// User message
		isSecure, err := CheckStringSecurity(userMessage)
		if err != nil {
			Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !(*isSecure) {
			msg := "Error. Found forbidden symbols in POST parameter \"userMessage\"."
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// Thing photo
		isSecure, err := CheckStringSecurity(thingPhoto)
		if err != nil {
			Logger.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !(*isSecure) {
			msg := "Error. Found forbidden symbols in POST parameter \"thingPhoto\"."
			Logger.Error(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// TODO: get advertisement owner from the JWT access here and check it's
		// security (maybe)

		thingId := uuid.New().String()
		if _, err := DB.Exec(
			"INSERT INTO found_thing (id, name, thing_location, user_message, verified, found, advertisement_owner) VALUES (?, ?, ?, ?);",
			thingId,
			thingName,
			thingLocation,
			userMessage,
			0,
			0,
			advertisementOwner,
		); err != nil {
			msg := "Error adding new found thing: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		if thingPhoto != "" {
			if err := SaveThingPhotoToStorage(thingPhoto, thingId); err != nil {
				msg := "Error saving thing photo to storage: " + err.Error()
				Logger.Error(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
		}

		msg := "Success. Added new found thing"
		Logger.Info(msg)
		w.Write([]byte(msg))
	} else {
		msg := "Error. POST parameter \"thingType\" can be \"lost\" or \"found\""
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
}
