package utils

import "net/http"

func SetupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Origin", "https://ltsservice.ru")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}
