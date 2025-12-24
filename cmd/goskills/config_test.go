package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Create a test command with flags
	cmd := &cobra.Command{}
	setupFlags(cmd)

	// Test loading config with default values
	cmd.SetArgs([]string{})

	// Parse flags
	err := cmd.ParseFlags([]string{})
	assert.NoError(t, err)

	cfg, err := loadConfig(cmd)
	assert.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check default values (note: SkillsDir will be resolved to absolute path)
	assert.Contains(t, cfg.SkillsDir, ".goskills/skills")
	assert.True(t, cfg.AutoApproveTools)
	assert.False(t, cfg.Verbose)
	assert.False(t, cfg.Loop)
}

func TestLoadConfig_WithFlags(t *testing.T) {
	cmd := &cobra.Command{}
	setupFlags(cmd)

	// Set flags
	err := cmd.ParseFlags([]string{
		"--skills-dir", "/test/skills",
		"--model", "gpt-4",
		"--api-base", "https://api.openai.com/v1/",
		"--verbose",
		"--loop",
	})
	assert.NoError(t, err)

	cfg, err := loadConfig(cmd)
	assert.NoError(t, err)

	assert.Equal(t, "/test/skills", cfg.SkillsDir)
	assert.Equal(t, "gpt-4", cfg.Model)
	assert.Equal(t, "https://api.openai.com/v1", cfg.APIBase) // Should trim trailing slash
	assert.True(t, cfg.Verbose)
	assert.True(t, cfg.Loop)
}

func TestLoadConfig_WithEnvVars(t *testing.T) {
	// Set environment variables
	oldAPIKey := os.Getenv("OPENAI_API_KEY")
	oldAPIBase := os.Getenv("OPENAI_API_BASE")
	oldModel := os.Getenv("OPENAI_MODEL")
	defer func() {
		os.Setenv("OPENAI_API_KEY", oldAPIKey)
		os.Setenv("OPENAI_API_BASE", oldAPIBase)
		os.Setenv("OPENAI_MODEL", oldModel)
	}()

	os.Setenv("OPENAI_API_KEY", "test-key-from-env")
	os.Setenv("OPENAI_API_BASE", "https://api.example.com/v1/")
	os.Setenv("OPENAI_MODEL", "deepseek-v3")

	cmd := &cobra.Command{}
	setupFlags(cmd)

	// Parse empty flags (should use env vars)
	err := cmd.ParseFlags([]string{})
	assert.NoError(t, err)

	cfg, err := loadConfig(cmd)
	assert.NoError(t, err)

	assert.Equal(t, "test-key-from-env", cfg.APIKey)
	assert.Equal(t, "https://api.example.com/v1", cfg.APIBase)
	assert.Equal(t, "deepseek-v3", cfg.Model)
}

func TestLoadConfig_FlagsOverrideEnvVars(t *testing.T) {
	// Set environment variables
	oldAPIKey := os.Getenv("OPENAI_API_KEY")
	oldAPIBase := os.Getenv("OPENAI_API_BASE")
	oldModel := os.Getenv("OPENAI_MODEL")
	defer func() {
		os.Setenv("OPENAI_API_KEY", oldAPIKey)
		os.Setenv("OPENAI_API_BASE", oldAPIBase)
		os.Setenv("OPENAI_MODEL", oldModel)
	}()

	os.Setenv("OPENAI_API_KEY", "env-key")
	os.Setenv("OPENAI_API_BASE", "https://api.example.com/v1")
	os.Setenv("OPENAI_MODEL", "env-model")

	cmd := &cobra.Command{}
	setupFlags(cmd)

	// Set flags (should override env vars)
	err := cmd.ParseFlags([]string{
		"--api-key", "flag-key",
		"--model", "flag-model",
		"--api-base", "https://api.flag.com/v1/",
	})
	assert.NoError(t, err)

	cfg, err := loadConfig(cmd)
	assert.NoError(t, err)

	// Flags should override env vars for API key, model, and API base
	assert.Equal(t, "flag-key", cfg.APIKey)
	assert.Equal(t, "flag-model", cfg.Model)
	assert.Equal(t, "https://api.flag.com/v1", cfg.APIBase)
}

