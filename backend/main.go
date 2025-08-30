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
	defer Valkey.Close()
	utils.GenKeysIfNotExist()

	http.HandleFunc("/thing/add", handlers.AddThing)
	http.HandleFunc("/thing/change_status", handlers.ChangeThingStatus)
	http.HandleFunc("/thing/get_data", handlers.GetThingData)
	http.HandleFunc("/things/get_list", handlers.GetThingsList)
	http.HandleFunc("/moderator/register", handlers.ModeratorRegister)
	http.HandleFunc("/moderator/login", handlers.ModeratorLogin)

	Logger.Info("Starting server")
	err := http.ListenAndServeTLS(":443", Cfg.SSL.PathToCert, Cfg.SSL.PathToKey, nil)
	if err != nil {
		Logger.Error("Error starting the server: " + err.Error())
	}
}
