import os
import subprocess
from .config import PATH_TO_PROJECT, TAG_FILE
import sys


def graceful_shutdown(signum, frame):
    del signum, frame  # ignore parameters
    print("\n───────────────────────────────────")
    print(f"Received a signal, shutting down...")
    sys.exit(0)


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
        cwd=PATH_TO_PROJECT,
    )


def is_tag_file_new() -> bool:
    if not os.path.exists(TAG_FILE):
        print_wait("Tag file not found")
        print_secondary(TAG_FILE)
        return True
    return False
