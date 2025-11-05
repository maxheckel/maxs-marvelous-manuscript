.PHONY: all build clean run-recorder run-web dev install-frontend build-frontend test

# Build everything
all: build build-frontend

# Build Go binaries
build:
	@echo "Building Go applications..."
	@mkdir -p bin
	go build -o bin/recorder ./cmd/recorder
	go build -o bin/web ./cmd/web
	@echo "Build complete!"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -rf web/frontend/dist/
	rm -rf web/frontend/node_modules/
	rm -f data/*.db
	@echo "Clean complete!"

# Install frontend dependencies
install-frontend:
	@echo "Installing frontend dependencies..."
	cd web/frontend && npm install
	@echo "Frontend dependencies installed!"

# Build frontend
build-frontend:
	@echo "Building frontend..."
	cd web/frontend && npm run build
	@echo "Frontend build complete!"

# Run the recorder
run-recorder: build
	@echo "Starting recorder..."
	./bin/recorder

# Run the web server
run-web: build build-frontend
	@echo "Starting web server..."
	./bin/web

# Development mode - runs web server and frontend dev server
dev:
	@echo "Starting development servers..."
	@echo "Backend will run on http://localhost:8080"
	@echo "Frontend will run on http://localhost:5173"
	@echo ""
	@trap 'kill 0' EXIT; \
		go run ./cmd/web & \
		cd web/frontend && npm run dev

# Run Go tests
test:
	go test ./...

# Download Go dependencies
deps:
	go mod download
	go mod tidy

# Format Go code
fmt:
	go fmt ./...

# Lint Go code (requires golangci-lint)
lint:
	golangci-lint run

# Create data directory
init:
	@mkdir -p data
	@echo "Data directory created!"

# Show help
help:
	@echo "Available targets:"
	@echo "  make all             - Build everything (Go + frontend)"
	@echo "  make build           - Build Go binaries"
	@echo "  make build-frontend  - Build Vue frontend"
	@echo "  make install-frontend- Install frontend npm dependencies"
	@echo "  make clean           - Clean build artifacts"
	@echo "  make run-recorder    - Run the recorder app"
	@echo "  make run-web         - Run the web server"
	@echo "  make dev             - Run in development mode"
	@echo "  make test            - Run Go tests"
	@echo "  make deps            - Download Go dependencies"
	@echo "  make fmt             - Format Go code"
	@echo "  make lint            - Lint Go code"
	@echo "  make init            - Create data directory"
	@echo "  make help            - Show this help"
