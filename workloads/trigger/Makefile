##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: lint
lint: ## Run Linters
	golangci-lint run -c ../.golangci.yml

.PHONY: lint-fix
lint-fix: ## Run Linters and fix issues
	golangci-lint run --fix -c ../.golangci.yml

.PHONY: build
build: ## Build the trigger
	go build -o bin/trigger-workload main.go

.PHONY: run
run:  ## Run the trigger
	go run main.go --help

.PHONY: test
test: ## Run tests
	go test ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	go test ./... -coverprofile=coverage.out
