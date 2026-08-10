from enum import Enum
import json
import os
import random
import signal
import subprocess
import sys
import time
import urllib.error
import urllib.parse
import urllib.request


path_to_project = f"{os.getenv('HOME')}/lost-things-search"


def load_env(file_path=f"{path_to_project}/.env"):
    if not os.path.exists(file_path):
        return
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
                os.environ[key] = value


load_env()


VK_ALERTS_API_KEY = os.environ.get("VK_ALERTS_API_KEY")
VK_ALERTS_CHAT_ID = os.environ.get("VK_ALERTS_CHAT_ID")


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


def print_wait(s: str) -> None:
    print(s, end=" ", flush=True)


def print_ok(s: str, line_break=True) -> None:
    end = "\n" if line_break else " "
    print(f"\033[32m{s}\033[0m", end=end, flush=True)


def print_err(s: str, line_break=True) -> None:
    end = "\n" if line_break else " "
    print(f"\033[31m{s}\033[0m", end=end, flush=True)


def print_secondary(s: str, line_break=True) -> None:
    end = "\n" if line_break else " "
    print(f"\033[2m{s}\033[0m", end=end, flush=True)


def run_command(command: list[str]) -> subprocess.CompletedProcess[bytes]:
    """
    Run shell command from the project dir.
    """
    return subprocess.run(
        command,
        capture_output=True,
        cwd=path_to_project,
    )


class Service(Enum):
    BACKEND = 0
    FRONTEND = 1
    ML = 2
    MIGRATE = 3


# Spinner
b = ["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"]


