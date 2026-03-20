ci-lint:
	@echo "Running linter..."
	@golangci-lint run

lint:
	@echo "Running linter..."
	@golangci-lint run --fix