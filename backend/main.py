import base64
import os
from pathlib import Path
import sqlite3

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import Literal, Optional


PATH_TO_DB = os.getenv('PATH_TO_DB', '/backend/data/db.sqlite3')
PATH_TO_STORAGE = os.getenv('PATH_TO_STORAGE', '/backend/data/storage')
PORT = os.getenv('PORT', 80)

app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=f'http://localhost:{PORT}',
    allow_methods=['*'],
    allow_headers=['*'],
)


class LostThingData(BaseModel):
    thing_name: str
    email: str
    custom_text: str
    thing_photo: Optional[str] = None


class FoundThingData(BaseModel):
    thing_name: str
    thing_location: str
    custom_text: str
    thing_photo: Optional[str] = None


# Creating storage directories
Path(f'{PATH_TO_STORAGE}/lost').mkdir(parents=True, exist_ok=True)
Path(f'{PATH_TO_STORAGE}/found').mkdir(parents=True, exist_ok=True)


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
        verified INTEGER NOT NULL,
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
        verified INTEGER NOT NULL,
        status INTEGER NOT NULL
    );
    """)


def write_photo_to_the_storage(
    type: Literal['lost'] | Literal['found'], id: int, photo_base64: str
):
    with open(f'{PATH_TO_STORAGE}/{type}/{id}.jpeg', 'wb') as photo:
        photo.write(base64.b64decode(photo_base64))


@app.get('/get_things_list')
def get_things_list(type: Literal['lost'] | Literal['found']):
    with sqlite3.connect(PATH_TO_DB) as connection:
        cursor = connection.cursor()
        data = cursor.execute(
            f"""
            SELECT * FROM {type}_thing WHERE verified=1 AND status=0 ORDER BY id DESC;
            """
        ).fetchall()
        formatted_data = []
        match type:
            case 'lost':
                for elem in data:
                    formatted_data.append(
                        {
                            'id': elem[0],
                            'publication_date': elem[1],
                            'publication_time': elem[2],
                            'thing_name': elem[3],
                            'email': elem[4],
                            'custom_text': elem[5],
                            'thing_photo': get_thing_photo('lost', elem[0]),
                            'verified': elem[6],
                            'status': elem[7],
                        }
                    )
            case 'found':
                for elem in data:
                    formatted_data.append(
                        {
                            'id': elem[0],
                            'publication_date': elem[1],
                            'publication_time': elem[2],
                            'thing_name': elem[3],
                            'thing_location': elem[4],
                            'custom_text': elem[5],
                            'thing_photo': get_thing_photo('found', elem[0]),
                            'verified': elem[6],
                            'status': elem[7],
                        }
                    )
    return formatted_data


@app.get('/get_thing_data')
def get_thing_data(type: Literal['lost'] | Literal['found'], id: int):
    with sqlite3.connect(PATH_TO_DB) as connection:
        cursor = connection.cursor()
        data = cursor.execute(
            f"""
            SELECT * FROM {type}_thing WHERE id={id};
            """
        ).fetchone()
        formatted_data = {}
        match type:
            case 'lost':
                formatted_data = {
                    'id': data[0],
                    'publication_date': data[1],
                    'publication_time': data[2],
                    'thing_name': data[3],
                    'email': data[4],
                    'custom_text': data[5],
                    'thing_photo': get_thing_photo('lost', data[0]),
                    'verified': data[6],
                    'status': data[7],
                }
            case 'found':
                formatted_data = {
                    'id': data[0],
                    'publication_date': data[1],
                    'publication_time': data[2],
                    'thing_name': data[3],
                    'thing_location': data[4],
                    'custom_text': data[5],
                    'thing_photo': get_thing_photo('found', data[0]),
                    'verified': data[6],
                    'status': data[7],
                }
    return formatted_data


@app.get('/get_thing_photo')
def get_thing_photo(type: Literal['lost'] | Literal['found'], id: int):
    path_to_photo = f'{PATH_TO_STORAGE}/{type}/{id}.jpeg'
    if os.path.exists(path_to_photo):
        with open(path_to_photo, 'rb') as photo:
            return base64.b64encode(photo.read())
    return ''


@app.post('/add_new_lost_thing')
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
                verified,
                status
            )
            VALUES (
                date('now'),
                substr(time('now'), 1, 5),
                '{data.thing_name}',
                '{data.email}',
                '{data.custom_text}',
                0,
                0
            );
            """
        )
        if data.thing_photo is not None and cursor.lastrowid is not None:
            write_photo_to_the_storage(
                'lost',
                cursor.lastrowid,
                data.thing_photo[23:],
            )


@app.post('/add_new_found_thing')
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
                verified,
                status
            )
            VALUES (
                date('now'),
                substr(time('now'), 1, 5),
                '{data.thing_name}',
                '{data.thing_location}',
                '{data.custom_text}',
                0,
                0
            );
            """
        )
        if data.thing_photo is not None and cursor.lastrowid is not None:
            write_photo_to_the_storage(
                'found',
                cursor.lastrowid,
                data.thing_photo[23:],
            )


@app.get('/change_thing_status')
def change_thing_status(type: Literal['lost'] | Literal['found'], id: int):
    with sqlite3.connect(PATH_TO_DB) as connection:
        cursor = connection.cursor()
        cursor.execute(
            f"""
            UPDATE {type}_thing SET status=1 WHERE id={id};
            """
        )
