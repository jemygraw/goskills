package tool

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPythonTool_Run(t *testing.T) {
	pythonTool := &PythonTool{}

	// Test case 1: Simple Python code
	args := map[string]any{}
	code := "print('Hello from Python!')"

	result, err := pythonTool.Run(args, code)
	if err != nil {
		t.Errorf("PythonTool.Run() error = %v", err)
		return
	}

	expected := "Hello from Python!\n"
	if result != expected {
		t.Errorf("PythonTool.Run() = %q, want %q", result, expected)
	}

	// Test case 2: Python code with template arguments
	args = map[string]any{
		"name":  "GoTest",
		"value": 42,
	}
	code = `print("Name: {{.name}}, Value: {{.value}}")`

	result, err = pythonTool.Run(args, code)
	if err != nil {
		t.Errorf("PythonTool.Run() with args error = %v", err)
		return
	}

	expected = "Name: GoTest, Value: 42\n"
	if result != expected {
		t.Errorf("PythonTool.Run() with args = %q, want %q", result, expected)
	}

	// Test case 3: Python code with syntax error
	args = map[string]any{}
	code = "print('unclosed string"

	_, err = pythonTool.Run(args, code)
	if err == nil {
		t.Error("PythonTool.Run() with syntax error expected error, got nil")
	}

	// Test case 4: Python code that writes to stderr
	args = map[string]any{}
	code = `import sys
print("This goes to stdout")
print("This goes to stderr", file=sys.stderr)`

	result, err = pythonTool.Run(args, code)
	if err != nil {
		t.Errorf("PythonTool.Run() with stderr output error = %v", err)
		return
	}

	// Result should contain both stdout and stderr
	if !containsString(result, "This goes to stdout") || !containsString(result, "This goes to stderr") {
		t.Errorf("PythonTool.Run() result should contain both stdout and stderr, got %q", result)
	}
}

func TestRunPythonScript(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test Python script
	scriptContent := `#!/usr/bin/env python3
import sys

print("Script started")
print(f"Arguments received: {len(sys.argv)}")
for i, arg in enumerate(sys.argv):
    print(f"Arg {i}: {arg}")
print("Script ended")`

	scriptPath := filepath.Join(tmpDir, "test_script.py")
	err := os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	if err != nil {
		t.Fatalf("Failed to create test Python script: %v", err)
	}

	// Test case 1: Run script without arguments
	result, err := RunPythonScript(scriptPath, nil)
	if err != nil {
		t.Errorf("RunPythonScript() error = %v", err)
		return
	}

	if !containsString(result, "Script started") || !containsString(result, "Script ended") {
		t.Errorf("RunPythonScript() result should contain start and end messages, got %q", result)
	}

	// Test case 2: Run script with arguments
	args := []string{"arg1", "arg2"}
	result, err = RunPythonScript(scriptPath, args)
	if err != nil {
		t.Errorf("RunPythonScript() with args error = %v", err)
		return
	}

	if !containsString(result, "arg1") || !containsString(result, "arg2") {
		t.Errorf("RunPythonScript() result should contain arguments, got %q", result)
	}

	// Test case 3: Non-existent script
	_, err = RunPythonScript(filepath.Join(tmpDir, "nonexistent.py"), nil)
	if err == nil {
		t.Error("RunPythonScript() with non-existent script expected error, got nil")
	}

	// Test case 4: Script with Python error
	errorScriptContent := `#!/usr/bin/env python3
print("Before error")
raise ValueError("Intentional error")
print("After error")`

	errorScriptPath := filepath.Join(tmpDir, "error_script.py")
	err = os.WriteFile(errorScriptPath, []byte(errorScriptContent), 0755)
	if err != nil {
		t.Fatalf("Failed to create error Python script: %v", err)
	}

	_, err = RunPythonScript(errorScriptPath, nil)
	if err == nil {
		t.Error("RunPythonScript() with failing script expected error, got nil")
	}
}

func TestPythonToolWithComplexCode(t *testing.T) {
	pythonTool := &PythonTool{}

	tests := []struct {
		name     string
		args     map[string]any
		code     string
		wantErr  bool
		contains string
	}{
		{
			name:     "Simple math",
			args:     map[string]any{},
			code:     "result = 2 + 3\nprint(f'Result: {result}')",
			contains: "Result: 5",
		},
		{
			name:     "String formatting with template",
			args:     map[string]any{"name": "Test", "age": 25},
			code:     `print("{{.name}} is {{.age}} years old")`,
			contains: "Test is 25 years old",
		},
		{
			name: "Import standard library",
			args: map[string]any{},
			code: `import datetime
print(f"Current year: {datetime.datetime.now().year}")`,
			contains: "Current year:",
		},
		{
			name: "JSON handling",
			args: map[string]any{},
			code: `import json
data = {"key": "value"}
print(json.dumps(data))`,
			contains: `{"key": "value"}`,
		},
		{
			name:    "Invalid template syntax",
			args:    map[string]any{},
			code:    "print('{{.unclosed')",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pythonTool.Run(tt.args, tt.code)

			if (err != nil) != tt.wantErr {
				t.Errorf("PythonTool.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.contains != "" && !containsString(result, tt.contains) {
				t.Errorf("PythonTool.Run() result = %q, want contains %q", result, tt.contains)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Example of how to benchmark Python execution
func BenchmarkPythonTool_Run(b *testing.B) {
	pythonTool := &PythonTool{}
	args := map[string]any{}
	code := "print('Benchmark test')"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := pythonTool.Run(args, code)
		if err != nil {
			b.Fatalf("PythonTool.Run() error = %v", err)
		}
	}
}
