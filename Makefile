.PHONY: test coverage lint fmt vet build clean help

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

test: ## Run tests
	go test -v -race $(shell go list ./... | grep -v /examples)

coverage: ## Generate coverage report
	go test -coverprofile=coverage.out -covermode=atomic $(shell go list ./... | grep -v /examples)
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@go tool cover -func=coverage.out | grep total | awk '{print "Total coverage: " $$3}'

lint: ## Run linter
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run $(shell go list ./... | grep -v /examples)

fmt: ## Format code
	go fmt $(shell go list ./... | grep -v /examples)

vet: ## Run go vet
	go vet $(shell go list ./... | grep -v /examples)

build: ## Build (validation)
	go build $(shell go list ./... | grep -v /examples)

clean: ## Clean build artifacts
	rm -f coverage.out coverage.html
	go clean

.DEFAULT_GOAL := help
