package mcp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseToolName(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedServer string
		expectedTool   string
		expectError    bool
	}{
		{
			name:           "valid name",
			input:          "server__tool",
			expectedServer: "server",
			expectedTool:   "tool",
			expectError:    false,
		},
		{
			name:        "missing separator",
			input:       "tool",
			expectError: true,
		},
		{
			name:        "too many separators",
			input:       "server__tool__extra",
			expectError: true,
		},
		{
			name:        "empty",
			input:       "",
			expectError: true,
		},
		{
			name:           "complex server name",
			input:          "my-server__my_tool",
			expectedServer: "my-server",
			expectedTool:   "my_tool",
			expectError:    false,
		},
		{
			name:           "single character names",
			input:          "a__b",
			expectedServer: "a",
			expectedTool:   "b",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, tool, err := parseToolName(tt.input)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedServer, server)
				assert.Equal(t, tt.expectedTool, tool)
			}
		})
	}
}

// TestNewClient tests client creation
func TestNewClient(t *testing.T) {
	config := &Config{
		MCPServers: map[string]MCPServer{
			"test": {
				Type:    "stdio",
				Command: "echo",
				Args:    []string{"test"},
			},
		},
	}

	client, err := NewClient(context.Background(), config)

	// Should create client successfully (even if connection fails)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.NotNil(t, client.sessions)
}

// TestNewClient_EmptyConfig tests client creation with empty config
func TestNewClient_EmptyConfig(t *testing.T) {
	config := &Config{
		MCPServers: map[string]MCPServer{},
	}

	client, err := NewClient(context.Background(), config)

	// Should create client successfully
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Empty(t, client.sessions)
}

// TestNewClient_EmptyMCPServers tests client creation with empty MCP servers
func TestNewClient_EmptyMCPServers(t *testing.T) {
	config := &Config{
		MCPServers: map[string]MCPServer{},
	}

	client, err := NewClient(context.Background(), config)

	// Should create client successfully
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.Empty(t, client.sessions)
}

// TestClient_GetTools tests getting tools from client
func TestClient_GetTools(t *testing.T) {
	config := &Config{
		MCPServers: map[string]MCPServer{
			"test": {
				Type:    "stdio",
				Command: "echo",
				Args:    []string{"test"},
			},
		},
	}

	client, err := NewClient(context.Background(), config)
	assert.NoError(t, err)

	// Test getting tools (will succeed but return empty tools since echo is not a real MCP server)
	tools, err := client.GetTools(context.Background())

	// GetTools returns success even with failed connections
	// It logs errors and continues
	assert.NoError(t, err)
	// Should return empty tools slice since no valid MCP servers
	if tools != nil {
		assert.Equal(t, 0, len(tools)) // Check length is 0 if not nil
	}
}

// TestClient_CallTool tests calling a tool
func TestClient_CallTool(t *testing.T) {
	config := &Config{
		MCPServers: map[string]MCPServer{
			"test": {
				Type:    "stdio",
				Command: "echo",
				Args:    []string{"test"},
			},
		},
	}

	client, err := NewClient(context.Background(), config)
	assert.NoError(t, err)

	// Test calling tool (may fail due to no real server)
	result, err := client.CallTool(context.Background(), "test_tool", map[string]interface{}{})

	// Result handling depends on whether server is available
	if err != nil {
		// Expected in test environment - echo is not a real MCP server
		assert.Contains(t, err.Error(), "test_tool")
	} else {
		assert.NotNil(t, result)
	}
}

// TestNewClient_WithRealStdioCommand tests with a command that actually exists
func TestNewClient_WithRealStdioCommand(t *testing.T) {
	config := &Config{
		MCPServers: map[string]MCPServer{
			"cat-server": {
				Type:    "stdio",
				Command: "cat",
				Args:    []string{},
			},
		},
	}

	client, err := NewClient(context.Background(), config)

	// Should create client successfully even if the command isn't a real MCP server
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
