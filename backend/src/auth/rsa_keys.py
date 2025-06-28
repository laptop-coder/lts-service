import os

from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric.rsa import generate_private_key

from .. import consts


# Create RSA keys if not exist
if not os.path.isfile(consts.PATH_TO_PRIVATE_KEY) and not os.path.isfile(
    consts.PATH_TO_PUBLIC_KEY
):
    private_key = generate_private_key(public_exponent=65537, key_size=4096)
    public_key = private_key.public_key()

    private_key_serialized = private_key.private_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PrivateFormat.PKCS8,
        encryption_algorithm=serialization.BestAvailableEncryption(
            consts.PRIVATE_KEY_ENCRYPTION_PASSWORD.encode()
        ),
    ).decode()
    public_key_serialized = public_key.public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo,
    ).decode()

    with open(consts.PATH_TO_PRIVATE_KEY, 'w') as file:
        file.write(private_key_serialized)
    with open(consts.PATH_TO_PUBLIC_KEY, 'w') as file:
        file.write(public_key_serialized)


# Read keys
with open(consts.PATH_TO_PRIVATE_KEY, 'rb') as file:
    private_key = (
        serialization.load_pem_private_key(
            file.read(),
            password=consts.PRIVATE_KEY_ENCRYPTION_PASSWORD.encode(),
        )
        .private_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PrivateFormat.PKCS8,
            encryption_algorithm=serialization.NoEncryption(),
        )
        .decode()
    )

with open(consts.PATH_TO_PUBLIC_KEY, 'rb') as file:
    public_key = (
        serialization.load_pem_public_key(file.read())
        .public_bytes(
            encoding=serialization.Encoding.PEM,
            format=serialization.PublicFormat.SubjectPublicKeyInfo,
        )
        .decode()
    )
