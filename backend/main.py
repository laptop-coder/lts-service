from PIL import Image
from argon2 import PasswordHasher
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric.rsa import RSAPrivateKey
from cryptography.hazmat.primitives.asymmetric.rsa import generate_private_key
from exceptions import MultipleModeratorsHaveTheSameUsername
from fastapi import FastAPI, Response
from fastapi.middleware.cors import CORSMiddleware
from pathlib import Path
from pydantic import BaseModel
from typing import Literal, Optional
import base64
import datetime
import jwt
import os
import sqlite3


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
    private_key = (
        serialization.load_pem_private_key(
            file.read(),
            password=PRIVATE_KEY_ENCRYPTION_PASSWORD.encode(),
        )
        .private_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PrivateFormat.PKCS8,
            encryption_algorithm=serialization.NoEncryption(),
        )
        .decode()
    )

with open(PATH_TO_PUBLIC_KEY, 'rb') as file:
    public_key = (
        serialization.load_pem_public_key(file.read())
        .public_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PublicFormat.SubjectPublicKeyInfo,
        )
        .decode()
    )


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
    """
    Return 1 if moderator with this username exists.
    Return 0 if moderator with this username doesn't exist.
    Raise error if more than one moderator with this username exist.
    """
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
            raise MultipleModeratorsHaveTheSameUsername


jwt_exp: dict[str, int] = {
    'access': 900,  # 900 seconds = 15 minutes
    'refresh': 2592000,  # 2592000 seconds = 30 days
}


def create_jwt(
    private_key: RSAPrivateKey | str,
    payload: dict[str, int | str],
    type: Literal['access'] | Literal['refresh'],
) -> str:
    payload['exp'] = int(datetime.datetime.now().timestamp()) + jwt_exp[type]
    return jwt.encode(payload, private_key, algorithm='RS256')


@app.post('/moderator/register')
def moderator_register(response: Response, data: ModeratorRegister):
    try:
        if not check_moderator_exists(data.username):
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
            jwt_payload: dict[str, int | str] = {
                'username': data.username,
                'password': password_hash,
            }
            response.set_cookie(
                key='jwt',
                value=create_jwt(private_key, jwt_payload, 'access'),
                httponly=True,
            )
        else:
            return {
                'Message': 'moderator with this username already exists, use a different username'
            }
    except MultipleModeratorsHaveTheSameUsername:
        return {'Error': 'multiple moderators have the same username'}
