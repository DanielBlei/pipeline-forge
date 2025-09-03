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

.PHONY: test
test: ## Run tests
	@uv run pytest . -v