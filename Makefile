# MSS-Bot Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary name
BINARY_NAME=mss-bot
BINARY_PATH=./cmd/bot

# Build flags
BUILD_FLAGS=-ldflags="-w -s"

# Default config file
CONFIG_FILE=configs/config.kdl

.PHONY: all build clean test coverage lint lint-install help deps

# Default target
all: deps test lint build

## Build the application
build:
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) $(BINARY_PATH)

## Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

## Run tests
test:
	$(GOTEST) -v ./...

## Run tests with coverage
coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

## Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

## Run golangci-lint
lint:
	golangci-lint run

## Install golangci-lint
lint-install:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin)

## Fix linting issues automatically where possible
lint-fix:
	golangci-lint run --fix

## Run the bot
run:
	./$(BINARY_NAME) -config $(CONFIG_FILE)

## Run the bot in development mode
dev: build
	./$(BINARY_NAME) -config configs/config.local.kdl

## Build and run
start: build run

## Docker targets
docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker-compose up

docker-stop:
	docker-compose down

## Full pipeline (what CI should run)
ci: deps test lint build

## Show help
help:
	@echo 'Management commands for $(BINARY_NAME):'
	@echo
	@echo 'Usage:'
	@echo '    make build           Compile the project.'
	@echo '    make clean           Clean build artifacts.'
	@echo '    make test            Run tests.'
	@echo '    make coverage        Run tests with coverage report.'
	@echo '    make deps            Install dependencies.'
	@echo '    make lint            Run golangci-lint.'
	@echo '    make lint-install    Install golangci-lint.'
	@echo '    make lint-fix        Fix linting issues automatically.'
	@echo '    make run             Run the compiled binary.'
	@echo '    make dev             Build and run with local config.'
	@echo '    make start           Build and run with default config.'
	@echo '    make docker-build    Build Docker image.'
	@echo '    make docker-run      Run with docker-compose.'
	@echo '    make docker-stop     Stop docker-compose services.'
	@echo '    make ci              Run full CI pipeline.'
	@echo '    make all             Run deps, test, lint, and build.'
	@echo