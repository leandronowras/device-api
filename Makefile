.PHONY: help build run test clean

clean:
	@echo "Cleaning up..."
	@rm -rf bin

test:
	@echo "Running tests..."
	@go test ./...

run: build
	@echo "Running the application..."
	@./bin/app

build:
	@echo "Building the application..."
	@go build -o bin/app ./cmd/app

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build    Compile the application"
	@echo "  run      Run the compiled application"
	@echo "  test     Run the tests"
	@echo "  clean    Remove the build artifacts"
