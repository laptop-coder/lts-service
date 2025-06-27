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
