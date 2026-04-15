package helpers

import (
	"backend/pkg/logger"
	"encoding/json"
	"net/http"
	"strings"
)

func JsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	// Set JSON content type
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// Convert data to JSON format
	encodedData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"error": "Failed to encode response to JSON",
		}); err != nil {
			panic(err)
		}
		return
	}
	// Status code
	w.WriteHeader(statusCode)
	// Response
	w.Write(encodedData)
}

func ErrorResponse(log logger.Logger, w http.ResponseWriter, message string, statusCode int) {
	log.Error(message, "status_code", statusCode)
	JsonResponse(w, map[string]string{
		"error": message,
	}, statusCode)
}

func SuccessResponse(w http.ResponseWriter, data interface{}) {
	JsonResponse(w, data, http.StatusOK)
}

func HandleServiceError(log logger.Logger, w http.ResponseWriter, err error) {
	errMsg := err.Error()
	statusCode := http.StatusInternalServerError
	switch {
	case strings.Contains(errMsg, "not found"):
		statusCode = http.StatusNotFound
	case strings.Contains(errMsg, "validation"),
		strings.Contains(errMsg, "invalid"),
		strings.Contains(errMsg, "bad request"),
		strings.Contains(errMsg, "required"):
		statusCode = http.StatusBadRequest
	case strings.Contains(errMsg, "unauthorized"),
		strings.Contains(errMsg, "expired"):
		statusCode = http.StatusUnauthorized
	case strings.Contains(errMsg, "conflict"),
		strings.Contains(errMsg, "already"):
		statusCode = http.StatusConflict
	case strings.Contains(errMsg, "forbidden"),
		strings.Contains(errMsg, "permission"),
		strings.Contains(errMsg, "revoked"),
		strings.Contains(errMsg, "verify"):
		statusCode = http.StatusForbidden
	}
	ErrorResponse(log, w, errMsg, statusCode)
}

func GetCookie(cookieKey string, r *http.Request) (string, error) {
	cookie, err := r.Cookie(cookieKey)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
