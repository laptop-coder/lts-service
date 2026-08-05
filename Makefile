.PHONY: migrate cron deploy first-run logs down dev
COMPOSE := $(shell command -v docker compose > /dev/null 2>&1 && echo "docker compose" || echo "docker-compose")

migrate:
	$(COMPOSE) --profile migrate up migrate --exit-code-from migrate
	$(COMPOSE) rm -f migrate

cron:
	sh -c "( crontab -l; cat ./crontab.tasks )" | crontab -

deploy:
	$(COMPOSE) up -d

first-run:
	$(MAKE) migrate
	$(MAKE) deploy
	$(MAKE) cron

logs:
	$(COMPOSE) logs -f

down:
	$(COMPOSE) down

dev:
	$(COMPOSE) -f ./dev.compose.yaml --profile migrate up --build

