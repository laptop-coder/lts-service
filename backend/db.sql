CREATE TABLE lost_thing (
    id INTEGER PRIMARY KEY,
    publication_date TEXT NOT NULL,
    publication_time TEXT NOT NULL,
    thing_name TEXT NOT NULL,
    user_contacts TEXT NOT NULL,
    custom_text TEXT NOT NULL,
    status INTEGER NOT NULL
);

CREATE TABLE found_thing (
    id INTEGER PRIMARY KEY,
    publication_date TEXT NOT NULL,
    publication_time TEXT NOT NULL,
    thing_name TEXT NOT NULL,
    thing_location TEXT NOT NULL,
    custom_text TEXT NOT NULL,
    status INTEGER NOT NULL
);
