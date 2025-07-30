# Подготовка к работе

1. Клонируйте репозиторий:
   ```
   git clone https://github.com/laptop-coder/lts-service.git
   ```
2. Перейдите в него:
   ```
   cd ./lts-service
   ```
3. Создайте SSL-сертификаты для `https` в рабочем окружении:
   1. Установите локальный центр сертификации в системное хранилище (trust
      store):

      ```
      mkcert -install
      ```

   2. Сгенерируйте сертификат и ключ для бэкенда и фронтенда:
      ```
      mkcert 172.16.1.2 172.16.1.3
      ```
   3. Переместите и переименуйте их:
      ```
      mkdir -p ./data/env
      mv ./172.16.1.2+1.pem ./data/env/certfile.crt
      mv ./172.16.1.2+1-key.pem ./data/env/keyfile.key
      ```

4. Скопируйте `.env.example` в `.env`:
   ```
   cp ./.env.example ./.env
   ```
5. Отредактируйте `.env`, например:
   ```
   vi ./.env
   ```
6. Запустите проект:
   ```
   docker compose -f ./compose-dev.yaml up
   ```
7. Готово!

- https://172.16.1.2 — бэкенд
- https://172.16.1.3 — фронтенд
