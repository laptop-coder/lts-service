package main

import (
	. "backend/logger"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var initialQueries = `
CREATE TABLE IF NOT EXISTS lost_thing (
    id INTEGER PRIMARY KEY,
    publication_date TEXT NOT NULL,
    publication_time TEXT NOT NULL,
    thing_name TEXT NOT NULL,
    email varchar(254) NOT NULL,
    custom_text TEXT NOT NULL,
    verified INTEGER NOT NULL,
    status INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS found_thing (
    id INTEGER PRIMARY KEY,
    publication_date TEXT NOT NULL,
    publication_time TEXT NOT NULL,
    thing_name TEXT NOT NULL,
    thing_location TEXT NOT NULL,
    custom_text TEXT NOT NULL,
    verified INTEGER NOT NULL,
    status INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS moderator (
    id INTEGER PRIMARY KEY,
    username varchar(32) NOT NULL,
    password TEXT NOT NULL
);
`

func initDB() *sql.DB {
	const PATH_TO_DB string = "./db.sqlite3" // TODO: move const
	db, err := sql.Open("sqlite3", PATH_TO_DB)
	if err != nil {
		Logger.Error("Error. Can't open database file: " + err.Error())
	} else {
		Logger.Info("The database file is open")
	}
	if err := db.Ping(); err != nil {
		Logger.Error("Error. Can't connect to the database: " + err.Error())
	} else {
		Logger.Info("Pinged successfully. Can connect to the database")
	}
	if _, err := db.Exec(initialQueries); err != nil {
		Logger.Error("Error in running initial SQL queries")
	} else {
		Logger.Info("Initial SQL queries completed")
	}

	return db
}

var db = initDB()
