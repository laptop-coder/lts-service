package handler

import (
	"encoding/json"
	"net/http"
	"strings"
)

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"error": "Failed to encode response"}`, http.StatusInternalServerError)
	}
}

func errorResponse(w http.ResponseWriter, message string, statusCode int) {
	jsonResponse(w, map[string]string{
		"error": message,
	}, statusCode)
}

func successResponse(w http.ResponseWriter, data interface{}) {
	jsonResponse(w, map[string]interface{}{
		"success": true,
		"data":    data,
	}, http.StatusOK)
}

func handleServiceError(w http.ResponseWriter, err error) {
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

	errorResponse(w, errMsg, statusCode)
}
