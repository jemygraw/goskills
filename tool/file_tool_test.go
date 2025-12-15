package tool

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestReadFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Test case 1: Successfully read a file
	testContent := "Hello, World!"
	testFile := filepath.Join(tmpDir, "test.txt")

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	content, err := ReadFile(testFile)
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
		return
	}

	if content != testContent {
		t.Errorf("ReadFile() = %q, want %q", content, testContent)
	}

	// Test case 2: File does not exist
	_, err = ReadFile(filepath.Join(tmpDir, "nonexistent.txt"))
	if err == nil {
		t.Error("ReadFile() expected error for nonexistent file, got nil")
	}
}

func TestWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "write_test.txt")
	testContent := "Test content for writing"

	// Test case 1: Successfully write to a new file
	err := WriteFile(testFile, testContent)
	if err != nil {
		t.Errorf("WriteFile() error = %v", err)
		return
	}

	// Verify the content was written correctly
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("WriteFile() content = %q, want %q", string(content), testContent)
	}

	// Test case 2: Overwrite an existing file
	newContent := "Overwritten content"
	err = WriteFile(testFile, newContent)
	if err != nil {
		t.Errorf("WriteFile() overwrite error = %v", err)
		return
	}

	content, err = os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read overwritten file: %v", err)
	}

	if string(content) != newContent {
		t.Errorf("WriteFile() overwrite content = %q, want %q", string(content), newContent)
	}

	// Test case 3: Write to a directory that doesn't exist (should fail)
	invalidPath := filepath.Join(tmpDir, "nonexistent", "file.txt")
	err = WriteFile(invalidPath, "test")
	if err == nil {
		t.Error("WriteFile() expected error for invalid path, got nil")
	}
}

// Test using in-memory file system for faster testing
func TestFileOperationsWithMemFS(t *testing.T) {
	memFS := fstest.MapFS{
		"test.txt": &fstest.MapFile{
			Data: []byte("In-memory file content"),
		},
	}

	// This test demonstrates how you could use in-memory filesystem for testing
	// Note: Our current implementation uses os.ReadFile directly, so this is just illustrative
	file, err := memFS.Open("test.txt")
	if err != nil {
		t.Fatalf("Failed to open file from memFS: %v", err)
	}
	defer file.Close()

	content := make([]byte, 25)
	n, err := file.Read(content)
	if err != nil {
		t.Fatalf("Failed to read from memFS: %v", err)
	}

	expected := "In-memory file content"
	if string(content[:n]) != expected {
		t.Errorf("memFS content = %q, want %q", string(content[:n]), expected)
	}
}
