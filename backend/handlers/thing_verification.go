package handlers

import (
	. "backend/database"
	. "backend/logger"
	. "backend/utils"
	"fmt"
	"io"
	"net/http"
)

func VerifyThing(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	thingId, thingType, action :=
		r.FormValue("thingId"),
		r.FormValue("thingType"),
		r.FormValue("action")
	if thingId == "" || thingType == "" || action == "" {
		msg := "error: the \"thingId\", \"thingType\" and \"action\" parameters are required"
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	if thingType != "lost" && thingType != "found" {
		msg := "Error. Thing type should be \"lost\" or \"found\""
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	if action != "accept" && action != "reject" {
		msg := "Error. Action should be \"accept\" or \"reject\""
		Logger.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// Verify JWT

	publicKey, _, err := GetPublicKey()
	if err != nil {
		msg := "Error getting public key: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	accessToken, err := r.Cookie("jwt_access")
	if err != nil {
		msg := "Error getting JWT access from the cookie: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	if err := VerifyJWTAccess(&accessToken.Value, publicKey); err != nil {
		msg := "JWT access verification error: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Database query

	var newVerificationStatus string // TODO: is it OK? (string instead of int)
	switch action {
	case "accept":
		newVerificationStatus = "1"
	case "reject":
		newVerificationStatus = "-1"
	}

	sqlQuery := fmt.Sprintf("UPDATE %s_thing SET verified=%s WHERE %s_thing_id=%s;",
		thingType,
		newVerificationStatus,
		thingType,
		thingId,
	)
	if _, err := DB.Exec(sqlQuery); err != nil {
		msg := "Thing verification error: " + err.Error()
		Logger.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	} else {
		msg := "Success. If a thing with this id and type exists, its verification status has been updated"
		Logger.Info(msg)
		io.WriteString(w, msg)
		return
	}
}
