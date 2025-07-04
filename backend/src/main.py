from pathlib import Path
from typing import Literal, Optional
import base64
import sqlite3

from PIL import Image
from argon2 import PasswordHasher
from argon2.exceptions import VerifyMismatchError
from fastapi import FastAPI, Response
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel

from auth.jwt_setup import create_jwt
from exceptions import MultipleModeratorsHaveTheSameUsername
import consts


ph = PasswordHasher()
app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=f'http://localhost:{consts.PORT}',
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


class ModeratorAuth(BaseModel):
    username: str
    password: str


# Creating storage directories
Path(f'{consts.PATH_TO_STORAGE}/lost').mkdir(parents=True, exist_ok=True)
Path(f'{consts.PATH_TO_STORAGE}/found').mkdir(parents=True, exist_ok=True)


# Initial SQL-requests (creating database tables)
with (
    sqlite3.connect(consts.PATH_TO_DB) as connection,
    open('./db.sql', 'r') as sql_requests,
):
    connection.cursor().executescript(sql_requests.read())


def write_photo_to_the_storage(
    type: Literal['lost'] | Literal['found'], id: int, photo_base64: str
):
    path_to_photo = f'{consts.PATH_TO_STORAGE}/{type}/{id}.jpeg'
    with open(path_to_photo, 'wb') as photo:
        photo.write(base64.b64decode(photo_base64))
    photo = Image.open(path_to_photo)
    photo.save(path_to_photo, quality=25)


@app.get('/get_things_list')
def get_things_list(type: Literal['lost'] | Literal['found']):
    with sqlite3.connect(consts.PATH_TO_DB) as connection:
        cursor = connection.cursor()
        data = cursor.execute(
            f"""
            SELECT * FROM {type}_thing WHERE status=0 ORDER BY id DESC;
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
                            'verified': elem[6],
                            'status': elem[7],
                        }
                    )
    return formatted_data


@app.get('/get_thing_data')
def get_thing_data(type: Literal['lost'] | Literal['found'], id: int):
    with sqlite3.connect(consts.PATH_TO_DB) as connection:
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
                    'verified': data[6],
                    'status': data[7],
                }
    return formatted_data


@app.post('/add_new_lost_thing')
def add_new_lost_thing(data: LostThingData):
    with sqlite3.connect(consts.PATH_TO_DB) as connection:
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
    with sqlite3.connect(consts.PATH_TO_DB) as connection:
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
    with sqlite3.connect(consts.PATH_TO_DB) as connection:
        cursor = connection.cursor()
        cursor.execute(
            f"""
            UPDATE {type}_thing SET status=1 WHERE id={id};
            """
        )


def check_moderator_exists(username: str):
    """
    Return 1 if moderator with this username exists.
    Return 0 if moderator with this username doesn't exist.
    Raise error if more than one moderator with this username exist.
    """
    with sqlite3.connect(consts.PATH_TO_DB) as connection:
        cursor = connection.cursor()
        [count] = cursor.execute(
            f"""
            SELECT COUNT(*) FROM moderator WHERE username='{username}';
            """
        ).fetchone()
        if count == 1:
            return True
        elif count == 0:
            return False
        else:
            raise MultipleModeratorsHaveTheSameUsername


@app.post('/moderator/register')
def moderator_register(response: Response, data: ModeratorAuth):
    try:
        if not check_moderator_exists(data.username):
            password_hash = ph.hash(data.password)
            with sqlite3.connect(consts.PATH_TO_DB) as connection:
                cursor = connection.cursor()
                cursor.execute(
                    f"""
                    INSERT INTO moderator (username, password) VALUES (
                    '{data.username}',
                    '{password_hash}'
                    );
                    """
                )
            jwt_payload: dict[str, int | str] = {
                'username': data.username,
                'password': password_hash,
            }
            response.set_cookie(
                key='jwt',
                value=create_jwt(jwt_payload, 'access'),
                httponly=True,
            )
        else:
            return {
                'Message': 'moderator with this username already exists, use a different username'
            }
    except MultipleModeratorsHaveTheSameUsername:
        return {'Error': 'multiple moderators have the same username'}


@app.post('/moderator/login')
def moderator_login(response: Response, data: ModeratorAuth):
    try:
        if check_moderator_exists(data.username):
            with sqlite3.connect(consts.PATH_TO_DB) as connection:
                cursor = connection.cursor()
                [password_hash] = cursor.execute(
                    f"""
                    SELECT password FROM moderator WHERE username='{data.username}' LIMIT 1;
                    """
                ).fetchone()
                try:
                    ph.verify(password_hash, data.password)
                    if ph.check_needs_rehash(password_hash):
                        # See https://argon2-cffi.readthedocs.io/en/stable/api.html#argon2.PasswordHasher.check_needs_rehash
                        password_hash = ph.hash(data.password)
                        cursor.execute(
                            f"""
                            UPDATE moderator SET password='{password_hash}' WHERE username='{data.username}';
                            """
                        )
                    jwt_payload: dict[str, int | str] = {
                        'username': data.username,
                        'password': password_hash,
                    }
                    response.set_cookie(
                        key='jwt',
                        value=create_jwt(jwt_payload, 'access'),
                        httponly=True,
                    )
                except VerifyMismatchError:
                    return {'Message': 'passwords do not match'}
        else:
            return {
                'Message': 'moderator with this username does not exist, please register first'
            }
    except MultipleModeratorsHaveTheSameUsername:
        return {'Error': 'multiple moderators have the same username'}
