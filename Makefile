# load env variables from .env
ENV_PATH ?= ./.env
ifneq ($(wildcard $(ENV_PATH)),)
    include .env
    export
endif


# Build commands =======================================================================================================

dep:
	go env -w GO111MODULE=on
	go env -w GOPRIVATE=github.com/travelata/*
	go mod tidy

check-lint-installed:
	@if ! [ -x "$$(command -v golangci-lint)" ]; then \
		echo "golangci-lint is not installed"; \
		exit 1; \
	fi; \

lint: check-lint-installed
	@echo Running golangci-lint
	golangci-lint run ./...
	go fmt

check-goose-installed:
	@if ! [ -x "$$(command -v goose)" ]; then \
		echo "goose is not installed"; \
		exit 1; \
	fi; \


build: lint dep  ## builds the main
	go build -o servgen main.go