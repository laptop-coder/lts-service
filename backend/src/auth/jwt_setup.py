from typing import Literal
import datetime
import uuid

import jwt

from .rsa_keys import private_key


jwt_exp: dict[str, int] = {
    'access': 900,  # 900 seconds = 15 minutes
    'refresh': 2592000,  # 2592000 seconds = 30 days
}


def create_jwt(
    payload: dict[str, int | str],
    type: Literal['access'] | Literal['refresh'],
) -> str:
    payload['exp'] = int(datetime.datetime.now().timestamp()) + jwt_exp[type]
    payload['iat'] = int(datetime.datetime.now().timestamp())
    payload['jti'] = str(uuid.uuid4())
    return jwt.encode(payload, private_key, algorithm='RS256')
