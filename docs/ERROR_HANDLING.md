# Simple Tool Error Handling

## Overview

This system provides a simple but effective approach to tool error handling by feeding error information back to the LLM, allowing it to intelligently decide how to recover from failures.

## How It Works

When a tool execution fails, the system:

1. **Logs the error** for debugging purposes
2. **Creates a detailed error message** including:
   - Tool name that failed
   - Error details
   - Original arguments
   - Suggestions for recovery
3. **Adds the error message** to the conversation history
4. **Lets the LLM decide** the next action

## Error Message Format

When a tool fails, the LLM receives a message like:

```
Tool execution failed: read_file
Error details: file not found: /path/to/nonexistent/file.txt
Tool name: read_file
Arguments: {"path": "/path/to/nonexistent/file.txt"}

You can try:
1. Retry with different parameters
2. Use a different tool
3. Modify your approach
```

## LLM Recovery Strategies

The LLM can intelligently respond to errors by:

### 1. **Parameter Correction**
```
LLM: The file path was incorrect. Let me try a different path:
{
  "tool": "read_file",
  "arguments": {"path": "./existing-file.txt"}
}
```

### 2. **Alternative Tools**
```
LLM: Reading the file failed. Let me try searching for it:
{
  "tool": "run_shell_code",
  "arguments": {"code": "find . -name '*.txt' | head -10"}
}
```

### 3. **Approach Modification**
```
LLM: The file doesn't exist. Let me create it instead:
{
  "tool": "write_file",
  "arguments": {"path": "/path/to/file.txt", "content": "new content"}
}
```

## Benefits

### ✅ **Simple & Reliable**
- No complex retry logic
- No additional configuration needed
- Works with all existing tools

### ✅ **LLM-Powered**
- Leverages LLM's reasoning abilities
- Context-aware error handling
- Flexible recovery strategies

### ✅ **Transparent**
- Full error information available to LLM
- Detailed logging for developers
- Clear recovery suggestions

### ✅ **Extensible**
- Works with any tool
- No tool-specific code needed
- Future-proof for new error types

## Implementation Details

The error handling is implemented in `runner.go` lines 349-358:

```go
if err != nil {
    log.Error("tool call failed: %v", err)
    // Provide detailed error information to help LLM understand what went wrong
    errorMsg := fmt.Sprintf("Tool execution failed: %s\nError details: %v\nTool name: %s\nArguments: %s\n\nYou can try:\n1. Retry with different parameters\n2. Use a different tool\n3. Modify your approach",
        tc.Function.Name, err, tc.Function.Name, tc.Function.Arguments)
    a.messages = append(a.messages, openai.ChatCompletionMessage{
        Role:       openai.ChatMessageRoleTool,
        ToolCallID: tc.ID,
        Content:    errorMsg,
    })
}
```

## Example Scenario

**User**: "Read the configuration file and show me the API key"

**LLM**: I'll read the config file to get the API key.
```json
{
  "tool": "read_file",
  "arguments": {"path": "/etc/config.json"}
}
```

**System**: Tool fails with error
```
Tool execution failed: read_file
Error details: permission denied: open /etc/config.json
Tool name: read_file
Arguments: {"path": "/etc/config.json"}

You can try:
1. Retry with different parameters
2. Use a different tool
3. Modify your approach
```

**LLM**: I don't have permission to read that file. Let me try a different approach and check for user config files:
```json
{
  "tool": "read_file",
  "arguments": {"path": "./config.json"}
}
```

This simple approach allows the LLM to intelligently handle errors without complex retry logic or failure classification systems.