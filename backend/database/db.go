package database

import (
	. "backend/config"
	. "backend/logger"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var initialQueries = `
CREATE TABLE IF NOT EXISTS lost_thing (
    lost_thing_id VARCHAR(36) PRIMARY KEY,
    publication_datetime DATETIME,
    thing_name TEXT NOT NULL,
    user_email VARCHAR(254) NOT NULL,
    custom_text TEXT NOT NULL,
    verified INTEGER NOT NULL,
    status INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS found_thing (
    found_thing_id VARCHAR(36) PRIMARY KEY,
    publication_datetime DATETIME,
    thing_name TEXT NOT NULL,
    thing_location TEXT NOT NULL,
    custom_text TEXT NOT NULL,
    verified INTEGER NOT NULL,
    status INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS moderator (
    moderator_id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(32) NOT NULL UNIQUE,
    email VARCHAR(254) NOT NULL UNIQUE,
    password VARCHAR(128) NOT NULL,
    credentials_version INTEGER NOT NULL DEFAULT 0
);
`

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", Cfg.DB.PathTo)
	if err != nil {
		Logger.Error("Error. Can't open database file: " + err.Error())
		return nil
	} else {
		Logger.Info("The database file is open")
	}
	if err := db.Ping(); err != nil {
		Logger.Error("Error. Can't connect to the database: " + err.Error())
		return nil
	} else {
		Logger.Info("Pinged successfully. Can connect to the database")
	}
	if _, err := db.Exec(initialQueries); err != nil {
		Logger.Error("Error in running initial SQL queries")
		return nil
	} else {
		Logger.Info("Initial SQL queries completed")
	}

	return db
}

var DB = initDB()
