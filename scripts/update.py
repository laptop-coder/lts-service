import os
import subprocess


def print_wait(s: str) -> None:
    print(s, end=" ")


def print_ok(s: str, line_break=True) -> None:
    end = "\n" if line_break else " "
    print(f"\033[32m{s}\033[0m", end=end)


def print_err(s: str, line_break=True) -> None:
    end = "\n" if line_break else " "
    print(f"\033[31m{s}\033[0m", end=end)


def print_secondary(s: str, line_break=True) -> None:
    end = "\n" if line_break else " "
    print(f"\033[2m{s}\033[0m", end=end)


def alg() -> bool:
    path_to_project = f"{os.getenv('HOME')}/lost-things-search"

    print_wait(f"Got the project dir.")
    print_secondary(path_to_project)

    tag_file = f"{path_to_project}/.tag"
    tag_file_is_new = False

    if not os.path.exists(tag_file):
        print_wait("Tag file not found")
        print_secondary(tag_file)
        tag_file_is_new = True

    service = "backend"
    print(f"Looking at the {service} service.")
    repo = f"laptop-coder/lost-things-search-{service}"

    # Get GHCR token
    print_wait("Trying to get GHCR token...")
    result = subprocess.run(
        f"""curl -s "https://ghcr.io/token?scope=repository:{repo}:pull" | grep -o '"token":"[^"]*"' | cut -d '"' -f4""",
        shell=True,
        capture_output=True,
    )
    if result.returncode != 0:
        err = result.stderr.decode("utf-8")
        print_err("ERROR")
        print_err(f"Failed to get GHCR token! Error: {err}")
        print("Running script one more time...")
        return False
    token = result.stdout.decode("utf-8")
    print_ok("OK")

    # Get the latest tag
    print_wait(f"Trying to get the latest {service} service tag...")
    result = subprocess.run(
        f"""
        curl -s -H "Authorization: Bearer {token}" "https://ghcr.io/v2/{repo}/tags/list" | grep -o '"tags":\\[[^]]*\\]' | grep -o '"[^"]*"' | tail -1 | tr -d '"'
        """,
        shell=True,
        capture_output=True,
    )
    if result.returncode != 0:
        err = result.stderr.decode("utf-8")
        print_err("ERROR")
        print_err(f"Failed to get tag! Error: {err}")
        print("Running script one more time...")
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
        print_wait(f"Getting current {service} tag...")
        with open(tag_file, "r") as file:
            current_tag = file.readline()[:-1]
        print_ok("OK", line_break=False)
        print_secondary(current_tag)
        print_wait("Compairing tags...")
        if current_tag != latest_tag:
            print_secondary(f"{current_tag} != {latest_tag}")
            print("Downloading new images:")
            for service in ["backend", "frontend", "ml", "migrate"]:
                print_wait(f"- {service}")
                result = subprocess.run(
                    f"""
                    docker pull "ghcr.io/{repo}:{latest_tag}"
                    """,
                    shell=True,
                    capture_output=True,
                )
                if result.returncode != 0:
                    err = result.stderr.decode("utf-8")
                    print_err("ERROR")
                    print_err(f"Failed to pull {service} image! Error: {err}")
                    print("Running script one more time...")
                    return False
                print_ok("OK")

            # Stop the project
            print_wait("Stopping the project...")
            result = subprocess.run(
                f"""
                make down
                """,
                shell=True,
                capture_output=True,
                cwd=path_to_project,
            )
            if result.returncode != 0:
                err = result.stderr.decode("utf-8")
                print_err("ERROR")
                print_err(f"Failed to stop the project! Error: {err}")
                print("Running script one more time...")
                return False
            print_ok("OK")

            # Pull the code changes
            print_wait("Pulling the code changes...")
            result = subprocess.run(
                f"""
                git pull
                """,
                shell=True,
                capture_output=True,
                cwd=path_to_project,
            )
            if result.returncode != 0:
                err = result.stderr.decode("utf-8")
                print_err("ERROR")
                print_err(f"Failed to pull the code changes! Error: {err}")
                print("Running script one more time...")
                return False
            print_ok("OK")

            # Deploy the project
            print_wait("Deploying the project...")
            result = subprocess.run(
                f"""
                make deploy
                """,
                shell=True,
                capture_output=True,
                cwd=path_to_project,
            )
            if result.returncode != 0:
                err = result.stderr.decode("utf-8")
                print_err("ERROR")
                print_err(f"Failed to deploy the project! Error: {err}")
                print("Running script one more time...")
                return False
            print_ok("OK")

            # Run migrations
            print_wait("Running migrations...")
            result = subprocess.run(
                f"""
                make migrate
                """,
                shell=True,
                capture_output=True,
                cwd=path_to_project,
            )
            if result.returncode != 0:
                err = result.stderr.decode("utf-8")
                print_err("ERROR")
                print_err(f"Failed to run migrations! Error: {err}")
                print("Running script one more time...")
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


for _ in range(10):
    success = alg()
    if success:
        print_ok("Done!")
        break
