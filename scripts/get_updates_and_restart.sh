#!/bin/sh

docker-compose -f "${HOME}/lts-service/compose.yaml" down

docker rmi laptopcoder/lts-service-backend:latest > /dev/null
docker rmi laptopcoder/lts-service-frontend:latest  > /dev/null

docker pull laptopcoder/lts-service-backend:latest > /dev/null
docker pull laptopcoder/lts-service-frontend:latest  > /dev/null

docker-compose -f "${HOME}/lts-service/compose.yaml" up -d
