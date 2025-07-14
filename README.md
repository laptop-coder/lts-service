# LostThingsSearch service

[![CodeFactor](https://www.codefactor.io/repository/github/laptop-coder/lts-service/badge)](https://www.codefactor.io/repository/github/laptop-coder/lts-service)

## Installation (dev)

0. Install requirements:
   - Docker
   - Docker Compose,
   - Git
   - [mkcert](https://github.com/FiloSottile/mkcert)
   - any text editor
1. `git clone https://github.com/laptop-coder/lts-service.git` — clone this
   repository
2. `cd ./lts-service` — go to the project directory
3. Make SSL-certificate for `https` in dev env:
   1. `mkcert -install` — install the local CA in the system trust store (you
      need to do this once, for example, after `mkcert` installation)
   2. `mkcert 172.16.1.2 172.16.1.3` — generate certificate and key for backend
      and frontend
   3. Move and rename certificate and key:
      ```
      mkdir -p ./data/env
      mv ./172.16.1.2+1.pem ./data/env/certfile.crt
      mv ./172.16.1.2+1-key.pem ./data/env/keyfile.key
      ```
4. `cp ./.env.example ./.env` — copy `.env.example` to `.env`
5. `vi ./.env` (for example) — edit `.env` to your preferences
6. `docker compose up` — run project
7. Open https://172.16.1.3 in your browser
