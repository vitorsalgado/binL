.ONESHELL:
.DEFAULT_GOAL := help

# allow user specific optional overrides
-include Makefile.overrides

export

.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: up
up: ## start all systems
	@docker-compose up --build --force-recreate

.PHONY: down
down: ## close everything
	@docker-compose down

.PHONY: build
build: ## builds the reader application
	@go build -o build/ ./...

.PHONY: rm
rm: ## cleanup everything
	@docker-compose down --remove-orphans --volumes

.PHONY: run
run: ## run api locally
	@go run ./...

.PHONY: deps
deps: ## download dependencies
	@go mod download
