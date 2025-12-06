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
	noticeEditor, err := GetUsername(accessToken, publicKey)
	if err != nil {
		msg := "Can't get username from JWT access: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	newThingName := r.FormValue("newThingName")
	newThingType := r.FormValue("newThingType")
	newUserMessage := r.FormValue("newUserMessage")
	newThingPhoto := r.FormValue("newThingPhoto")

	if newThingType == "" {
		msg := "Error. POST parameter \"newThingType\" is required"
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	if newThingType != "lost" && newThingType != "found" {
		msg := "Error. POST parameter \"newThingType\" must be \"lost\" or \"found\""
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if newThingName == "" {
		msg := "Error. \"thingType\" is \"lost\", so POST parameter \"newThingName\" is required"
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

	// Thing type
	isSecure, err = CheckStringSecurity(newThingType)
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !(*isSecure) {
		msg := "Error. Found forbidden symbols in POST parameter \"newThingType\"."
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

	// Check if notice belongs to registered user (compare username
	// in database and username in JWT)
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

	// Update notice
	if _, err := DB.Exec(
		"UPDATE thing SET name=?, type=?, user_message=? WHERE id=?;",
		newThingName,
		newThingType,
		newUserMessage,
		thingId,
	); err != nil {
		msg := "Error editing thing: " + err.Error()
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

	msg := "Success. Edited thing"
	Logger.Info(msg)
	w.Write([]byte(msg))
}
