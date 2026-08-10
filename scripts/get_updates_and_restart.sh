#!/bin/sh
PATH_TO_PROJECT="$HOME/lost-things-search"
TAGS_FILE="${PATH_TO_PROJECT}/.tags"
NEED_RESTART=0
TAGS_FILE_IS_NEW=0

if [ ! -f "${TAGS_FILE}" ]; then
    echo -e "The new tags file will be created: \033[2m${TAGS_FILE}\033[0m"
    TAGS_FILE_IS_NEW=1
fi

for service in backend migrate ml frontend; do
    echo "Looking at the ${service} service..."
    REPO="laptop-coder/lost-things-search-${service}"
    echo "Trying to get GHCR token..."
    if ! TOKEN=$(curl -s "https://ghcr.io/token?scope=repository:${REPO}:pull" | grep -o '"token":"[^"]*"' | cut -d '"' -f4); then
        echo -e "\033[31mFailed to get GHCR token!\033[0m"
        continue
    fi
    echo -e "\033[1ATrying to get GHCR token... \033[32mOK\033[0m"
    echo "Trying to get the latest ${service} tag..."
    if ! LATEST_TAG=$(curl -s -H "Authorization: Bearer ${TOKEN}" "https://ghcr.io/v2/${REPO}/tags/list" | grep -o '"tags":\[[^]]*\]' | grep -o '"[^"]*"' | tail -1 | tr -d '"'); then
        echo -e "\033[31mFailed to get the latest tag of the ${service} service!\033[0m"
        continue
    fi
    echo -e "\033[1ATrying to get the latest ${service} tag... \033[32mOK\033[0m \033[2m${LATEST_TAG}\033[0m"

    if [ "${TAGS_FILE_IS_NEW}" -ne 0 ]; then
        echo "Adding ${service} latest tag to the tags file..."
        echo "${service}:${LATEST_TAG}" >> "${TAGS_FILE}"
    else
        echo "Getting current ${service} tag..."
        CURRENT_TAG=$(grep "^${service}:" "${TAGS_FILE}" | cut -d ':' -f2)
        echo -e "\033[1AGetting current ${service} tag... \033[32mOK\033[0m \033[2m${CURRENT_TAG}\033[0m"
        if [ "${CURRENT_TAG}" != "${LATEST_TAG}" ]; then
            echo "The tags are different, updating."
            echo "Downloading new image..."
            if ! docker pull "ghcr.io/${REPO}:${LATEST_TAG}"; then
                echo -e "\033[31mFailed to pull image from GHCR!\033[0m"
                continue
            fi
            echo -e "\033[1ADownloading new image... \033[32mOK\033[0m"
            NEED_RESTART=1
            echo "Updating current tag in the tags file..."
            sed -i "s/^${service}:.*/${service}:${LATEST_TAG}/" "${TAGS_FILE}"
            echo -e "\033[1AUpdating current tag in the tags file... \033[32mOK\033[0m"
        fi
    fi
done


if [ "${NEED_RESTART}" -ne 0 ]; then
    echo "Going to the project directory..."
    if ! cd "${PATH_TO_PROJECT}"; then
        echo -e "\033[31mFailed to change directory to ${PATH_TO_PROJECT}!\033[0m Exiting..."
        exit 1
    fi
    echo -e "\033[1AGoing to the project directory... \033[32mOK\033[0m (\033[2m${PATH_TO_PROJECT}\033[0m)"
    echo "Stopping the project..."
    if ! make down; then
        echo -e "\033[31mFailed to stop the project!\033[0m Trying to up and exit..."
        make deploy
        exit 1
    fi
    echo -e "\033[1AStopping the project... \033[32mOK\033[0m"
    echo "Pulling the code changes..."
    git pull || echo -e "\033[31mFailed to pull the code changes!\033[0m"
    echo "Deploying the project..."
    if ! make deploy; then
        echo -e "\033[31mFailed to deploy the project!\033[0m Exiting..."
        exit 1
    fi
    echo -e "\033[1ADeploying the project... \033[32mOK\033[0m"
fi

