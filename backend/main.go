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
	utils.GenKeysIfNotExist()

	mux := http.NewServeMux()

	// For all users
	mux.Handle("/user/register", http.HandlerFunc(handlers.UserRegister))
	mux.Handle("/user/login", http.HandlerFunc(handlers.UserLogin))
	mux.Handle("/user/logout", http.HandlerFunc(handlers.UserLogout))
	mux.Handle("/moderator/register", http.HandlerFunc(handlers.ModeratorRegister))
	mux.Handle("/moderator/login", http.HandlerFunc(handlers.ModeratorLogin))
	mux.Handle("/moderator/logout", http.HandlerFunc(handlers.ModeratorLogout))
	mux.Handle("/things/get_list", http.HandlerFunc(handlers.GetThingsList))

	// For registered users
	mux.Handle("/thing/get_data", utils.AuthMiddleware(&Cfg.Role.User, http.HandlerFunc(handlers.GetThingData)))
	mux.Handle("/thing/add", utils.AuthMiddleware(&Cfg.Role.User, http.HandlerFunc(handlers.AddThing)))
	mux.Handle("/thing/edit", utils.AuthMiddleware(&Cfg.Role.User, http.HandlerFunc(handlers.EditThing)))
	mux.Handle("/thing/delete", utils.AuthMiddleware(&Cfg.Role.User, http.HandlerFunc(handlers.DeleteThing)))
	mux.Handle("/thing/delete_photo", utils.AuthMiddleware(&Cfg.Role.User, http.HandlerFunc(handlers.DeleteThingPhoto)))
	mux.Handle("/thing/mark_as_found", utils.AuthMiddleware(&Cfg.Role.User, http.HandlerFunc(handlers.MarkThingAsFound)))

	// For moderators
	mux.Handle("/thing/verify", utils.AuthMiddleware(&Cfg.Role.Moderator, http.HandlerFunc(handlers.VerifyThing)))

	Logger.Info("Starting server via HTTP...")
	err := http.ListenAndServe(":"+Cfg.App.PortBackend, mux)
	if err != nil {
		Logger.Error("Error starting the server: " + err.Error())
	}
}
