import datetime
import os
import sqlite3
import base64

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import Optional

import config


class LostThingData(BaseModel):
    thing_name: str
    user_contacts: str
    custom_text: str
    thing_photo: Optional[ str ] = None
    

class FoundThingData(BaseModel):
    thing_name: str
    thing_location: str
    custom_text: str
    thing_photo: str


app = FastAPI()

origins = [
    "http://localhost:4000"
]
app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_methods=["*"],
    allow_headers=["*"]
)


@app.get("/get_things_list")
def get_things_list(type: str):
    connection = sqlite3.connect(config.PATH_TO_DB)
    with connection:
        cursor = connection.cursor()
        data = cursor.execute(
            f"""
            SELECT * FROM {type}_thing WHERE status=0;
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
                    "user_contacts": data[i][4],
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
        with open(f"{config.PATH_TO_STORAGE}/{type}/{id}.jpeg", "rb") as photo:
            photo_base64 = base64.b64encode(photo.read())
            return photo_base64
    except:
        return ""


@app.post("/add_new_lost_thing")
def add_new_lost_thing(data: LostThingData):
    connection = sqlite3.connect(config.PATH_TO_DB)
    with connection:
        cursor = connection.cursor()
        cursor.execute(
            f"""
            INSERT INTO lost_thing (
                publication_date,
                publication_time,
                thing_name,
                user_contacts,
                custom_text,
                status
            )
            VALUES (
                '{str(datetime.datetime.now())[0:10]}',
                '{str(datetime.datetime.now())[11:16]}',
                '{data.thing_name}',
                '{data.user_contacts}',
                '{data.custom_text}',
                0
            );
            """
        )
    try:
        thing_photo = f"{data.thing_photo[23:]}".encode()
        with open(f"./storage/lost/{cursor.lastrowid}.jpeg", "wb") as file:
            file.write(base64.decodebytes(thing_photo))
    except TypeError:
        pass


@app.post("/add_new_found_thing")
def add_new_found_thing(data: FoundThingData):
    connection = sqlite3.connect(config.PATH_TO_DB)
    with connection:
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
    thing_photo = f"{data.thing_photo[23:]}".encode()
    with open(f"./storage/found/{cursor.lastrowid}.jpeg", "wb") as file:
        file.write(base64.decodebytes(thing_photo))


@app.get("/change_thing_status")
def change_thing_status(type: str, id: int):
    connection = sqlite3.connect(config.PATH_TO_DB)
    with connection:
        cursor = connection.cursor()
        cursor.execute(
            f"""
            UPDATE {type}_thing SET status=1 WHERE id={id};
            """
        )
        try:
            os.remove(f"./storage/{type}/{id}.jpeg")
        except FileNotFoundError:
            pass

