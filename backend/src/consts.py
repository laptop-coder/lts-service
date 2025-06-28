import os

PATH_TO_DB = os.getenv('PATH_TO_DB', '/backend/data/db/db.sqlite3')
PATH_TO_ENV = os.getenv('PATH_TO_ENV', '/env')
PATH_TO_PRIVATE_KEY = f'{PATH_TO_ENV}/private_key.pem'
PATH_TO_PUBLIC_KEY = f'{PATH_TO_ENV}/public_key.pem'
PATH_TO_STORAGE = os.getenv('PATH_TO_STORAGE', '/backend/data/storage')
PORT = os.getenv('PORT', 80)
PRIVATE_KEY_ENCRYPTION_PASSWORD = os.getenv(
    'PRIVATE_KEY_ENCRYPTION_PASSWORD', ''
)
