

##@ Development

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
