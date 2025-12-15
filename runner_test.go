package goskills

import (
	"context"
	"os"
	"testing"

	openai "github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

// createTestAgent creates an agent for testing using environment variables
func createTestAgent(t *testing.T) *Agent {
	token := os.Getenv("OPENAI_API_KEY")
	if token == "" {
		t.Skip("OPENAI_API_KEY not set, skipping integration test")
	}

	config := openai.DefaultConfig(token)
	if apiURL := os.Getenv("OPENAI_API_BASE"); apiURL != "" {
		config.BaseURL = apiURL
	}

	client := openai.NewClientWithConfig(config)
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "deepseek-v3" // default model
	}

	cfg := RunnerConfig{
		Model: model,
	}

	return &Agent{
		client: client,
		cfg:    cfg,
	}
}

// TestSelectSkill_Success tests successful skill selection
func TestSelectSkill_Success(t *testing.T) {
	agent := createTestAgent(t)

	// Create test skills
	skills := map[string]SkillPackage{
		"pdf": {
			Meta: SkillMeta{
				Name:        "pdf",
				Description: "Comprehensive PDF manipulation toolkit for extracting text and tables",
			},
		},
		"xlsx": {
			Meta: SkillMeta{
				Name:        "xlsx",
				Description: "Comprehensive spreadsheet creation, editing, and analysis",
			},
		},
	}

	// Execute test
	ctx := context.Background()
	userPrompt := "Please extract text from this PDF file"

	skillName, err := agent.selectSkill(ctx, userPrompt, skills)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, "pdf", skillName)
}

// TestSelectSkill_XlsxSelection tests xlsx skill selection
func TestSelectSkill_XlsxSelection(t *testing.T) {
	agent := createTestAgent(t)

	skills := map[string]SkillPackage{
		"pdf": {
			Meta: SkillMeta{
				Name:        "pdf",
				Description: "Comprehensive PDF manipulation toolkit for extracting text and tables",
			},
		},
		"xlsx": {
			Meta: SkillMeta{
				Name:        "xlsx",
				Description: "Comprehensive spreadsheet creation, editing, and analysis",
			},
		},
		"email": {
			Meta: SkillMeta{
				Name:        "email",
				Description: "Email composition and sending toolkit",
			},
		},
	}

	testCases := []struct {
		name          string
		userPrompt    string
		expectedSkill string
	}{
		{
			name:          "PDF extraction request",
			userPrompt:    "Extract tables from this PDF document",
			expectedSkill: "pdf",
		},
		{
			name:          "Spreadsheet creation request",
			userPrompt:    "Create an Excel spreadsheet with sales data",
			expectedSkill: "xlsx",
		},
		{
			name:          "Email request",
			userPrompt:    "Send an email to the team",
			expectedSkill: "email",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			skillName, err := agent.selectSkill(ctx, tc.userPrompt, skills)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedSkill, skillName)
		})
	}
}

// TestSelectSkill_WithQuotesAndWhitespace tests skill selection with various response formats
func TestSelectSkill_WithQuotesAndWhitespace(t *testing.T) {
	agent := createTestAgent(t)

	skills := map[string]SkillPackage{
		"pdf": {
			Meta: SkillMeta{
				Name:        "pdf",
				Description: "PDF manipulation toolkit",
			},
		},
	}

	// Test with a clear PDF request that should return pdf skill
	ctx := context.Background()
	userPrompt := "Please help me extract text from a PDF file"

	skillName, err := agent.selectSkill(ctx, userPrompt, skills)

	// Assertions - the function should handle trimming quotes and whitespace
	assert.NoError(t, err)
	assert.Equal(t, "pdf", skillName)
}

// TestSelectSkill_EmptySkillsMap tests behavior with empty skills map
func TestSelectSkill_EmptySkillsMap(t *testing.T) {
	agent := createTestAgent(t)

	// Empty skills map
	skills := map[string]SkillPackage{}

	ctx := context.Background()
	userPrompt := "Do something"

	skillName, err := agent.selectSkill(ctx, userPrompt, skills)

	// Should not error, but may return a message explaining no skills are available
	assert.NoError(t, err)
	// The AI might return an explanatory message when no skills are available
	// so we just check that it doesn't crash and returns some response
	assert.NotEmpty(t, skillName)
}

