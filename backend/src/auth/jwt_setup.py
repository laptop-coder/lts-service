from typing import Literal
import datetime

from cryptography.hazmat.primitives.asymmetric.rsa import RSAPrivateKey
import jwt


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
