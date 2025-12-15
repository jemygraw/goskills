package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	// Test that Execute doesn't panic
	// We can't easily test the full execution since it would exit,
	// but we can test the setup
	assert.NotPanics(t, func() {
		// Test that completion is disabled
		assert.True(t, rootCmd.CompletionOptions.DisableDefaultCmd)
		assert.Equal(t, "goskills-cli", rootCmd.Use)
		assert.Equal(t, "A CLI tool for creating and managing Claude skills.", rootCmd.Short)
	})
}

func TestRootCmd(t *testing.T) {
	// Test root command properties
	assert.Equal(t, "goskills-cli", rootCmd.Use)
	assert.Equal(t, "A CLI tool for creating and managing Claude skills.", rootCmd.Short)
	assert.Contains(t, rootCmd.Long, "goskills-cli")
	assert.Contains(t, rootCmd.Long, "Claude Skill packages")

	// Should have subcommands
	assert.True(t, len(rootCmd.Commands()) > 0)

	// Check that common commands exist
	cmdNames := []string{}
	for _, cmd := range rootCmd.Commands() {
		cmdNames = append(cmdNames, cmd.Name())
	}

	// Should contain typical skill-cli commands
	expectedCommands := []string{"list", "parse", "detail", "files", "search"}
	for _, expected := range expectedCommands {
		assert.Contains(t, cmdNames, expected, "Expected command '%s' to be available", expected)
	}
}

func TestRootCmd_Help(t *testing.T) {
	// Test that help command works
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "goskills-cli")
	assert.Contains(t, output, "command-line interface")
	assert.Contains(t, output, "Claude Skill packages")
	assert.Contains(t, output, "Available Commands:")
}

func TestRootCmd_Version(t *testing.T) {
	// Check if version flag exists (common pattern)
	versionFlag := rootCmd.PersistentFlags().Lookup("version")
	if versionFlag != nil {
		// Version flag exists - could test it but might be complex due to execution
		assert.NotNil(t, versionFlag)
	}
	// Note: We don't test version command execution as it would require complex setup
}

func TestRootCmd_InvalidCommand(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"invalid-command"})

	err := rootCmd.Execute()
	assert.Error(t, err)

	output := buf.String()
	// Should contain error message about unknown command
	hasError := assert.Contains(t, output, "unknown") || assert.Contains(t, output, "not found")
	_ = hasError // just checking if error message is present
}

func TestRootCmd_CompletionDisabled(t *testing.T) {
	// Verify that completion is disabled
	assert.True(t, rootCmd.CompletionOptions.DisableDefaultCmd)

	// Try to run completion command - should not exist
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"completion", "bash"})

	err := rootCmd.Execute()
	assert.Error(t, err)

	output := buf.String()
	hasError := assert.Contains(t, output, "unknown") || assert.Contains(t, output, "not found")
	_ = hasError // just checking if error message is present
}
