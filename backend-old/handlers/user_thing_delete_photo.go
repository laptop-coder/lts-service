package handlers

import (
	. "backend/config"
	. "backend/database"
	. "backend/logger"
	. "backend/utils"
	"fmt"
	"net/http"
)

func DeleteThingPhoto(w http.ResponseWriter, r *http.Request) {
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

	// Check if advertisement belongs to registered user (compare username
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

	pathToPhoto := fmt.Sprintf(
		"%s/%s.jpeg",
		Cfg.Storage.PathTo,
		thingId,
	)
	if err := DeleteThingPhotoFromStorageIfExists(pathToPhoto); err != nil {
		msg := "Error deleting thing photo from storage: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	msg := "Success. If a thing with the passed \"thingId\" had a photo, it was deleted"
	Logger.Info(msg)
	w.Write([]byte(msg))
}
