import os
from dataclasses import dataclass
from enum import Enum
import gc

PATH_TO_PROJECT = f"{os.getenv('HOME')}/lost-things-search"


# Load env variables
def load_env(file_path: str) -> dict[str, str]:
    env = {}
    if not os.path.exists(file_path):
        return env
    with open(file_path, "r") as file:
        for line in file.readlines():
            line = line[:-1]
            # skip empty lines and comments
            if line and line[0] != "#":
                line = line.split("=", 1)
                key = line[0]
                line = line[1]
                # remove comments
                q = ""
                if line[0] == '"':
                    q = '"'
                elif line[0] == "'":
                    q = "'"
                if not q:  # without quotes
                    if line.count("#") > 0:
                        value = line[: line.find("#")].strip()
                    else:
                        value = line
                else:
                    value = line[1 : line[1:].find(q) + 1]
                env[key] = value
    return env


env = load_env(os.path.join(PATH_TO_PROJECT, ".env"))

VK_ALERTS_API_KEY = env["VK_ALERTS_API_KEY"]
VK_ALERTS_CHAT_ID = env["VK_ALERTS_CHAT_ID"]
VK_ALERTS_GROUP_ID = env["VK_ALERTS_GROUP_ID"]
GITHUB_PAT = env["GITHUB_PAT"]

del env
gc.collect()


class Service(Enum):
    BACKEND = 0
    FRONTEND = 1
    ML = 2
    MIGRATE = 3


# Spinner
b = ["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"]


MAIN_SERVICE = Service.BACKEND  # service used to determine the latest tag
TAG_FILE = f"{PATH_TO_PROJECT}/.tag"


@dataclass(slots=True, frozen=True)
class DigestDTO:
    old_tag: str | None
    latest_tag: str
    downloading_time: list[int] | None
    script_time: int
    changelog: str | None


class DigestDTOBuilder:
    def __init__(self):
        self._old_tag = None
        self._latest_tag = None
        self._downloading_time = None
        self._script_time = None
        self._changelog = None

    @property
    def old_tag(self):
        return self._old_tag

    @property
    def latest_tag(self):
        return self._latest_tag

    @property
    def downloading_time(self):
        return self._downloading_time

    @property
    def script_time(self):
        return self._script_time

    @property
    def changelog(self):
        return self._changelog

    @old_tag.setter
    def old_tag(self, value):
        self._old_tag = value

    @latest_tag.setter
    def latest_tag(self, value):
        self._latest_tag = value

    @downloading_time.setter
    def downloading_time(self, value):
        self._downloading_time = value

    @script_time.setter
    def script_time(self, value):
        self._script_time = value

    @changelog.setter
    def changelog(self, value):
        self._changelog = value

    def build(self) -> DigestDTO:

        if self._latest_tag is None:
            raise ValueError("Latest tag cannot be None")
        if self._script_time is None:
            raise ValueError("Script time cannot be None")

        return DigestDTO(
            old_tag=self._old_tag,
            latest_tag=self._latest_tag,
            downloading_time=self._downloading_time,
            script_time=self._script_time,
            changelog=self._changelog,
        )
