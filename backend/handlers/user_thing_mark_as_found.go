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
	thingType := r.FormValue("thingType")
	thingId := r.FormValue("thingId")
	// TODO: refactor other checks like this (combine several conditions into
	// one)
	if thingType == "" || thingId == "" {
		msg := "Error. POST parameters \"thingType\" and \"thingId\" are required"
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// TODO: check security of all GET and POST parameters (including thingType
	// and thingId)

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
		// Check if advertisement belongs to registered user (compare username
		// in database and username in JWT)
		// TODO: put this code in a separate function (repeated many times)
		row := DB.QueryRow(
			"SELECT advertisement_owner FROM lost_thing WHERE id=?;",
			thingId,
		)
		var advertisementOwner string
		err = row.Scan(
			&advertisementOwner,
		)
		if err != nil {
			// TODO: handle situation when thingId is incorrect (instead of
			// 500 HTTP error)
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

		// Check that the advertisement has been verified
		row = DB.QueryRow(
			"SELECT verified FROM lost_thing WHERE id=?;",
			thingId,
		)
		var verified string
		err = row.Scan(
			&verified,
		)
		if err != nil {
			msg := "Error getting lost thing verification status: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		// TODO: replace "1" with smth like Consts.Verified.True
		if verified != "1" {
			msg := "Access denied: the advertisement has not been verified"
			Logger.Error(msg)
			http.Error(w, msg, http.StatusForbidden)
			return
		}

		// Update advertisement
		if _, err := DB.Exec(
			"UPDATE lost_thing SET found=1 WHERE id=?;",
			thingId,
		); err != nil {
			msg := "Error marking lost thing as found: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		msg := "Success. Marked lost thing as found"
		Logger.Info(msg)
		w.Write([]byte(msg))
	case "found":

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

		// Check that the advertisement has been verified
		row = DB.QueryRow(
			"SELECT verified FROM found_thing WHERE id=?;",
			thingId,
		)
		var verified string
		err = row.Scan(
			&verified,
		)
		if err != nil {
			msg := "Error getting found thing verification status: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		// TODO: replace "1" with smth like Consts.Verified.True
		if verified != "1" {
			msg := "Access denied: the advertisement has not been verified"
			Logger.Error(msg)
			http.Error(w, msg, http.StatusForbidden)
			return
		}

		// Update advertisement
		if _, err := DB.Exec(
			"UPDATE found_thing SET found=1 WHERE id=?;",
			thingId,
		); err != nil {
			msg := "Error marking found thing as found: " + err.Error()
			Logger.Error(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		msg := "Success. Marked found thing as found"
		Logger.Info(msg)
		w.Write([]byte(msg))
	default:
		msg := "Error. POST parameter \"thingType\" can be \"lost\" or \"found\""
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
}
