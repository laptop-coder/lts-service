package handlers

import (
	. "backend/database"
	. "backend/logger"
	. "backend/utils"
	"net/http"
)

func MarkThingAsFound(w http.ResponseWriter, r *http.Request) {
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
	// TODO: refactor other checks like this (combine several conditions into
	// one)
	// MAYBE DEPRECATED
	if thingId == "" {
		msg := "Error. POST parameter \"thingId\" is required"
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// TODO: check security of all GET and POST parameters (including thingType
	// and thingId)
	// MAYBE DEPRECATED

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

	// Check if notice belongs to registered user (compare username
	// in database and username in JWT)
	// TODO: put this code in a separate function (repeated many times)
	// MAYBE DEPRECATED
	row := DB.QueryRow(
		"SELECT notice_owner FROM lost_thing WHERE id=?;",
		thingId,
	)
	var noticeOwner string
	err = row.Scan(
		&noticeOwner,
	)
	if err != nil {
		// TODO: handle situation when thingId is incorrect (instead of
		// 500 HTTP error)
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

	// Check that the notice has been verified
	row = DB.QueryRow(
		"SELECT verified FROM thing WHERE id=?;",
		thingId,
	)
	var verified string
	err = row.Scan(
		&verified,
	)
	if err != nil {
		msg := "Error getting thing verification status: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	// TODO: replace "1" with smth like Consts.Verified.True
	if verified != "1" {
		msg := "Access denied: the notice has not been verified"
		Logger.Error(msg)
		http.Error(w, msg, http.StatusForbidden)
		return
	}

	// Update notice
	if _, err := DB.Exec(
		"UPDATE thing SET found=1 WHERE id=?;",
		thingId,
	); err != nil {
		msg := "Error marking thing as found: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	msg := "Success. Marked thing as found"
	Logger.Info(msg)
	w.Write([]byte(msg))
}
