package helpers

import (
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

func ErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	JsonResponse(w, map[string]string{
		"error": message,
	}, statusCode)
}

func SuccessResponse(w http.ResponseWriter, data interface{}) {
	JsonResponse(w, map[string]interface{}{
		"success": true,
		"data":    data,
	}, http.StatusOK)
}

func HandleServiceError(w http.ResponseWriter, err error) {
	errMsg := err.Error()
	statusCode := http.StatusInternalServerError
	switch {
	case strings.Contains(errMsg, "not found"):
		statusCode = http.StatusNotFound
	case strings.Contains(errMsg, "already exists"):
		statusCode = http.StatusConflict
	case strings.Contains(errMsg, "validation"),
		strings.Contains(errMsg, "invalid"),
		strings.Contains(errMsg, "required"):
		statusCode = http.StatusBadRequest
	case strings.Contains(errMsg, "unauthorized"),
		strings.Contains(errMsg, "permission"):
		statusCode = http.StatusUnauthorized
	}
	ErrorResponse(w, errMsg, statusCode)
}

func GetCookie(cookieKey string, r *http.Request) (string, error) {
	cookie, err := r.Cookie(cookieKey)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