def main() -> bool:
    print("\033c", end="")
    print_wait(f"Got the project dir.")
    print_secondary(path_to_project)

    tag_file = f"{path_to_project}/.tag"
    tag_file_is_new = False

    if not os.path.exists(tag_file):
        print_wait("Tag file not found")
        print_secondary(tag_file)
        tag_file_is_new = True

    service = Service.BACKEND
    print(f"Looking at the {service.name} service.")
    repo = f"laptop-coder/lost-things-search-{service.name}".lower()

    # Get GHCR token
    print_wait("Trying to get GHCR token...")
    try:
        with urllib.request.urlopen(
            f"https://ghcr.io/token?scope=repository:{repo}:pull"
        ) as response:
            token = json.loads(response.read().decode("utf-8"))["token"]
    except json.JSONDecodeError:
        print_err("ERROR")
        msg = "Failed to parse JSON response"
        print_err(msg)
        send_alert(msg)
        return False
    except urllib.error.HTTPError as e:
        print_err("ERROR")
        msg = f"Failed to get GHCR token! Status code: {e.code}. Error: {e.reason}"
        print_err(msg)
        send_alert(msg)
        return False
    except urllib.error.URLError as e:
        print_err("ERROR")
        msg = f"Failed to get GHCR token! Error: {e.reason}"
        print_err(msg)
        send_alert(msg)
        return False
    print_ok("OK")

    # Get the latest tag
    print_wait(f"Trying to get the latest {service.name} service tag...")
    try:
        req = urllib.request.Request(
            f"https://ghcr.io/v2/{repo}/tags/list",
            headers={"Authorization": f"Bearer {token}"},
        )
        with urllib.request.urlopen(req) as response:
            latest_tag = json.loads(response.read().decode("utf-8"))["tags"][-1]
    except json.JSONDecodeError:
        print_err("ERROR")
        msg = "Failed to parse JSON response"
        print_err(msg)
        send_alert(msg)
        return False
    except urllib.error.HTTPError as e:
        print_err("ERROR")
        msg = f"Failed to get the latest {service.name} service tag! Status code: {e.code}. Error: {e.reason}"
        print_err(msg)
        send_alert(msg)
        return False
    except urllib.error.URLError as e:
        print_err("ERROR")
        msg = f"Failed to get the latest {service.name} service tag! Error: {e.reason}"
        print_err(msg)
        send_alert(msg)
        return False
    print_ok("OK", line_break=False)
    print_secondary(latest_tag)

    if tag_file_is_new:
        print_wait("Adding latest tag to the file...")
        with open(tag_file, "w") as file:
            file.write(f"{latest_tag}\n")
        print_ok("OK")
    else:
        print_wait(f"Getting current {service.name} tag...")
        with open(tag_file, "r") as file:
            current_tag = file.readline()[:-1]
        print_ok("OK", line_break=False)
        print_secondary(current_tag)
        print_wait("Compairing tags...")
        if current_tag != latest_tag:
            print_secondary(f"{current_tag} != {latest_tag}")
            print("Downloading new images:")

            processes = []

            print("╭" + "─" * (max([len(s.name) for s in Service]) + 4) + "╮")

            # Run parallel processes (download docker images)
            for s in Service:
                print(
                    "    "
                    + s.name
                    + " " * (max([len(x.name) for x in Service]) - len(s.name) + 1)
                    + "│"
                )
                repo = f"laptop-coder/lost-things-search-{s.name}".lower()
                processes.append(
                    subprocess.Popen(
                        ["docker", "pull", "-q", f"ghcr.io/{repo}:{latest_tag}"],
                        stdout=subprocess.DEVNULL,
                        stderr=subprocess.DEVNULL,
                    )
                )

            print("╰" + "─" * (max([len(x.name) for x in Service]) + 4) + "╯")

            # Print downloading status
            i = 0
            while True:
                return_codes = [p.poll() for p in processes]
                for j in range(len(return_codes)):
                    if return_codes[j] is None:
                        # Spinner
                        print(
                            "\r"
                            + "\033[1A" * (len(Service) - j + 1)
                            + f"│ \033[33m{b[i]}\033[0m"
                            + "\n" * (len(Service) - j + 1),
                            end="",
                            flush=True,
                        )
                    elif return_codes[j] == 0:
                        # OK status
                        print(
                            "\r"
                            + "\033[1A" * (len(Service) - j + 1)
                            + f"│ \033[32m✔\033[0m"
                            + "\n" * (len(Service) - j + 1),
                            end="",
                            flush=True,
                        )
                    else:
                        # ERROR status
                        print(
                            "\r"
                            + "\033[1A" * (len(Service) - j + 1)
                            + f"│ \033[31m✗\033[0m"
                            + "\n" * (len(Service) - j + 1),
                            end="",
                            flush=True,
                        )
                if return_codes.count(None) == 0:
                    break
                i += 1
                if i > len(b) - 1:
                    i = 0
                time.sleep(0.1)

            # Stop the project
            print_wait("Stopping the project...")
            result = run_command(["make", "down"])
            if result.returncode != 0:
                err = result.stderr.decode("utf-8")
                print_err("ERROR")
                msg = f"Failed to stop the project! Error: {err}"
                print_err(msg)
                send_alert(msg)
                return False
            print_ok("OK")

            # Pull the code changes
            print_wait("Pulling the code changes...")
            result = run_command(["git", "pull"])
            if result.returncode != 0:
                err = result.stderr.decode("utf-8")
                print_err("ERROR")
                msg = f"Failed to pull the code changes! Error: {err}"
                print_err(msg)
                send_alert(msg)
                return False
            print_ok("OK")

            # Deploy the project
            print_wait("Deploying the project...")
            result = run_command(["make", "deploy"])
            if result.returncode != 0:
                err = result.stderr.decode("utf-8")
                print_err("ERROR")
                msg = f"Failed to deploy the project! Error: {err}"
                print_err(msg)
                send_alert(msg)
                return False
            print_ok("OK")

            # Run migrations
            print_wait("Running migrations...")
            result = run_command(["make", "migrate"])
            if result.returncode != 0:
                err = result.stderr.decode("utf-8")
                print_err("ERROR")
                msg = f"Failed to run migrations! Error: {err}"
                print_err(msg)
                send_alert(msg)
                return False
            print_ok("OK")

            # Update current tag in file
            print_wait("Updating current tag in the file...")
            with open(tag_file, "w") as file:
                file.write(f"{latest_tag}\n")
            print_ok("OK")

        else:
            print_secondary(f"{current_tag} == {latest_tag}")
    return True


def graceful_shutdown(signum, frame):
    del signum, frame  # ignore parameters
    print("\n───────────────────────────────────")
    print(f"Received a signal, shutting down...")
    sys.exit(0)


if __name__ == "__main__":
    signal.signal(signal.SIGINT, graceful_shutdown)
    signal.signal(signal.SIGTERM, graceful_shutdown)
    for i in range(10):
        if main():
            print_ok("Done!")
            break
        msg = (
            f"{i + 1}/10 attempt. Waiting for 10 seconds to run script one more time..."
        )
        send_alert(msg)
        print(
            msg,
            flush=True,
        )
        time.sleep(10)
