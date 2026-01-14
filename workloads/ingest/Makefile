##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Development

.PHONY: check
check: ## Run linters (Ruff and MyPy)
	@uv run ruff check . 
	@uv run ruff format --check .
	@uv run mypy .

.PHONY: fix
fix: ## Attempt to fix linting errors (Ruff)
	@uv run ruff check . --fix
	@uv run ruff format .

.PHONY: install
install: ## Install the project
	@uv sync
	@uv pip install -e .

.PHONY: run-help
run-help: ## Run the project and show the help
	@uv run ingest --help


##@ Build

.PHONY: build
build: ## Build the project
	@uv run build

.PHONY: build-docker
build-docker: ## Build the project as a Docker image
	@docker build -t ingest .

##@ Tests

.PHONY: test-docker
test-docker: ## Run the project as a Docker container
	@docker run -it ingest ingest --help

.PHONY: test
test: ## Run tests
	@uv run pytest . -v --junitxml=test-results.xml

.PHONY: test-coverage
test-coverage: ## Run tests with coverage and open HTML report
	@uv run pytest . --cov=ingest --cov-report=html --cov-report=term-missing
	@echo "Opening coverage report in browser..."
	@xdg-open htmlcov/index.html || open htmlcov/index.html || start htmlcov/index.html
