import datetime
import time
import subprocess
from .vk_api import send_alert, post_on_wall
from .config import Service, MAIN_SERVICE, b, DigestDTO, GITHUB_PAT
import urllib.error
import urllib.request
import json
from .utils import print_wait, print_ok, print_err, print_secondary, run_command


def publish_update_digest(dto: DigestDTO) -> None:
    print_wait("Publishing digest...")

    current_date = (
        datetime.datetime.now()
        .astimezone(datetime.timezone(datetime.timedelta(hours=3)))
        .strftime("%d.%m.%Y")
    )
    current_time = datetime.datetime.now(
        datetime.timezone(datetime.timedelta(hours=3))
    ).strftime("%H:%M")

    content = "LostThingsSearch\n"
    content += "\n"
    content += f"Обновление завершено {current_date} в {current_time}\n"
    content += "\n"

    if dto.old_tag:
        content += f"Старая версия: {dto.old_tag}\n"

    content += f"Новая версия: {dto.latest_tag}\n"

    content += "\n"

    if dto.downloading_time:
        content += "Время загрузки Docker-образов:\n"
        for i in range(len(Service)):
            content += (
                f"• {list(Service)[i].name.lower()}: {dto.downloading_time[i]} с\n"
            )

    content += "\n"

    content += f"Время работы скрипта: {dto.script_time} с\n"

    if dto.changelog:
        content += dto.changelog

    try:
        post_on_wall(content)
    except:
        print_err("ERROR")
        msg = "Failed to publish digest!"
        print_err(msg)
        send_alert(msg)
        raise Exception()
    print_ok("OK")


def get_ghcr_token() -> str:
    print_wait("Getting GHCR token...")
    try:
        with urllib.request.urlopen(
            f"https://ghcr.io/token?scope=repository:laptop-coder/lost-things-search-{MAIN_SERVICE.name}:pull".lower()
        ) as response:
            token = json.loads(response.read().decode("utf-8"))["token"]
    except json.JSONDecodeError:
        print_err("ERROR")
        msg = "Failed to parse JSON response"
        print_err(msg)
        send_alert(msg)
        raise Exception()
    except urllib.error.HTTPError as e:
        print_err("ERROR")
        msg = f"Failed to get GHCR token! Status code: {e.code}. Error: {e.reason}"
        print_err(msg)
        send_alert(msg)
        raise Exception()
    except urllib.error.URLError as e:
        print_err("ERROR")
        msg = f"Failed to get GHCR token! Error: {e.reason}"
        print_err(msg)
        send_alert(msg)
        raise Exception()
    print_ok("OK")
    return token


def get_latest_tag(token: str) -> str:
    print_wait(f"Getting the latest {MAIN_SERVICE.name} service tag...")
    try:
        req = urllib.request.Request(
            f"https://ghcr.io/v2/laptop-coder/lost-things-search-{MAIN_SERVICE.name}/tags/list".lower(),
            headers={"Authorization": f"Bearer {token}"},
        )
        with urllib.request.urlopen(req) as response:
            latest_tag = json.loads(response.read().decode("utf-8"))["tags"][-1]
    except json.JSONDecodeError:
        print_err("ERROR")
        msg = "Failed to parse JSON response"
        print_err(msg)
        send_alert(msg)
        raise Exception()
    except urllib.error.HTTPError as e:
        print_err("ERROR")
        msg = f"Failed to get the latest {MAIN_SERVICE.name} service tag! Status code: {e.code}. Error: {e.reason}"
        print_err(msg)
        send_alert(msg)
        raise Exception()
    except urllib.error.URLError as e:
        print_err("ERROR")
        msg = f"Failed to get the latest {MAIN_SERVICE.name} service tag! Error: {e.reason}"
        print_err(msg)
        send_alert(msg)
        raise Exception()
    print_ok("OK", line_break=False)
    print_secondary(latest_tag)
    return latest_tag


def download_docker_images(latest_tag: str) -> list[int]:
    print("Downloading new images:")

    processes = []
    processes_time = []  # docker images downloading time
    # was start time in processes_time[i] replaced with delta?
    processes_time_calculated = [False] * len(Service)

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
        processes_time.append(time.monotonic())

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
                if not processes_time_calculated[j]:
                    processes_time[j] = round(time.monotonic() - processes_time[j])
                    processes_time_calculated[j] = True
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
                if not processes_time_calculated[j]:
                    processes_time[j] = round(time.monotonic() - processes_time[j])
                    processes_time_calculated[j] = True
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

    return processes_time


def stop_project():
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


def pull_code_changes():
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


def deploy_project():
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


def run_migrations():
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


def update_current_tag_in_file(tag_file, latest_tag):
    print_wait("Updating current tag in the file...")
    with open(tag_file, "w") as file:
        file.write(f"{latest_tag}\n")
    print_ok("OK")


def get_changelog() -> str:
    print_wait("Getting changelog...")
    try:
        req = urllib.request.Request(
            "https://api.github.com/repos/laptop-coder/lost-things-search/releases",
            headers={
                "Authorization": f"Bearer {GITHUB_PAT}",
                "Accept": "application/vnd.github+json",
                "X-GitHub-Api-Version": "2026-03-10",
            },
        )
        with urllib.request.urlopen(req) as response:
            body = json.loads(response.read().decode("utf-8"))[0]["body"]
            lines = [x for x in repr(body)[1:-1].split("\\n") if x]

            content = "\n==== CHANGELOG ====\n"

            header = lines[0]

            content += f"Сравнение с предыдущей версией: {header[header.find('(') + 1 : header.find(')')]}\n"

            for line in lines:
                if line.startswith("###"):  # header (e.g., Bug Fixes)
                    line = line.replace("### ", "")
                    content += f"\n{line}:\n"
                elif line.startswith("*"):  # commit (name and link)
                    line = line.replace("**", "", 2).replace("*", "•", 1)
                    line = line[: line.rfind("[")] + line[line.rfind("]") + 2 : -1]
                    content += f"{line}\n"

    except json.JSONDecodeError:
        print_err("ERROR")
        msg = "Failed to parse JSON response"
        print_err(msg)
        send_alert(msg)
        raise Exception()
    except urllib.error.HTTPError as e:
        print_err("ERROR")
        msg = f"Failed to get changelog from GitHub releases! Status code: {e.code}. Error: {e.reason}"
        print_err(msg)
        send_alert(msg)
        raise Exception()
    except urllib.error.URLError as e:
        print_err("ERROR")
        msg = f"Failed to get changelog from GitHub releases! Error: {e.reason}"
        print_err(msg)
        send_alert(msg)
        raise Exception()
    print_ok("OK")
    return content
