from typing import Literal
import base64

from PIL import Image

import consts


def write_photo_to_the_storage(
    type: Literal['lost'] | Literal['found'], id: int, photo_base64: str
):
    path_to_photo = f'{consts.PATH_TO_STORAGE}/{type}/{id}.jpeg'
    with open(path_to_photo, 'wb') as photo:
        photo.write(base64.b64decode(photo_base64))
    photo = Image.open(path_to_photo)
    photo.save(path_to_photo, quality=25)
