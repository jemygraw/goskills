# GoSkills - Claude Skills Management Tool

English | [简体中文](README_CN.md)

A powerful command-line tool to parse, manage, and execute Claude Skill packages. GoSkills is designed according to the specifications found in the [official Claude documentation](https://docs.claude.com/en/docs/agents-and-tools/agent-skills/).

[![License](https://img.shields.io/:license-MIT-blue.svg)](https://opensource.org/licenses/MIT) [![github actions](https://github.com/smallnest/goskills/actions/workflows/go.yml/badge.svg)](https://github.com/smallnest/goskills/actions) [![Go Report Card](https://goreportcard.com/badge/github.com/smallnest/goskills)](https://goreportcard.com/report/github.com/smallnest/goskills)

## Features

- **Skill Management**: List, search, parse, and inspect Claude skills from local directories
- **Runtime Execution**: Execute skills with LLM integration (OpenAI, Claude, and compatible APIs)
- **Web Interface**: Interactive chat UI with real-time updates, session replay, and rich artifact rendering (PPT, Podcasts)
- **Rich Content Generation**: Generate PowerPoint presentations (via Slidev) and Podcast audio
- **Deep Research**: Recursive analysis and self-correction capabilities for in-depth investigation
- **Built-in Tools**: Shell commands, Python execution, file operations, web fetching, and search
- **MCP Support**: Model Context Protocol (MCP) server integration
- **Internationalization**: Full support for English and Chinese languages
- **Comprehensive Testing**: Full test suite with coverage reports

## Installation

### From Source

```shell
git clone https://github.com/smallnest/goskills.git
cd goskills
make
```

### Using Homebrew

```shell
brew install smallnest/goskills/goskills
```

Or:

```shell
# add tap
brew tap smallnest/goskills

# install goskills
brew install goskills
```

## Quick Start

```shell
# Set your OpenAI API key
export OPENAI_API_KEY="YOUR_OPENAI_API_KEY"

# Start the Web Interface
./agent-web

# List available skills
./goskills-cli list ./skills

# Run a skill using the runner
./goskills run "create a react component for a todo app"
```

## Built-in Tools

GoSkills includes a comprehensive set of built-in tools for skill execution:

- **Shell Tools**: Execute shell commands and scripts
- **Python Tools**: Run Python code and scripts
- **File Tools**: Read, write, and manage files
- **Web Tools**: Fetch and process web content
- **Search Tools**: Wikipedia and Tavily search integration
- **MCP Tools**: Integration with Model Context Protocol servers

## CLI Tools

GoSkills provides a suite of tools for different purposes:

### 1. Web Interface (`agent-web`)

A modern web-based interface for interacting with GoSkills agents.
- **Chat**: Real-time conversation with agents.
- **Artifacts**: View generated Reports, PPTs, and Podcasts directly in the browser.
- **History**: Replay and review past sessions.
- **Localization**: Toggle between English and Chinese interfaces.

### 2. Skill Management CLI (`goskills-cli`)

Located in `cmd/goskills-cli`, this tool helps you inspect and manage your local Claude skills.

#### Building `goskills-cli`

```shell
make cli
# or
go build -o goskills-cli ./cmd/goskills-cli
```

#### Available Commands

- **list**: Lists all valid skills in a given directory.
- **parse**: Parses a single skill and displays a summary of its structure.
- **detail**: Displays the full, detailed information for a single skill.
- **files**: Lists all the files that make up a skill package.
- **search**: Searches for skills by name or description.

### 3. Skill Runner CLI (`goskills`)

Located in `cmd/goskills`, this tool simulates the Claude skill-use workflow by integrating with Large Language Models (LLMs) like OpenAI's models.

#### Building `goskills` runner

```shell
make runner
# or
go build -o goskills ./cmd/goskills
```

#### Available Commands

#### download
Downloads a skill package from a GitHub directory URL to `~/.goskills/skills`.

```shell
# Download a skill from GitHub
./goskills download https://github.com/ComposioHQ/awesome-claude-skills/tree/master/meeting-insights-analyzer

# Download a skill with subdirectories
./goskills download https://github.com/ComposioHQ/awesome-claude-skills/tree/master/artifacts-builder
```

The download command will:
- Automatically create the `~/.goskills/skills` directory if it doesn't exist
- Recursively download all files and subdirectories
- Extract the skill name from the URL and use it as the target directory name
- Prevent duplicate downloads with error messages

#### run
Processes a user request by first discovering available skills, then asking an LLM to select the most appropriate one, and finally executing the selected skill.

**Requires the `OPENAI_API_KEY` environment variable to be set.**

```shell
# Example with default OpenAI model (gpt-4o)
export OPENAI_API_KEY="YOUR_OPENAI_API_KEY"
./goskills run "create an algorithm that generates abstract art"

# Example with a custom OpenAI-compatible model and API base URL using environment variables
export OPENAI_API_KEY="YOUR_OPENAI_API_KEY"
export OPENAI_API_BASE="https://qianfan.baidubce.com/v2"
export OPENAI_MODEL="deepseek-v3"
./goskills run "create an algorithm that generates abstract art"

# Example with a custom OpenAI-compatible model and API base URL using command-line flags
export OPENAI_API_KEY="YOUR_OPENAI_API_KEY"
./goskills run --model deepseek-v3 --api-base https://qianfan.baidubce.com/v2 "create an algorithm that generates abstract art"

# Example with a custom OpenAI-compatible model and API base URL using command-line flags, auto-approve without human-in-the-loop
./goskills run --auto-approve --model deepseek-v3 --api-base https://qianfan.baidubce.com/v2 "使用markitdown 工具解析网页 https://baike.baidu.com/item/%E5%AD%94%E5%AD%90/1584"

# Example with a custom OpenAI-compatible model and API base URL using command-line flags, in a loop mode and not exit automatically
./goskills run --auto-approve --model deepseek-v3 --api-base https://qianfan.baidubce.com/v2 --skills-dir=./testdata/skills "使用markitdown 工具解析网 页 https://baike.baidu.com/item/%E5%AD%94%E5%AD%90/1584" -l
```

## Development

### Make Commands

The project includes a comprehensive Makefile for development tasks:

```shell
# Help - shows all available commands
make help

# Build
make build          # Build both CLI and runner
make cli            # Build CLI only
make runner         # Build runner only

# Testing
make test           # Run all tests
make test-race      # Run tests with race detector
make test-coverage  # Run tests with coverage report
make benchmark      # Run benchmarks

# Code Quality
make check          # Run fmt-check, vet, and lint
make fmt            # Format all Go files
make vet            # Run go vet
make lint           # Run golangci-lint

# Dependencies
make deps           # Download dependencies
make tidy           # Tidy and verify dependencies

# Other
make clean          # Clean build artifacts
make install-tools  # Install development tools
make info           # Display project information
```

### Running All Tests

```shell
# Run comprehensive test suite
make test-coverage

# Run specific tool tests
cd tool && ./test_all.sh
```

## Configuration

### Environment Variables

- `OPENAI_API_KEY`: OpenAI API key for LLM integration
- `OPENAI_API_BASE`: Custom API base URL (optional)
- `OPENAI_MODEL`: Custom model name (optional)
- `TAVILY_API_KEY`: Tavily search API key
- `MCP_CONFIG`: Path to MCP configuration file

### MCP Integration

Configure Model Context Protocol (MCP) servers by creating a `mcp.json` file:

```json
{
  "mcpServers": {
    "server-name": {
      "command": "path/to/server",
      "args": ["arg1", "arg2"]
    }
  }
}
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Run linting (`make check`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
