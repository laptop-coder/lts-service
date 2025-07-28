package main

import (
	. "backend/database"
	"backend/handlers"
	. "backend/logger"
	// "flag"
	"net/http"
)

func main() {
	defer DB.Close()

	// var pathToCert string
	// var pathToKey string
	// flag.StringVar(&pathToCert, "path-to-cert", "", "Путь к файлу SSL-сертификата")
	// flag.StringVar(&pathToKey, "path-to-key", "", "Путь к файлу SSL-ключа")

	http.HandleFunc("/thing/change_status", handlers.ChangeThingStatus)
	http.HandleFunc("/thing/get_data", handlers.GetThingData)
	http.HandleFunc("/things/get_list", handlers.GetThingsList)
	http.HandleFunc("/thing/add", handlers.AddThing)

	Logger.Info("Starting server")
	// err := http.ListenAndServeTLS(":443", pathToCert, pathToKey, nil)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		Logger.Error("Error starting the server: " + err.Error())
	}
}
