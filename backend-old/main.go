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
	// TODO: refactor handlers (like in frontend)

	mux := http.NewServeMux()

	// For all users
	mux.Handle("/user/register", http.HandlerFunc(handlers.UserRegister))
	mux.Handle("/user/login", http.HandlerFunc(handlers.UserLogin))
	mux.Handle("/moderator/register", http.HandlerFunc(handlers.ModeratorRegister))
	mux.Handle("/moderator/login", http.HandlerFunc(handlers.ModeratorLogin))
	mux.Handle("/logout", http.HandlerFunc(handlers.Logout))
	mux.Handle("/things/get_list/without_auth", http.HandlerFunc(handlers.GetThingsListWithoutAuth))

	// For registered users
	mux.Handle("/user/get_username", utils.AuthMiddleware(&Cfg.Role.User, http.HandlerFunc(handlers.UserGetUsername)))
	mux.Handle("/user/get_email", utils.AuthMiddleware(&Cfg.Role.User, http.HandlerFunc(handlers.UserGetEmail)))
	mux.Handle("/user/get_email/other", utils.AuthMiddleware(&Cfg.Role.User, http.HandlerFunc(handlers.UserOtherGetEmail)))
	mux.Handle("/things/get_list/user", utils.AuthMiddleware(&Cfg.Role.User, http.HandlerFunc(handlers.GetThingsListUser)))
	mux.Handle("/thing/get_data", utils.AuthMiddleware(&Cfg.Role.User, http.HandlerFunc(handlers.GetThingData)))
	mux.Handle("/thing/add", utils.AuthMiddleware(&Cfg.Role.User, http.HandlerFunc(handlers.AddThing)))
	mux.Handle("/thing/edit", utils.AuthMiddleware(&Cfg.Role.User, http.HandlerFunc(handlers.EditThing)))
	mux.Handle("/thing/delete/user", utils.AuthMiddleware(&Cfg.Role.User, http.HandlerFunc(handlers.UserDeleteThing)))
	mux.Handle("/thing/delete_photo", utils.AuthMiddleware(&Cfg.Role.User, http.HandlerFunc(handlers.DeleteThingPhoto)))
	mux.Handle("/thing/mark_as_found", utils.AuthMiddleware(&Cfg.Role.User, http.HandlerFunc(handlers.MarkThingAsFound)))

	// For moderators
	mux.Handle("/things/get_list/moderator", utils.AuthMiddleware(&Cfg.Role.Moderator, http.HandlerFunc(handlers.GetThingsListModerator)))
	mux.Handle("/moderator/get_username", utils.AuthMiddleware(&Cfg.Role.Moderator, http.HandlerFunc(handlers.ModeratorGetUsername)))
	mux.Handle("/thing/change_verification", utils.AuthMiddleware(&Cfg.Role.Moderator, http.HandlerFunc(handlers.ThingChangeVerification)))
	mux.Handle("/thing/delete/moderator", utils.AuthMiddleware(&Cfg.Role.Moderator, http.HandlerFunc(handlers.ModeratorDeleteThing)))

	Logger.Info("Starting server via HTTP...")
	err := http.ListenAndServe(":"+Cfg.App.PortBackend, mux)
	if err != nil {
		Logger.Error("Error starting the server: " + err.Error())
	}
}
