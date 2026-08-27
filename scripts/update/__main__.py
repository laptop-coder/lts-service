from .utils import (
    graceful_shutdown,
    print_wait,
    print_ok,
    print_secondary,
    is_tag_file_new,
)
from .tools import (
    get_ghcr_token,
    get_latest_tag,
    download_docker_images,
    stop_project,
    pull_code_changes,
    deploy_project,
    run_migrations,
    update_current_tag_in_file,
    publish_update_digest,
    get_changelog,
)
from .config import PATH_TO_PROJECT, MAIN_SERVICE, TAG_FILE, DigestDTOBuilder
import signal
import time
from .vk_api import send_alert


def main(digest_dto_builder: DigestDTOBuilder) -> bool:
    print("\033c", end="")
    print_wait(f"Got the project dir.")
    print_secondary(PATH_TO_PROJECT)
    print(f"Looking at the {MAIN_SERVICE.name} service.")

    try:
        token = get_ghcr_token()
    except:
        return False

    try:
        latest_tag = get_latest_tag(token)
    except:
        return False

    digest_dto_builder.latest_tag = latest_tag

    need_update = False
    if is_tag_file_new():
        print_wait("Adding latest tag to the file...")
        with open(TAG_FILE, "w") as file:
            file.write(f"{latest_tag}\n")
        print_ok("OK")
        need_update = True
    else:
        print_wait(f"Getting current {MAIN_SERVICE.name} tag...")
        with open(TAG_FILE, "r") as file:
            current_tag = file.readline()[:-1]
        print_ok("OK", line_break=False)
        digest_dto_builder.old_tag = current_tag
        print_secondary(current_tag)
        print_wait("Compairing tags...")
        if current_tag != latest_tag:
            print_secondary(f"{current_tag} != {latest_tag}")
            need_update = True
        else:
            print_secondary(f"{current_tag} == {latest_tag}")

    if need_update:
        digest_dto_builder.downloading_time = download_docker_images(latest_tag)

        stop_project()
        pull_code_changes()
        deploy_project()
        run_migrations()
        update_current_tag_in_file(TAG_FILE, latest_tag)
        digest_dto_builder.changelog = get_changelog()

    return True


if __name__ == "__main__":
    signal.signal(signal.SIGINT, graceful_shutdown)
    signal.signal(signal.SIGTERM, graceful_shutdown)
    start_time = time.monotonic()
    digest_dto_builder = DigestDTOBuilder()
    for i in range(10):
        success = main(digest_dto_builder)
        if success:
            digest_dto_builder.script_time = round(time.monotonic() - start_time)
            try:
                digest_dto = digest_dto_builder.build()
            except:
                # This code branch should not be executed
                msg = "Script error"
                send_alert(msg)
                print(msg, flush=True)
                time.sleep(5)
                continue
            try:
                publish_update_digest(digest_dto)
            except:
                time.sleep(5)
                continue
            print_ok("Done!")
            break
        if i + 1 < 10:
            msg = f"{i + 1}/10 attempt. Waiting for 10 seconds to run script one more time..."
            send_alert(msg)
            print(
                msg,
                flush=True,
            )
            time.sleep(10)
        else:
            msg = f"{i + 1}/10 attempt. Failed to update. Shutting down..."
            send_alert(msg)
            print(msg, flush=True)
