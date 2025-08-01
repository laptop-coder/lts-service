package main

import (
	"backend/config"
	. "backend/database"
	"backend/handlers"
	. "backend/logger"
	"net/http"
)

func main() {
	defer DB.Close()
	cfg := config.New()

	http.HandleFunc("/thing/change_status", handlers.ChangeThingStatus)
	http.HandleFunc("/thing/get_data", handlers.GetThingData)
	http.HandleFunc("/things/get_list", handlers.GetThingsList)
	http.HandleFunc("/thing/add", handlers.AddThing)

	Logger.Info("Starting server")
	err := http.ListenAndServeTLS(":443", cfg.SSL.PathToCert, cfg.SSL.PathToKey, nil)
	if err != nil {
		Logger.Error("Error starting the server: " + err.Error())
	}
}
