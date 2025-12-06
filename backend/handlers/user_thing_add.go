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
	noticeOwner, err := GetUsername(accessToken, publicKey)
	if err != nil {
		msg := "Can't get username from JWT access: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	thingId := uuid.New().String()

	thingName := r.FormValue("thingName")
	thingType := r.FormValue("thingType")
	userMessage := r.FormValue("userMessage")
	thingPhoto := r.FormValue("thingPhoto")
	if thingName == "" {
		msg := "Error. POST parameter \"thingName\" is required"
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	if thingType == "" {
		msg := "Error. POST parameter \"thingType\" is required"
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

	// Thing type
	isSecure, err = CheckStringSecurity(thingType)
	if err != nil {
		Logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !(*isSecure) {
		msg := "Error. Found forbidden symbols in POST parameter \"thingType\"."
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// User message
	isSecure, err = CheckStringSecurity(userMessage)
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
	isSecure, err = CheckStringSecurity(thingPhoto)
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

	if _, err := DB.Exec(
		"INSERT INTO thing (id, type, publication_datetime, name, user_message, verified, found, notice_owner) VALUES (?, ?, datetime('now'), ?, ?, ?, ?, ?);",
		thingId,
		thingType,
		thingName,
		userMessage,
		0,
		0,
		noticeOwner,
	); err != nil {
		msg := "Error adding new thing: " + err.Error()
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

	msg := "Success. Added new thing"
	Logger.Info(msg)
	w.Write([]byte(msg))

}
