# Go Skills Tool

A collection of Go tools for various operations including file manipulation, shell execution, Python scripting, web scraping, and API integrations.

## Features

- **File Operations**: Read and write files with error handling
- **Shell Tools**: Execute shell scripts and commands with template support
- **Python Tools**: Run Python code and scripts with template support
- **Web Tools**: Fetch and parse web pages, extracting readable content
- **Search Tools**:
  - Wikipedia search integration
  - Tavily search API integration for web searches
- **OpenAI Tool Definitions**: Pre-defined tool schemas for AI integration

## Installation

```bash
go get github.com/chaoyuepan/goskills/tool
```

## Usage

### File Operations

```go
import "github.com/chaoyuepan/goskills/tool"

// Read a file
content, err := tool.ReadFile("/path/to/file.txt")
if err != nil {
    log.Fatal(err)
}

// Write a file
err = tool.WriteFile("/path/to/output.txt", "Hello, World!")
if err != nil {
    log.Fatal(err)
}
```

### Shell Tools

```go
// Create shell tool
shellTool := &tool.ShellTool{}

// Execute shell code with templates
args := map[string]any{
    "name": "John",
    "age": 30,
}
result, err := shellTool.Run(args, "echo 'Hello {{.name}}, you are {{.age}} years old'")
if err != nil {
    log.Fatal(err)
}
fmt.Println(result)

// Execute shell script
result, err := tool.RunShellScript("/path/to/script.sh", []string{"arg1", "arg2"})
if err != nil {
    log.Fatal(err)
}
fmt.Println(result)
```

### Python Tools

```go
// Create Python tool
pythonTool := &tool.PythonTool{}

// Execute Python code with templates
args := map[string]any{
    "value": 42,
}
result, err := pythonTool.Run(args, "print('The answer is {{.value}}')")
if err != nil {
    log.Fatal(err)
}
fmt.Println(result)

// Execute Python script
result, err := tool.RunPythonScript("/path/to/script.py", []string{"arg1", "arg2"})
if err != nil {
    log.Fatal(err)
}
fmt.Println(result)
```

### Web Tools

```go
// Fetch web page content
content, err := tool.WebFetch("https://example.com")
if err != nil {
    log.Fatal(err)
}
fmt.Println(content)
```

### Search Tools

```go
// Wikipedia search
result, err := tool.WikipediaSearch("Albert Einstein")
if err != nil {
    log.Fatal(err)
}
fmt.Println(result)

// Tavily search (requires TAVILY_API_KEY environment variable)
result, err := tool.TavilySearch("latest Go programming news")
if err != nil {
    log.Fatal(err)
}
fmt.Println(result)
```

### OpenAI Tool Definitions

```go
// Get all base tools for OpenAI integration
tools := tool.GetBaseTools()
for _, tool := range tools {
    fmt.Printf("Tool: %s\n", tool.Function.Name)
    fmt.Printf("Description: %s\n", tool.Function.Description)
}
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run only short tests
make test-short

# Run tests with coverage
make coverage

# Run benchmarks
make benchmark
```

### Test Coverage

The project includes comprehensive unit tests for all major components:

- File operations with in-memory filesystem testing
- Shell command execution with temporary scripts
- Python script execution with error handling
- Web fetching with mock HTTP servers
- Search API integration testing
- Tool definition validation

### Running Specific Tests

```bash
# Run tests for a specific package
go test -v ./... -run TestReadFile

# Run tests for a specific function
make test-file FILE=TestShellTool

# Run tests with race detection
go test -race -v ./...
```

## Environment Variables

- `TAVILY_API_KEY`: Required for Tavily search functionality
- `PYTHON_PATH`: Optional path to Python executable (defaults to python3 or python)

## Dependencies

- `github.com/sashabaranov/go-openai`: OpenAI API integration
- `github.com/PuerkitoBio/goquery`: HTML parsing and manipulation

## Development

### Prerequisites

- Go 1.21 or later

### Development Commands

```bash
# Install dependencies
make deps

# Format code
make fmt

# Run linter
make lint

# Run security check
make security

# Watch for changes and run tests (requires entr)
make watch
```

### Project Structure

```
tool/
├── definitions.go          # OpenAI tool definitions
├── file_tool.go           # File operations
├── file_tool_test.go      # File operations tests
├── shell_tool.go          # Shell execution
├── shell_tool_test.go     # Shell execution tests
├── python_tool.go         # Python execution
├── python_tool_test.go    # Python execution tests
├── web_tool.go            # Web fetching
├── web_tool_test.go       # Web fetching tests
├── tavily_tool.go         # Tavily search
├── tavily_tool_test.go    # Tavily search tests
├── knowledge_tool.go      # Wikipedia search
├── knowledge_tool_test.go # Wikipedia search tests
├── definitions_test.go    # Tool definitions tests
├── Makefile               # Build and test commands
├── go.mod                 # Go module file
└── README.md              # This file
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Write tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## License

This project is licensed under the MIT License.