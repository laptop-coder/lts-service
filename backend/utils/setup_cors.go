package utils

import "net/http"

func SetupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST")
	(*w).Header().Set("Access-Control-Allow-Origin", "https://172.16.1.3")
}
