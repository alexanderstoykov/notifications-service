-include .env

.PHONY: build
build:
	@docker build  -f Dockerfile -t notifications-service .

.PHONY: init
init:
	@cp .env.dist .env
	@cp .env.test.dist .env.test
.PHONY: up

up:
	$(eval SERVICE = ${s})
	@docker-compose up -d --no-build --remove-orphans ${SERVICE}
	@docker-compose ps

down:
	@docker-compose down --remove-orphans --volumes

test:
	@docker-compose up -d --no-build --remove-orphans postgres-test
	@go test ./... -v --tags=integration -count=1