# Binary name
API_BIN=api
SUMMARIZER_WORKER_BIN=summarizer-worker
QUEUE_GENERATOR_BIN=queue-generator

# Build directory
BUILD_DIR=build

# Get the current commit hash
COMMIT=$(shell git rev-parse --short HEAD)

# Get the current date
DATE=$(shell date +%Y-%m-%d_%H:%M:%S)

# Common build flags
BUILD_FLAGS=-ldflags "-X main.commit=${COMMIT} -X main.buildTime=${DATE}"

# Tools versions
GOLANGCI_LINT_VERSION=v1.55.2

# Install required tools
.PHONY: install-tools
install-tools:
	@echo "Installing tools..."
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.61.0
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest

# Lint the code
.PHONY: lint
lint:
	@echo "Running linters..."
	golangci-lint run ./...

# Check for security issues
.PHONY: security-check
security-check:
	@echo "Running security checks..."
	gosec ./...
	govulncheck ./...

# Run go vet
.PHONY: vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run all checks
.PHONY: check-all
check-all: fmt vet 
	@echo "All checks completed!"

# Ensure build directory exists
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Clean build directory
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# Run tests with race detection and coverage
.PHONY: test
test:
	@echo "Running tests..."
	go test -v -race -cover ./...

# Build for current platform
.PHONY: build
build: check-all $(BUILD_DIR)
	go build ${BUILD_FLAGS} -o $(BUILD_DIR)/$(API_BIN)

# Build for Linux (amd64)
.PHONY: build-linux-amd64
build-linux-amd64: check-all $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build ${BUILD_FLAGS} -o $(BUILD_DIR)/$(API_BIN)-linux-amd64

# Build for Linux (arm64)
.PHONY: build-linux-arm64
build-linux-arm64: check-all $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 go build ${BUILD_FLAGS} -o $(BUILD_DIR)/$(API_BIN)-linux-arm64

# Build for macOS (amd64)
.PHONY: build-darwin-amd64
build-darwin-amd64: check-all $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 go build ${BUILD_FLAGS} -o $(BUILD_DIR)/$(API_BIN)-darwin-amd64

# Build for macOS (arm64/M1)
.PHONY: build-darwin-arm64
build-darwin-arm64: check-all $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 go build ${BUILD_FLAGS} -o $(BUILD_DIR)/$(API_BIN)-darwin-arm64

# Build for Windows (amd64)
.PHONY: build-windows-amd64
build-windows-amd64: check-all $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build ${BUILD_FLAGS} -o $(BUILD_DIR)/$(API_BIN)-windows-amd64.exe

# Build all platforms
.PHONY: build-all
build-all: check-all build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64

# Build the Docker image
.PHONY: build-docker-${API_BIN}
build-docker-${API_BIN}:
	docker build -t ${API_BIN}:latest .

.PHONY: build-docker-${SUMMARIZER_WORKER_BIN}
build-docker-${SUMMARIZER_WORKER_BIN}:
	docker build -t ${SUMMARIZER_WORKER_BIN}:latest .

.PHONY: build-docker-${QUEUE_GENERATOR_BIN}
build-docker-${QUEUE_GENERATOR_BIN}:
	docker build -t ${QUEUE_GENERATOR_BIN}:latest . --build-arg MAIN_PATH=./cmd/queue-generator/main.go \
		--build-arg APP_NAME=${QUEUE_GENERATOR_BIN}

# Run the application
.PHONY: run-api
run-api:
	go run ./cmd/api/main.go

.PHONY: run-queue-generator
run-queue-generator:
	go run ./cmd/queue-generator/main.go

# Default target
.DEFAULT_GOAL := build 