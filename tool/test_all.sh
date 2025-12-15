#!/bin/bash

echo "Running all tests for Go Skills Tool..."
echo "====================================="

# Set colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
        exit 1
    fi
}

# Run tests with coverage
echo -e "\n${YELLOW}Running tests with coverage...${NC}"
go test -coverprofile=coverage.out ./...
print_status $? "Tests completed"

# Show coverage summary
echo -e "\n${YELLOW}Coverage Summary:${NC}"
go tool cover -func=coverage.out | tail -1

# Run benchmarks
echo -e "\n${YELLOW}Running benchmarks...${NC}"
go test -bench=. -benchmem ./... > /dev/null 2>&1
print_status $? "Benchmarks completed"

# Check for race conditions
echo -e "\n${YELLOW}Running race condition tests...${NC}"
go test -race ./... > /dev/null 2>&1
print_status $? "Race condition tests completed"

# Vet code
echo -e "\n${YELLOW}Running go vet...${NC}"
go vet ./...
print_status $? "Code vet completed"

# Format check
echo -e "\n${YELLOW}Checking code formatting...${NC}"
UNFORMATTED=$(gofmt -l -s .)
if [ -z "$UNFORMATTED" ]; then
    print_status 0 "Code is properly formatted"
else
    echo -e "${RED}Unformatted files:${NC}"
    echo "$UNFORMATTED"
    print_status 1 "Code formatting check failed"
fi

# Build check
echo -e "\n${YELLOW}Building project...${NC}"
go build -o /tmp/test-tool .
print_status $? "Build completed"

# Clean up
rm -f /tmp/test-tool coverage.out coverage.html

echo -e "\n${GREEN}All tests passed successfully!${NC}"