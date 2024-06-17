DSN="host=localhost port=5432 user=postgres password=password dbname=goapi sslmode=disable timezone=UTC connect_timeout=5"
BINARY_NAME=goapi


build: ## Build will build binary for the application
	@echo "Building back end..."
	go build -o ${BINARY_NAME} ./cmd/api/
	@echo "Binary built!"


run: build ## Run builds and runs the application
	@echo "Starting back end..."
	@env DSN=${DSN} ./${BINARY_NAME} &
	@echo "Back end started!"


clean: ## Clean runs go clean and deletes binaries
	@echo "Cleaning..."
	@go clean
	@rm ${BINARY_NAME}
	@echo "Cleaned!"


start: run ## Start the application


stop: ## Stops the running application
	@echo "Stopping back end..."
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
	@echo "Stopped back end!"


restart: stop start ## Stops and starts the running application

.PHONY: build run clean start stop

help: ## Display details on all commands
	@awk 'BEGIN {FS = ":.*?##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n%s\n", substr($$0, 5) } ' $(MAKEFILE_LIST)