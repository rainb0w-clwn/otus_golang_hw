GOLANGCI_LINT_VERSION=v1.63.4
BIN := "./bin/calendar"
BIN_MIGRATE := "./bin/calendar_migrate"
BIN_SENDER := "./bin/calendar_sender"
BIN_SCHEDULER := "./bin/calendar_scheduler"
COMPOSE_FILE := "deployments/docker-compose.yml"
INTEGRATION_COMPOSE_FILE := "deployments/docker-compose.integration.yml"
CONFIG_PATH := "./configs/config.yml"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

generate: generate-proto generate-gateway

generate-proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/*.proto

generate-gateway:
	protoc --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
		--grpc-gateway_opt=generate_unbound_methods=true \
		--openapiv2_out=. --openapiv2_opt=generate_unbound_methods=true \
		api/*.proto

build-migrations:
	go build -v -o $(BIN_MIGRATE) ./cmd/migrations

migrate-status: build-migrations
	$(BIN_MIGRATE) -config=$(CONFIG_PATH) status

migrate-up: build-migrations
	$(BIN_MIGRATE) -config=$(CONFIG_PATH) up

migrate-down: build-migrations
	$(BIN_MIGRATE) -config=$(CONFIG_PATH) down

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler

run: migrate-up build
	$(BIN) -config $(CONFIG_PATH)

run-sender: migrate-up build
	$(BIN_SENDER) -config $(CONFIG_PATH)

run-scheduler: migrate-up build
	$(BIN_SCHEDULER) -config $(CONFIG_PATH)

up:
	docker compose -f $(COMPOSE_FILE) up

down:
	docker compose -f $(COMPOSE_FILE) down

reset-db:
	docker compose -f ${INTEGRATION_COMPOSE_FILE} down --volumes
	docker compose -f ${INTEGRATION_COMPOSE_FILE} up -d db rabbitmq

integration-tests: reset-db
	docker compose -f $(INTEGRATION_COMPOSE_FILE) up \
		--abort-on-container-exit \
		--exit-code-from test
	docker compose -f ${INTEGRATION_COMPOSE_FILE} down --rmi local

version: build
	$(BIN) version

test:
	CGO_ENABLED=1 go test --count=1 -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin ${GOLANGCI_LINT_VERSION}

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run run-sender run-scheduler version test lint generate generate-proto generate-gateway up down integration-tests
