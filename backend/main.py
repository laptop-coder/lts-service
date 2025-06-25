import base64
import os
import sys
from pathlib import Path
import sqlite3
from argon2 import PasswordHasher
from cryptography.hazmat.primitives.asymmetric.rsa import generate_private_key
from cryptography.hazmat.primitives import serialization
import jwt

from fastapi import FastAPI, Response
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import Literal, Optional
from PIL import Image


PATH_TO_DB = os.getenv('PATH_TO_DB', '/backend/data/db/db.sqlite3')
PATH_TO_ENV = os.getenv('PATH_TO_ENV', '/env')
PATH_TO_PRIVATE_KEY = f'{PATH_TO_ENV}/rsa_key'
PATH_TO_PUBLIC_KEY = f'{PATH_TO_ENV}/rsa_key.pub'
PATH_TO_STORAGE = os.getenv('PATH_TO_STORAGE', '/backend/data/storage')
PORT = os.getenv('PORT', 80)
PRIVATE_KEY_ENCRYPTION_PASSWORD = os.getenv(
    'PRIVATE_KEY_ENCRYPTION_PASSWORD', ''
)


# Create RSA keys if not exist
if not os.path.isfile(PATH_TO_PRIVATE_KEY) and not os.path.isfile(
    PATH_TO_PUBLIC_KEY
):
    private_key = generate_private_key(public_exponent=65537, key_size=4096)
    public_key = private_key.public_key()

    private_key_serialized = private_key.private_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PrivateFormat.PKCS8,
        encryption_algorithm=serialization.BestAvailableEncryption(
            PRIVATE_KEY_ENCRYPTION_PASSWORD.encode()
        ),
    ).decode()
    public_key_serialized = public_key.public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo,
    ).decode()

    with open(PATH_TO_PRIVATE_KEY, 'w') as file:
        file.write(private_key_serialized)
    with open(PATH_TO_PUBLIC_KEY, 'w') as file:
        file.write(public_key_serialized)


# Read keys
with open(PATH_TO_PRIVATE_KEY, 'rb') as file:
    private_key = serialization.load_pem_private_key(
        file.read(),
        password=PRIVATE_KEY_ENCRYPTION_PASSWORD.encode(),
    )
with open(PATH_TO_PUBLIC_KEY, 'rb') as file:
    public_key = serialization.load_pem_public_key(file.read())


app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=f'http://localhost:{PORT}',
    allow_methods=['*'],
    allow_headers=['*'],
)


ph = PasswordHasher()


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


class ModeratorRegister(BaseModel):
    username: str
    password: str


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
    cursor.execute("""
    CREATE TABLE IF NOT EXISTS moderator (
        id INTEGER PRIMARY KEY,
        username varchar(32) NOT NULL,
        password TEXT NOT NULL
    );
    """)


def write_photo_to_the_storage(
    type: Literal['lost'] | Literal['found'], id: int, photo_base64: str
):
    path_to_photo = f'{PATH_TO_STORAGE}/{type}/{id}.jpeg'
    with open(path_to_photo, 'wb') as photo:
        photo.write(base64.b64decode(photo_base64))
    photo = Image.open(path_to_photo)
    photo.save(path_to_photo, quality=25)


@app.get('/get_things_list')
def get_things_list(type: Literal['lost'] | Literal['found']):
    with sqlite3.connect(PATH_TO_DB) as connection:
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


def check_moderator_exists(username: str):
    with sqlite3.connect(PATH_TO_DB) as connection:
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
            sys.exit('Error! Multiple moderators have the same username')


@app.post('/moderator/register')
def moderator_register(response: Response, data: ModeratorRegister):
    password_hash = ph.hash(data.password)
    with sqlite3.connect(PATH_TO_DB) as connection:
        cursor = connection.cursor()
        cursor.execute(
            f"""
            INSERT INTO moderator (username, password) VALUES (
            '{data.username}',
            '{password_hash}'
            );
            """
        )
    jwt_payload = {
        'username': data.username,
        'password': password_hash,
    }
    jwt_encoded = jwt.encode(
        jwt_payload,
        private_key.private_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PrivateFormat.PKCS8,
            encryption_algorithm=serialization.NoEncryption(),
        ),
        algorithm='RS256',
    )
    response.set_cookie(key='jwt', value=jwt_encoded)