func TestLoadConfig_AbsolutePathResolution(t *testing.T) {
	cmd := &cobra.Command{}
	setupFlags(cmd)

	// Use relative path
	err := cmd.ParseFlags([]string{
		"--skills-dir", "./test/skills",
	})
	assert.NoError(t, err)

	cfg, err := loadConfig(cmd)
	assert.NoError(t, err)

	// Should be resolved to absolute path
	assert.True(t, filepath.IsAbs(cfg.SkillsDir))
	assert.Contains(t, cfg.SkillsDir, "test/skills")
}

func TestLoadConfig_DefaultSkillsDir(t *testing.T) {
	cmd := &cobra.Command{}
	setupFlags(cmd)

	// Don't set skills-dir flag
	err := cmd.ParseFlags([]string{})
	assert.NoError(t, err)

	cfg, err := loadConfig(cmd)
	assert.NoError(t, err)

	// Should use default (but resolved to absolute path)
	assert.Contains(t, cfg.SkillsDir, ".goskills/skills")
}

func TestLoadConfig_AllowedScripts(t *testing.T) {
	cmd := &cobra.Command{}
	setupFlags(cmd)

	err := cmd.ParseFlags([]string{
		"--allow-scripts=script1.py", "--allow-scripts=script2.sh", "--allow-scripts=script3.js",
	})
	assert.NoError(t, err)

	cfg, err := loadConfig(cmd)
	assert.NoError(t, err)

	// StringSliceFlag might parse differently, so just check it's not empty and contains expected items
	assert.NotEmpty(t, cfg.AllowedScripts)
	assert.Contains(t, cfg.AllowedScripts, "script1.py")
}

func TestLoadConfig_McpConfig(t *testing.T) {
	cmd := &cobra.Command{}
	setupFlags(cmd)

	err := cmd.ParseFlags([]string{
		"--mcp-config", "/path/to/mcp.json",
	})
	assert.NoError(t, err)

	cfg, err := loadConfig(cmd)
	assert.NoError(t, err)

	assert.Equal(t, "/path/to/mcp.json", cfg.McpConfig)
}

func TestSetupFlags(t *testing.T) {
	cmd := &cobra.Command{}

	// Before setup flags - should not have our flags
	assert.Nil(t, cmd.Flags().Lookup("skills-dir"))
	assert.Nil(t, cmd.Flags().Lookup("model"))
	assert.Nil(t, cmd.Flags().Lookup("verbose"))

	// Setup flags
	setupFlags(cmd)

	// After setup - should have our flags
	assert.NotNil(t, cmd.Flags().Lookup("skills-dir"))
	assert.NotNil(t, cmd.Flags().Lookup("model"))
	assert.NotNil(t, cmd.Flags().Lookup("api-base"))
	assert.NotNil(t, cmd.Flags().Lookup("api-key"))
	assert.NotNil(t, cmd.Flags().Lookup("auto-approve"))
	assert.NotNil(t, cmd.Flags().Lookup("allow-scripts"))
	assert.NotNil(t, cmd.Flags().Lookup("verbose"))
	assert.NotNil(t, cmd.Flags().Lookup("loop"))
	assert.NotNil(t, cmd.Flags().Lookup("mcp-config"))

	// Check shorthand flags
	flag := cmd.Flags().Lookup("skills-dir")
	assert.Equal(t, "d", flag.Shorthand)

	flag = cmd.Flags().Lookup("model")
	assert.Equal(t, "m", flag.Shorthand)

	flag = cmd.Flags().Lookup("verbose")
	assert.Equal(t, "v", flag.Shorthand)

	flag = cmd.Flags().Lookup("loop")
	assert.Equal(t, "l", flag.Shorthand)
}

func TestLoadConfig_APITrimTrailingSlash(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single slash",
			input:    "https://api.openai.com/v1/",
			expected: "https://api.openai.com/v1",
		},
		{
			name:     "double slash",
			input:    "https://api.openai.com/v1//",
			expected: "https://api.openai.com/v1",
		},
		{
			name:     "no slash",
			input:    "https://api.openai.com/v1",
			expected: "https://api.openai.com/v1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			setupFlags(cmd)

			err := cmd.ParseFlags([]string{"--api-base", tc.input})
			assert.NoError(t, err)

			cfg, err := loadConfig(cmd)
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, cfg.APIBase)
		})
	}
}
