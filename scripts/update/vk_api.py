import urllib.error
import random
import urllib.parse
import urllib.request
from .config import VK_ALERTS_GROUP_ID, VK_ALERTS_API_KEY, VK_ALERTS_CHAT_ID
from .utils import print_err


def post_on_wall(message: str) -> None:
    data = urllib.parse.urlencode(
        {
            "owner_id": VK_ALERTS_GROUP_ID,
            "message": message,
            "from_group": "1",
            "signed": "0",
            "v": "5.199",
        },
    ).encode("utf-8")
    req = urllib.request.Request(
        "https://api.vk.ru/method/wall.post",
        headers={"Authorization": f"Bearer {VK_ALERTS_API_KEY}"},
        data=data,
    )
    try:
        response = urllib.request.urlopen(req)
        response.close()
    except urllib.error.HTTPError as e:
        print_err("ERROR")
        print_err(
            f"Failed to post on VK wall! Status code: {e.code}. Error: {e.reason}"
        )
    except urllib.error.URLError as e:
        print_err("ERROR")
        print_err(f"Failed to post on VK wall! Error: {e.reason}")


def send_alert(message: str) -> None:
    data = urllib.parse.urlencode(
        {
            "peer_id": VK_ALERTS_CHAT_ID,
            "random_id": random.randint(0, 2_147_483_647),
            "message": message,
            "v": "5.199",
        },
    ).encode("utf-8")
    req = urllib.request.Request(
        "https://api.vk.ru/method/messages.send",
        headers={"Authorization": f"Bearer {VK_ALERTS_API_KEY}"},
        data=data,
    )
    try:
        response = urllib.request.urlopen(req)
        response.close()
    except urllib.error.HTTPError as e:
        print_err("ERROR")
        print_err(
            f"Failed to send alert to VK! Status code: {e.code}. Error: {e.reason}"
        )
    except urllib.error.URLError as e:
        print_err("ERROR")
        print_err(f"Failed to send alert to VK! Error: {e.reason}")
