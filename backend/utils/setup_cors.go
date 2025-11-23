package utils

import (
	. "backend/config"
	"net/http"
)

func SetupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:"+Cfg.App.PortFrontend)
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}
