import urllib.error
import json
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
        with urllib.request.urlopen(req) as response:
            try:
                # API returned the error
                if json.loads(response.read().decode("utf-8"))["error"]:
                    print_err("ERROR")
                    msg = "Failed to post on wall in VK"
                    print_err(msg)
                    raise Exception()
            except KeyError:
                # There are no errors
                pass
    except json.JSONDecodeError:
        print_err("ERROR")
        msg = "Failed to parse JSON response"
        print_err(msg)
        raise Exception()
    except urllib.error.HTTPError as e:
        print_err("ERROR")
        print_err(
            f"Failed to post on VK wall! Status code: {e.code}. Error: {e.reason}"
        )
        raise Exception
    except urllib.error.URLError as e:
        print_err("ERROR")
        print_err(f"Failed to post on VK wall! Error: {e.reason}")
        raise Exception


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
        with urllib.request.urlopen(req) as response:
            try:
                # API returned the error
                if json.loads(response.read().decode("utf-8"))["error"]:
                    print_err("ERROR")
                    msg = "Failed to send an alert in VK"
                    print_err(msg)
                    raise Exception()
            except KeyError:
                # There are no errors
                pass
    except json.JSONDecodeError:
        print_err("ERROR")
        msg = "Failed to parse JSON response"
        print_err(msg)
        raise Exception()
    except urllib.error.HTTPError as e:
        print_err("ERROR")
        print_err(
            f"Failed to send alert in VK! Status code: {e.code}. Error: {e.reason}"
        )
    except urllib.error.URLError as e:
        print_err("ERROR")
        print_err(f"Failed to send alert in VK! Error: {e.reason}")
