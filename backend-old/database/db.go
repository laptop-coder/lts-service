package database

import (
	. "backend/config"
	. "backend/logger"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// TODO: add length restrictions (VARCHAR(n) instead of TEXT)
// TODO: rename thing type to thing category
var initialQueries = `
CREATE TABLE IF NOT EXISTS thing (
    id VARCHAR(36) PRIMARY KEY,
	type VARCHAR(5) CHECK (type IN ('lost', 'found')),
    publication_datetime DATETIME,
    name TEXT NOT NULL,
    user_message TEXT NOT NULL,
    verified INTEGER NOT NULL,
    found INTEGER NOT NULL,
	notice_owner VARCHAR(36) NOT NULL,
	FOREIGN KEY(notice_owner) REFERENCES user(username)
);

CREATE TABLE IF NOT EXISTS user (
    username TEXT PRIMARY KEY,
    email VARCHAR(254) NOT NULL,
	password_hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS moderator (
    username TEXT PRIMARY KEY,
	password_hash TEXT NOT NULL
);

CREATE TRIGGER IF NOT EXISTS limit_moderator_accounts_count
BEFORE INSERT ON moderator
FOR EACH ROW
BEGIN
   SELECT CASE
	   WHEN (SELECT COUNT(*) FROM moderator) = 1 THEN
		   RAISE(ABORT, 'the moderator account has already been created (you can''t create more than one account)')
   end;
end;
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
		Logger.Error("Error in running initial SQL queries: " + err.Error())
		return nil
	} else {
		Logger.Info("Initial SQL queries completed")
	}

	return db
}

var DB = initDB()
