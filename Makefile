##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Pipeline Forge Development

.PHONY: test
test: test-operator ## Run Pipeline Forge tests

.PHONY: dev-up
dev-up: ## Start development environment with docker-compose (dev/)
	docker-compose -f dev/docker-compose.yml up -d

.PHONY: dev-logs
dev-logs: ## Follow the logs of the development environment
	docker-compose -f dev/docker-compose.yml logs -f

.PHONY: dev-down
dev-down: ## Stop development environment with docker-compose (dev/)
	docker-compose -f dev/docker-compose.yml down

.PHONY: dev-cleanup
dev-cleanup: ## Clean up development environment (containers, volumes, networks)
	docker-compose -f dev/docker-compose.yml down -v
	docker volume ls | grep -E "(dev_|pipeline-forge)" | awk '{print $$2}' | xargs -r docker volume rm
	docker ps -a | grep -E "(pipeline-forge|mysql|postgres)" | awk '{print $$1}' | xargs -r docker rm -f


##@ Operator Development

.PHONY: generate
generate: ## Generate the operator code
	${MAKE} -C operator generate

.PHONY: gen-and-install-crds
gen-and-install-crds: ## Generate and install CRDs in the current context
	${MAKE} -C operator generate install 

.PHONY: uninstall-crds
uninstall-crds: ## Uninstall CRDs in the current context
	${MAKE} -C operator uninstall

.PHONY: run-operator
run-operator: ## Run the operator locally
	${MAKE} -C operator run

.PHONY: test-operator
test-operator: ## Run the operator tests (not the e2e tests)
	${MAKE} -C operator test

.PHONY: build-and-deploy-operator
build-and-deploy-operator: ## Build and deploy the operator binary
	${MAKE} -C operator build deploy

.PHONY: undeploy-operator
undeploy-operator: ## Undeploy the operator binary
	${MAKE} -C operator undeploy
