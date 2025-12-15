package tool

import (
	"os"
	"path/filepath"
	"testing"
)

func TestShellTool_Run(t *testing.T) {
	shellTool := &ShellTool{}

	// Test case 1: Simple command with no arguments
	args := map[string]any{}
	code := "echo 'Hello World'"

	result, err := shellTool.Run(args, code)
	if err != nil {
		t.Errorf("ShellTool.Run() error = %v", err)
		return
	}

	expected := "Hello World\n"
	if result != expected {
		t.Errorf("ShellTool.Run() = %q, want %q", result, expected)
	}

	// Test case 2: Command with template arguments
	args = map[string]any{
		"name":  "GoTest",
		"count": 3,
	}
	code = `echo "Hello {{.name}}! Count is {{.count}}"`

	result, err = shellTool.Run(args, code)
	if err != nil {
		t.Errorf("ShellTool.Run() with args error = %v", err)
		return
	}

	expected = "Hello GoTest! Count is 3\n"
	if result != expected {
		t.Errorf("ShellTool.Run() with args = %q, want %q", result, expected)
	}

	// Test case 3: Invalid shell template
	args = map[string]any{}
	code = "echo {{.invalid.property}}"

	_, err = shellTool.Run(args, code)
	if err == nil {
		t.Error("ShellTool.Run() with invalid template expected error, got nil")
	}
}

func TestRunShellScript(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test shell script
	scriptContent := `#!/bin/bash
echo "Script started"
echo "Argument 1: $1"
echo "Argument 2: $2"
echo "Script ended"`

	scriptPath := filepath.Join(tmpDir, "test_script.sh")
	err := os.WriteFile(scriptPath, []byte(scriptContent), 0755)
	if err != nil {
		t.Fatalf("Failed to create test script: %v", err)
	}

	// Test case 1: Run script without arguments
	result, err := RunShellScript(scriptPath, nil)
	if err != nil {
		t.Errorf("RunShellScript() error = %v", err)
		return
	}

	expected := "Script started\nArgument 1: \nArgument 2: \nScript ended\n"
	if result != expected {
		t.Errorf("RunShellScript() = %q, want %q", result, expected)
	}

	// Test case 2: Run script with arguments
	args := []string{"arg1", "arg2"}
	result, err = RunShellScript(scriptPath, args)
	if err != nil {
		t.Errorf("RunShellScript() with args error = %v", err)
		return
	}

	expected = "Script started\nArgument 1: arg1\nArgument 2: arg2\nScript ended\n"
	if result != expected {
		t.Errorf("RunShellScript() with args = %q, want %q", result, expected)
	}

	// Test case 3: Non-existent script
	_, err = RunShellScript(filepath.Join(tmpDir, "nonexistent.sh"), nil)
	if err == nil {
		t.Error("RunShellScript() with non-existent script expected error, got nil")
	}

	// Test case 4: Script with error
	errorScriptContent := `#!/bin/bash
echo "Before error"
exit 1
echo "After error"`

	errorScriptPath := filepath.Join(tmpDir, "error_script.sh")
	err = os.WriteFile(errorScriptPath, []byte(errorScriptContent), 0755)
	if err != nil {
		t.Fatalf("Failed to create error script: %v", err)
	}

	_, err = RunShellScript(errorScriptPath, nil)
	if err == nil {
		t.Error("RunShellScript() with failing script expected error, got nil")
	}
}

func TestShellToolWithComplexCommands(t *testing.T) {
	shellTool := &ShellTool{}

	tests := []struct {
		name     string
		args     map[string]any
		code     string
		wantErr  bool
		contains string
	}{
		{
			name:     "List files",
			args:     map[string]any{},
			code:     "ls /tmp",
			wantErr:  false,
			contains: "",
		},
		{
			name:     "Environment variable",
			args:     map[string]any{"path": "/tmp"},
			code:     "echo 'Directory: {{.path}}'",
			wantErr:  false,
			contains: "Directory: /tmp\n",
		},
		{
			name:     "Pipeline command",
			args:     map[string]any{},
			code:     "echo 'hello world' | tr '[:lower:]' '[:upper:]'",
			wantErr:  false,
			contains: "HELLO WORLD\n",
		},
		{
			name:     "Empty command",
			args:     map[string]any{},
			code:     "",
			wantErr:  false,
			contains: "",
		},
		{
			name:    "Invalid template syntax",
			args:    map[string]any{},
			code:    "echo {{.unclosed",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := shellTool.Run(tt.args, tt.code)

			if (err != nil) != tt.wantErr {
				t.Errorf("ShellTool.Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.contains != "" && result != tt.contains {
				t.Errorf("ShellTool.Run() result = %q, want contains %q", result, tt.contains)
			}
		})
	}
}

// Example of how to benchmark shell execution
func BenchmarkShellTool_Run(b *testing.B) {
	shellTool := &ShellTool{}
	args := map[string]any{}
	code := "echo 'Benchmark test'"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := shellTool.Run(args, code)
		if err != nil {
			b.Fatalf("ShellTool.Run() error = %v", err)
		}
	}
}
