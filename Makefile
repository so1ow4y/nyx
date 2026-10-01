BIN_DIR := ./bin
CLI_EXE := nyx-cli
SERVER_EXE := nyx-server

.PHONY: init build-cli build-server test clean all

init:
	@echo "Create a dirrectory $(BIN_DIR)"
	@echo ""
	@echo "Use build-cli/build-server/run-server/test..."
	@echo "run-server not build binary file. Only run file"
	@mkdir -p $(BIN_DIR)
	@echo ""

build-cli: init
	@echo "Building CLI interface..."
	@go build -o $(BIN_DIR)/$(CLI_EXE) ./cmd/cli

build-server: init
	@echo "Building server interface..."
	@go build -o $(BIN_DIR)/$(SERVER_EXE) ./cmd/server

test:
	@echo "Testing project..."
	@go test -v ./... | grep -v "no test files"

run-server: build-server
	@echo "Starting server..."
	@$(BIN_DIR)/$(SERVER_EXE)

all: build-cli build-server test

clean:
	@echo "Cleaning binaries..."
	@rm -rf $(BIN_DIR)/*
	@go clean
