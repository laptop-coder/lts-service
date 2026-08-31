COMPOSE := $(shell command -v docker compose > /dev/null 2>&1 && echo "docker compose" || echo "docker-compose")
PYTHON := $(shell command -v python > /dev/null 2>&1 && echo "python" || echo "python3")

.PHONY: migrate
migrate: ## run database migrations using the migrate profile
	$(COMPOSE) --profile migrate up migrate --exit-code-from migrate
	$(COMPOSE) rm -f migrate

.PHONY: cron
cron: ## install nightly cron job from crontab.tasks for automatic updates
	sh -c "( crontab -l; cat ./crontab.tasks )" | crontab -

.PHONY: deploy
deploy: ## start all application services in detached mode
	$(COMPOSE) up -d

.PHONY: first-run
first-run: ## perform initial setup: migrate, start services, and install cron
	$(MAKE) deploy
	$(MAKE) migrate
	$(MAKE) cron

.PHONY: logs
logs: ## follow logs from all running docker compose services
	$(COMPOSE) logs -f

.PHONY: down
down: ## stop and remove docker compose services
	$(COMPOSE) down

.PHONY: dev
dev: ## run development stack with build and migration profile enabled
	$(COMPOSE) -f ./dev.compose.yaml --profile migrate up --build

.PHONY: dev-down
dev-down: ## stop and remove development docker compose services
	$(COMPOSE) -f ./dev.compose.yaml down

.PHONY: update
update: ## update docker images, pull the code changes and restart the project
	cd ./scripts && $(PYTHON) -m update

.PHONY: help
help: ## show available make targets with short descriptions
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

