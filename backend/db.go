package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func initDB() *sql.DB {
	logger := initLogger()
	const PATH_TO_DB string = "./db.sqlite3" // TODO: move const
	db, err := sql.Open("sqlite3", PATH_TO_DB)
	if err != nil {
		logger.Error("Error. Can't open database file: " + err.Error())
	} else {
		logger.Info("The database file is open")
	}
	if err := db.Ping(); err != nil {
		logger.Error("Error. Can't connect to the database: " + err.Error())
	} else {
		logger.Info("Pinged successfully. Can connect to the database")
	}
	return db
}

var db = initDB()
