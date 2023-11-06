##help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST)

download: ## Populate the module cache based on the go.{mod,sum} files.
	go mod download

## run/api: runs the cmd/api application
.PHONY: run/api
run/api:
	go run ./cmd/api

## test: runs all the tests
.PHONY: test
test:
	go test ./...
