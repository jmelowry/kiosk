.PHONY: test build

test:
	@echo "Running tests..."
	@go test ./... -v
	@if [ $$? -ne 0 ]; then \
		echo "Tests failed. Commit aborted."; \
		exit 1; \
	fi
	@echo "All tests passed. Proceeding with commit."