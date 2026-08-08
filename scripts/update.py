from collections.abc import Callable
from enum import Enum
from typing import cast
import os
import signal
import subprocess
import sys
import threading
import time


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


path_to_project = f"{os.getenv('HOME')}/lost-things-search"


def run_command(command: str) -> subprocess.CompletedProcess[bytes]:
    """
    Run shell command from the project dir.
    """
    return subprocess.run(
        command,
        shell=True,
        capture_output=True,
        cwd=path_to_project,
    )


class Service(Enum):
    BACKEND = 0
    FRONTEND = 1
    ML = 2
    MIGRATE = 3


current_service = Service.BACKEND
docker_images_downloaded_successfully = False
spinner_th: threading.Thread | None = None


# Spinner
b = ["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"]
spinner = b[0]
prev_n = len(Service) + 1


def animate_spinner(stop_signal: threading.Event):
    global spinner
    global prev_n
    while not stop_signal.is_set():
        for c in b:
            spinner = c
            time.sleep(0.05)
            n = len(Service) - current_service.value + 1
            if n != prev_n:
                print(
                    "\r" + "\033[1A" * prev_n + f"│ \033[32m✔\033[0m" + "\n" * prev_n,
                    end="",
                    flush=True,
                )
                prev_n = n
            print(
                "\r" + "\033[1A" * n + f"│ \033[33m{spinner}\033[0m" + "\n" * n,
                end="",
                flush=True,
            )
    else:
        if docker_images_downloaded_successfully:
            print(
                "\r" + "\033[1A" * prev_n + f"│ \033[32m✔\033[0m" + "\n" * prev_n,
                end="",
                flush=True,
            )


# Create stop signal
stop_signal = threading.Event()


def start_new_thread(fun: Callable, stop_signal: threading.Event) -> threading.Thread:
    th = threading.Thread(target=fun, args=(stop_signal,))
    th.start()
    return th


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
    result = run_command(
        f"""curl -s "https://ghcr.io/token?scope=repository:{repo}:pull" | grep -o '"token":"[^"]*"' | cut -d '"' -f4""",
    )
    if result.returncode != 0:
        err = result.stderr.decode("utf-8")
        print_err("ERROR")
        print_err(f"Failed to get GHCR token! Error: {err}")
        return False
    token = result.stdout.decode("utf-8")
    print_ok("OK")

    # Get the latest tag
    print_wait(f"Trying to get the latest {service.name} service tag...")
    result = run_command(
        f"""
        curl -s -H "Authorization: Bearer {token}" "https://ghcr.io/v2/{repo}/tags/list" | grep -o '"tags":\\[[^]]*\\]' | grep -o '"[^"]*"' | tail -1 | tr -d '"'
        """,
    )
    if result.returncode != 0:
        err = result.stderr.decode("utf-8")
        print_err("ERROR")
        print_err(f"Failed to get tag! Error: {err}")
        return False
    latest_tag = result.stdout.decode("utf-8")[:-1]
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

            # Print services
            print("╭" + "─" * (max([len(s.name) for s in Service]) + 4) + "╮")
            for s in Service:
                print(
                    f"│ \033[33m⠼\033[0m {s.name}{' ' * (max([len(x.name) for x in Service]) - len(s.name) + 1)}│"
                )
            print("╰" + "─" * (max([len(x.name) for x in Service]) + 4) + "╯")

            global current_service

            # Start animation
            th = start_new_thread(animate_spinner, stop_signal)
            global spinner_th
            spinner_th = th
            # Download
            for s in Service:
                current_service = s
                repo = f"laptop-coder/lost-things-search-{s.name}".lower()
                result = run_command(
                    f"""
                    docker pull "ghcr.io/{repo}:{latest_tag}"
                    """,
                )
                if result.returncode != 0:
                    err = result.stderr.decode("utf-8")
                    print_err("ERROR")
                    print_err(f"Failed to pull {s.name} image! Error: {err}")
                    return False

            # Stop animation
            global docker_images_downloaded_successfully
            docker_images_downloaded_successfully = True
            stop_signal.set()
            th.join()

            # Stop the project
            print_wait("Stopping the project...")
            result = run_command("make down")
            if result.returncode != 0:
                err = result.stderr.decode("utf-8")
                print_err("ERROR")
                print_err(f"Failed to stop the project! Error: {err}")
                return False
            print_ok("OK")

            # Pull the code changes
            print_wait("Pulling the code changes...")
            result = run_command("git pull")
            if result.returncode != 0:
                err = result.stderr.decode("utf-8")
                print_err("ERROR")
                print_err(f"Failed to pull the code changes! Error: {err}")
                return False
            print_ok("OK")

            # Deploy the project
            print_wait("Deploying the project...")
            result = run_command("make deploy")
            if result.returncode != 0:
                err = result.stderr.decode("utf-8")
                print_err("ERROR")
                print_err(f"Failed to deploy the project! Error: {err}")
                return False
            print_ok("OK")

            # Run migrations
            print_wait("Running migrations...")
            result = run_command("make migrate")
            if result.returncode != 0:
                err = result.stderr.decode("utf-8")
                print_err("ERROR")
                print_err(f"Failed to run migrations! Error: {err}")
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
    stop_signal.set()
    sys.exit(0)


if __name__ == "__main__":
    signal.signal(signal.SIGINT, graceful_shutdown)
    signal.signal(signal.SIGTERM, graceful_shutdown)
    for i in range(10):
        if main():
            print_ok("Done!")
            break
        stop_signal.set()
        spinner_th = cast(threading.Thread | None, spinner_th)
        if spinner_th is not None:
            spinner_th.join()
        print(
            f"{i + 1}/10 attempt. Waiting for 10 seconds to run script one more time...",
            flush=True,
        )
        time.sleep(10)
        # Reset variables
        stop_signal.clear()
        spinner_th = None
        prev_n = len(Service) + 1
        docker_images_downloaded_successfully = False
