.PHONY: build test

build:
	@echo "Building the project..."
	go build -o ./bin/envcraft .

test:
	@echo "Running tests..."
	go test ./...