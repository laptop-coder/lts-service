import base64
import datetime
import os
from pathlib import Path
import sqlite3

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import Optional


PATH_TO_DB = os.getenv("PATH_TO_DB")
PATH_TO_STORAGE = os.getenv("PATH_TO_STORAGE")

app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins="http://localhost",
    allow_methods=["*"],
    allow_headers=["*"]
)


class LostThingData(BaseModel):
    thing_name: str
    email: str
    custom_text: str
    thing_photo: Optional[ str ] = None
    

class FoundThingData(BaseModel):
    thing_name: str
    thing_location: str
    custom_text: str
    thing_photo: Optional[ str ] = None


# Creating storage directories
try:
    Path(f"{PATH_TO_STORAGE}/lost").mkdir(parents=True)
    Path(f"{PATH_TO_STORAGE}/found").mkdir(parents=True)
except FileExistsError:
    pass


# Creating database tables
with sqlite3.connect(PATH_TO_DB) as connection:
    cursor = connection.cursor()
    cursor.execute("""
    CREATE TABLE IF NOT EXISTS lost_thing (
        id INTEGER PRIMARY KEY,
        publication_date TEXT NOT NULL,
        publication_time TEXT NOT NULL,
        thing_name TEXT NOT NULL,
        email varchar(254) NOT NULL,
        custom_text TEXT NOT NULL,
        status INTEGER NOT NULL
    );
    """)
    cursor.execute("""
    CREATE TABLE IF NOT EXISTS found_thing (
        id INTEGER PRIMARY KEY,
        publication_date TEXT NOT NULL,
        publication_time TEXT NOT NULL,
        thing_name TEXT NOT NULL,
        thing_location TEXT NOT NULL,
        custom_text TEXT NOT NULL,
        status INTEGER NOT NULL
    );
    """)


@app.get("/get_things_list")
def get_things_list(type: str):
    with sqlite3.connect(PATH_TO_DB) as connection:
        cursor = connection.cursor()
        data = cursor.execute(
            f"""
            SELECT * FROM {type}_thing WHERE status=0 ORDER BY id DESC;
            """
        ).fetchall()
        formatted_data = []
        if type == "lost":
            for i in range(len(data)):
                formatted_data.append({
                    "id": data[i][0],
                    "publication_date": data[i][1],
                    "publication_time": data[i][2],
                    "thing_name": data[i][3],
                    "email": data[i][4],
                    "custom_text": data[i][5],
                    "thing_photo": get_thing_photo("lost", data[i][0]),
                    "status": data[i][6]
                })
        elif type == "found":
            for i in range(len(data)):
                formatted_data.append({
                    "id": data[i][0],
                    "publication_date": data[i][1],
                    "publication_time": data[i][2],
                    "thing_name": data[i][3],
                    "thing_location": data[i][4],
                    "custom_text": data[i][5],
                    "thing_photo": get_thing_photo("found", data[i][0]),
                    "status": data[i][6]
                })
    return formatted_data
    

@app.get("/get_thing_photo")
def get_thing_photo(type: str, id: int):
    try:
        with open(f"{PATH_TO_STORAGE}/{type}/{id}.jpeg", "rb") as photo:
            photo_base64 = base64.b64encode(photo.read())
            return photo_base64
    except FileNotFoundError:
        return ""


@app.post("/add_new_lost_thing")
def add_new_lost_thing(data: LostThingData):
    with sqlite3.connect(PATH_TO_DB) as connection:
        cursor = connection.cursor()
        cursor.execute(
            f"""
            INSERT INTO lost_thing (
                publication_date,
                publication_time,
                thing_name,
                email,
                custom_text,
                status
            )
            VALUES (
                '{str(datetime.datetime.now())[0:10]}',
                '{str(datetime.datetime.now())[11:16]}',
                '{data.thing_name}',
                '{data.email}',
                '{data.custom_text}',
                0
            );
            """
        )
    try:
        thing_photo = f"{data.thing_photo[23:]}".encode()
        with open(f"{PATH_TO_STORAGE}/lost/{cursor.lastrowid}.jpeg", "wb") as file:
            file.write(base64.decodebytes(thing_photo))
    except TypeError:
        pass


@app.post("/add_new_found_thing")
def add_new_found_thing(data: FoundThingData):
    with sqlite3.connect(PATH_TO_DB) as connection:
        cursor = connection.cursor()
        cursor.execute(
            f"""
            INSERT INTO found_thing (
                publication_date,
                publication_time,
                thing_name,
                thing_location,
                custom_text,
                status
            )
            VALUES (
                '{str(datetime.datetime.now())[0:10]}',
                '{str(datetime.datetime.now())[11:16]}',
                '{data.thing_name}',
                '{data.thing_location}',
                '{data.custom_text}',
                0
            );
            """
        )
    try:
        thing_photo = f"{data.thing_photo[23:]}".encode()
        with open(f"{PATH_TO_STORAGE}/found/{cursor.lastrowid}.jpeg", "wb") as file:
            file.write(base64.decodebytes(thing_photo))
    except TypeError:
        pass


@app.get("/change_thing_status")
def change_thing_status(type: str, id: int):
    with sqlite3.connect(PATH_TO_DB) as connection:
        cursor = connection.cursor()
        cursor.execute(
            f"""
            UPDATE {type}_thing SET status=1 WHERE id={id};
            """
        )
        try:
            os.remove(f"{PATH_TO_STORAGE}/{type}/{id}.jpeg")
        except FileNotFoundError:
            pass

