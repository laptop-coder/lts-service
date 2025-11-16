package handlers

import (
	. "backend/database"
	. "backend/logger"
	. "backend/utils"
	"fmt"
	"io"
	"net/http"
	"time"
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

	accessTokenCookie, err := r.Cookie("jwt_access")
	var accessToken *string
	if err != nil {
		Logger.Info("Can't get JWT access from the cookie: " + err.Error() + ". Trying to refresh it")
		refreshToken, err := r.Cookie("jwt_refresh")
		if err != nil {
			msg := "Can't get JWT refresh from the cookie: " + err.Error() + ". Please log in by password"
			Logger.Info(msg)
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}
		newAccessToken, err := RefreshJWTAccess(&refreshToken.Value)
		if err != nil {
			msg := "Error refreshing JWT access: " + err.Error()
			Logger.Info(msg)
			http.Error(w, msg, http.StatusUnauthorized)
			return
		}
		http.SetCookie(
			w,
			&http.Cookie{
				Name:        "jwt_access",
				Value:       *newAccessToken,
				Secure:      true,
				HttpOnly:    true,
				Partitioned: true,
				SameSite:    http.SameSiteNoneMode,
				Path:        "/", // TODO: is it OK?
				Domain:      "server.ltsservice.ru",
				Expires:     time.Now().Add(5 * time.Minute),
			},
		)
		accessToken = newAccessToken
	} else {
		accessToken = &accessTokenCookie.Value
	}
	if err := VerifyJWTAccess(accessToken, publicKey); err != nil {
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
