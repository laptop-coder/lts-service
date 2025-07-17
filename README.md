# LostThingsSearch service

[![CodeFactor](https://www.codefactor.io/repository/github/laptop-coder/lts-service/badge)](https://www.codefactor.io/repository/github/laptop-coder/lts-service)

## Installation (dev)

0. Install requirements:
   - Docker
   - Docker Compose,
   - Git
   - [mkcert](https://github.com/FiloSottile/mkcert)
   - any text editor
1. Clone this repository:
   ```
   git clone https://github.com/laptop-coder/lts-service.git
   ```
2. Go to the project directory:
   ```
   cd ./lts-service
   ```
3. Make SSL-certificate for `https` in dev env:
   1. Install the local CA in the system trust store (you need to do this once,
      for example, after `mkcert` installation):
      ```
      mkcert -install
      ```
   2. Generate certificate and key for backend and frontend:
      ```
      mkcert 172.16.1.2 172.16.1.3
      ```
   3. Move and rename certificate and key:
      ```
      mkdir -p ./data/env
      mv ./172.16.1.2+1.pem ./data/env/certfile.crt
      mv ./172.16.1.2+1-key.pem ./data/env/keyfile.key
      ```
4. Copy `.env.example` to `.env`:
   ```
   cp ./.env.example ./.env
   ```
5. Edit `.env` to your preferences (example):
   ```
   vi ./.env
   ```
6. Run project:
   ```
   docker compose -f ./compose-dev.yaml up
   ```
7. Open https://172.16.1.3 in your browser
