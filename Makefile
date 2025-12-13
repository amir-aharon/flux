.PHONY: help tidy run/%

help:
	@echo "Flux Build System"
	@echo ""
	@echo "Usage:"
	@echo "  make run/tempo    - Run the Rate Limiter demo"
	@echo "  make run/swarm    - Run the Worker Pool demo"
	@echo "  make tidy         - Clean up go.mod and go.sum"
	@echo "  make test         - Run all tests in the project"

run/%:
	@echo "Running $* demo..."
	@go run ./examples/$*-demo/

tidy:
	@echo "Tidying modules..."
	@go mod tidy

test:
	@echo "Running all tests..."
	@go test ./...