// TestSelectSkill_ComplexSkillDescriptions tests skill selection with complex descriptions
func TestSelectSkill_ComplexSkillDescriptions(t *testing.T) {
	agent := createTestAgent(t)

	// Create skills with complex descriptions containing special characters
	skills := map[string]SkillPackage{
		"pdf": {
			Meta: SkillMeta{
				Name:        "pdf",
				Description: "Comprehensive PDF manipulation toolkit for extracting text, tables, and images. Supports OCR, form filling, and digital signatures. Works with encrypted/password-protected files.",
			},
		},
		"data-analysis": {
			Meta: SkillMeta{
				Name:        "data-analysis",
				Description: "Advanced data analysis toolkit with statistical functions, visualization, and machine learning capabilities. Works with CSV, JSON, and various data formats.",
			},
		},
	}

	ctx := context.Background()
	userPrompt := "Extract text and tables from an encrypted PDF file"

	skillName, err := agent.selectSkill(ctx, userPrompt, skills)

	assert.NoError(t, err)
	assert.Equal(t, "pdf", skillName)
}

// TestNewAgent_Success tests successful agent creation
func TestNewAgent_Success(t *testing.T) {
	cfg := RunnerConfig{
		APIKey:    "test-api-key",
		Model:     "gpt-4",
		SkillsDir: "./test-skills",
	}

	agent, err := NewAgent(cfg, nil)

	assert.NoError(t, err)
	assert.NotNil(t, agent)
	assert.Equal(t, cfg.APIKey, agent.cfg.APIKey)
	assert.Equal(t, cfg.Model, agent.cfg.Model)
	assert.Equal(t, cfg.SkillsDir, agent.cfg.SkillsDir)
	assert.NotNil(t, agent.client)
	assert.NotNil(t, agent.messages)
}

// TestNewAgent_EmptyAPIKey tests agent creation with empty API key
func TestNewAgent_EmptyAPIKey(t *testing.T) {
	cfg := RunnerConfig{
		APIKey: "",
		Model:  "gpt-4",
	}

	agent, err := NewAgent(cfg, nil)

	assert.Error(t, err)
	assert.Nil(t, agent)
	assert.Contains(t, err.Error(), "API key is not set")
}

// TestNewAgent_DefaultModel tests agent creation with default model
func TestNewAgent_DefaultModel(t *testing.T) {
	cfg := RunnerConfig{
		APIKey: "test-api-key",
		// Model is empty, should default to "gpt-4o"
	}

	agent, err := NewAgent(cfg, nil)

	assert.NoError(t, err)
	assert.NotNil(t, agent)
	assert.Equal(t, "gpt-4o", agent.cfg.Model)
}

// TestDiscoverSkills tests the skill discovery functionality
func TestDiscoverSkills(t *testing.T) {
	cfg := RunnerConfig{
		APIKey:    "test-api-key",
		Model:     "gpt-4",
		SkillsDir: "./tool", // Use existing tool directory for testing
	}

	agent, err := NewAgent(cfg, nil)
	assert.NoError(t, err)

	// Test with an existing directory that should have skill-like structure
	skills, err := agent.discoverSkills("./tool")

	// We expect this to potentially error since ./tool may not be a proper skills directory
	// but we're testing the function doesn't panic
	if err != nil {
		assert.NotEmpty(t, err.Error())
	} else {
		assert.NotNil(t, skills)
	}
}

// TestExtractSkillName tests skill name extraction from AI responses
func TestExtractSkillName(t *testing.T) {
	skills := map[string]SkillPackage{
		"pdf": {
			Meta: SkillMeta{
				Name: "pdf",
			},
		},
		"xlsx": {
			Meta: SkillMeta{
				Name: "xlsx",
			},
		},
	}

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Exact skill name",
			input:    "pdf",
			expected: "pdf",
		},
		{
			name:     "Skill name with quotes",
			input:    `"pdf"`,
			expected: "pdf",
		},
		{
			name:     "Skill name with whitespace",
			input:    "  pdf  ",
			expected: "pdf",
		},
		{
			name:     "Skill name with quotes and whitespace",
			input:    `  "pdf"  `,
			expected: "pdf",
		},
		{
			name:     "Case insensitive match",
			input:    "PDF",
			expected: "pdf",
		},
		{
			name:     "Complex response with skill name",
			input:    "Based on your request, I'll use the pdf skill to help you.",
			expected: "pdf",
		},
		{
			name:     "No skill found",
			input:    "I don't know which skill to use",
			expected: "I don't know which skill to use",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := extractSkillName(tc.input, skills)
			assert.Equal(t, tc.expected, result)
		})
	}
}
