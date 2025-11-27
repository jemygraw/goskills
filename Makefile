.PHONY: all build clean agent agent-web cli runner test

# Binary names
BINARY_CLI=goskills-cli
BINARY_RUNNER=goskills
BINARY_AGENT=agent-cli
BINARY_WEB=agent-web

# Build directory
BUILD_DIR=.

all: build

build: cli runner agent agent-web

agent:
	go build -o $(BUILD_DIR)/$(BINARY_AGENT) ./cmd/agent-cli

agent-web:
	go build -o $(BUILD_DIR)/$(BINARY_WEB) ./cmd/agent-web

cli:
	go build -o $(BUILD_DIR)/$(BINARY_CLI) ./cmd/skill-cli

runner:
	go build -o $(BUILD_DIR)/$(BINARY_RUNNER) ./cmd/skill-runner

clean:
	rm -f $(BUILD_DIR)/$(BINARY_CLI)
	rm -f $(BUILD_DIR)/$(BINARY_RUNNER)
	rm -f $(BUILD_DIR)/$(BINARY_AGENT)
	rm -f $(BUILD_DIR)/$(BINARY_WEB)

test:
	go test ./...