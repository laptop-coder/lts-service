package main

import (
	. "backend/config"
	. "backend/database"
	"backend/handlers"
	. "backend/logger"
	"backend/utils"
	"net/http"
)

func main() {
	defer DB.Close()
	// defer Valkey.Close()
	utils.GenKeysIfNotExist()

	mux := http.NewServeMux()

	// For all users
	mux.Handle("/thing/add", http.HandlerFunc(handlers.AddThing))
	mux.Handle("/thing/change_status", http.HandlerFunc(handlers.ChangeThingStatus))
	mux.Handle("/thing/verify", http.HandlerFunc(handlers.VerifyThing))
	mux.Handle("/thing/get_data", http.HandlerFunc(handlers.GetThingData))
	mux.Handle("/things/get_list", http.HandlerFunc(handlers.GetThingsList))
	mux.Handle("/moderator/register", http.HandlerFunc(handlers.ModeratorRegister))
	mux.Handle("/moderator/login", http.HandlerFunc(handlers.ModeratorLogin))

	Logger.Info("Starting server via HTTP...")
	err := http.ListenAndServe(":"+Cfg.App.PortBackend, mux)
	if err != nil {
		Logger.Error("Error starting the server: " + err.Error())
	}
}
