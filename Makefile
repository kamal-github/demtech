APP_NAME := demtech
BUILD_DIR := bin
MAIN_PACKAGE := ./cmd/main.go
GO_FILES := $(shell find . -name '*.go' -not -path "./vendor/*")

.PHONY: all build run test fmt lint clean install-deps up down restart dbuild e2e

all: build

## Runs the server in docker
up: down
	@echo "Running docker compose up"
	docker-compose up --build

## Destroys all the containers
down:
	@echo "Running docker compose down"
	docker-compose down -v

## Restart the docker containers
restart: down up

## Build docker image of Go application
dbuild:
	docker build -t $(APP_NAME) .

## Build the Go application
build:
	@echo "Building the application..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PACKAGE)

## Run the compiled binary
run: build
	@echo "Running the application..."
	@$(BUILD_DIR)/$(APP_NAME)

## Run unit tests with coverage
unit-test:
	@echo "Running unit tests only..."
	go test ./... -cover -race

## Run all tests (unit & integration)
test:
	@echo "Running all the tests (unit & integration)..."
	docker-compose -f docker-compose.test.yml down -v
	docker-compose -f docker-compose.test.yml up -d --remove-orphans
	go test -tags=integration ./... -cover -race
	docker-compose -f docker-compose.test.yml down -v

e2e:
	@echo "Running E2E tests..."
	docker-compose -f docker-compose.e2e.yml down -v
	docker-compose -f docker-compose.e2e.yml up --abort-on-container-exit --build
	docker-compose -f docker-compose.e2e.yml down -v

## Format the Go code
fmt:
	@echo "Formatting code..."
	go fmt ./... ## Format code

## Run linter (requires golangci-lint)
lint:
	@echo "Linting code..."
	golangci-lint run ## Run linter

## Clean build artifacts
clean:
	@echo "Cleaning up..."
	rm -rf $(BUILD_DIR)/*
	go clean ## Clean build artifacts

## Install dependencies
install-deps:
	@echo "Installing dependencies..."
	go mod tidy ## Install dependencies
