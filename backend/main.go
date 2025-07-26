package main

import (
	. "backend/database"
	"backend/handlers"
	. "backend/logger"
	"net/http"
)

func main() {
	defer DB.Close()

	http.HandleFunc("/thing/change_status", handlers.ChangeThingStatus)
	http.HandleFunc("/thing/get_data", handlers.GetThingData)
	http.HandleFunc("/things/get_list", handlers.GetThingsList)
	Logger.Info("Starting server at port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		Logger.Error("Error starting the server: " + err.Error())
	}
}
