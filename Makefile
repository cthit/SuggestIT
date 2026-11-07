# Developer shortcuts. The backend binary is built on the host and
# bind-mounted into the backend container (see docker-compose.yml).

COMPOSE ?= docker compose

.PHONY: build up down restart logs frontend

## build: build the backend binary (backend/bin/suggestit)
build:
	$(MAKE) -C backend build

## up: build the backend and start db, adminer and backend in the background
up: build
	$(COMPOSE) up -d --build

## down: stop and remove the containers (the database volume is kept)
down:
	$(COMPOSE) down

## restart: rebuild the backend binary and restart the backend container
restart: build
	$(COMPOSE) restart backend

## logs: follow the logs of all services
logs:
	$(COMPOSE) logs -f

## frontend: start the Vite dev server on http://localhost:3000
frontend:
	cd frontend && pnpm dev
