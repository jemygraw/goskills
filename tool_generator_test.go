package goskills

import (
	"testing"

	openai "github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

// TestGenerateToolDefinitions_AllowedTools tests tool generation with allowed tools filter
func TestGenerateToolDefinitions_AllowedTools(t *testing.T) {
	skill := SkillPackage{
		Path: "/test/skill",
		Meta: SkillMeta{
			AllowedTools: []string{"read_file", "write_file"}, // Only allow these tools
		},
		Resources: SkillResources{
			Scripts: []string{"test.py", "setup.sh"},
		},
	}

	tools, scriptMap := GenerateToolDefinitions(&skill)

	// Should have 2 allowed base tools + 2 script tools
	assert.Len(t, tools, 4)
	assert.Len(t, scriptMap, 2)

	// Check that only allowed tools are included
	toolNames := make(map[string]bool)
	for _, tool := range tools {
		toolNames[tool.Function.Name] = true
	}
	assert.True(t, toolNames["read_file"])
	assert.True(t, toolNames["write_file"])
	assert.True(t, toolNames["run_test_py"])
	assert.True(t, toolNames["run_setup_sh"])
}

// TestGenerateToolDefinitions_NoAllowedTools tests tool generation without allowed tools filter
func TestGenerateToolDefinitions_NoAllowedTools(t *testing.T) {
	skill := SkillPackage{
		Path: "/test/skill",
		Meta: SkillMeta{
			AllowedTools: []string{}, // Empty allowed tools
		},
		Resources: SkillResources{
			Scripts: []string{"deploy.sh"},
		},
	}

	tools, scriptMap := GenerateToolDefinitions(&skill)

	// Should have all base tools + 1 script tool
	assert.Greater(t, len(tools), 1) // At least one script tool
	assert.Len(t, scriptMap, 1)      // One script tool

	// Check script map contains correct path
	assert.Contains(t, scriptMap, "run_deploy_sh")
	assert.Equal(t, "/test/skill/deploy.sh", scriptMap["run_deploy_sh"])
}

// TestGenerateToolDefinitions_NoScripts tests tool generation with no scripts
func TestGenerateToolDefinitions_NoScripts(t *testing.T) {
	skill := SkillPackage{
		Path: "/test/skill",
		Meta: SkillMeta{
			AllowedTools: []string{"read_file"},
		},
		Resources: SkillResources{
			Scripts: []string{}, // No scripts
		},
	}

	tools, scriptMap := GenerateToolDefinitions(&skill)

	// Should have only the allowed base tool
	assert.Len(t, tools, 1)
	assert.Len(t, scriptMap, 0) // No script tools

	// Check it's the correct tool
	assert.Equal(t, "read_file", tools[0].Function.Name)
}

// TestGenerateScriptTool_PythonScript tests Python script tool generation
func TestGenerateScriptTool_PythonScript(t *testing.T) {
	skillPath := "/test/skill"
	scriptRelPath := "scripts/test.py"

	tool, toolName := generateScriptTool(skillPath, scriptRelPath)

	// Check tool name generation
	assert.Equal(t, "run_scripts_test_py", toolName)

	// Check tool definition
	assert.Equal(t, openai.ToolTypeFunction, tool.Type)
	assert.Equal(t, toolName, tool.Function.Name)
	assert.Contains(t, tool.Function.Description, "python script")
	assert.Contains(t, tool.Function.Description, "scripts/test.py")

	// Verify tool structure
	assert.NotNil(t, tool.Function)
	assert.NotNil(t, tool.Function.Parameters)
}

// TestGenerateScriptTool_ShellScript tests shell script tool generation
func TestGenerateScriptTool_ShellScript(t *testing.T) {
	skillPath := "/test/skill"
	scriptRelPath := "deploy.sh"

	tool, toolName := generateScriptTool(skillPath, scriptRelPath)

	// Check tool name generation
	assert.Equal(t, "run_deploy_sh", toolName)

	// Check tool definition
	assert.Equal(t, openai.ToolTypeFunction, tool.Type)
	assert.Equal(t, toolName, tool.Function.Name)
	assert.Contains(t, tool.Function.Description, "shell script")
	assert.Contains(t, tool.Function.Description, "deploy.sh")

	// Verify tool structure
	assert.NotNil(t, tool.Function)
	assert.NotNil(t, tool.Function.Parameters)
}

// TestGenerateScriptTool_SpecialCharacters tests script tool generation with special characters
func TestGenerateScriptTool_SpecialCharacters(t *testing.T) {
	testCases := []struct {
		scriptRelPath string
		expectedName  string
	}{
		{
			scriptRelPath: "my-script.sh",
			expectedName:  "run_my_script_sh",
		},
		{
			scriptRelPath: "test script.py",
			expectedName:  "run_test_script_py",
		},
		{
			scriptRelPath: "setup-v1.0.sh",
			expectedName:  "run_setup_v1_0_sh",
		},
		{
			scriptRelPath: "data@process.py",
			expectedName:  "run_data_process_py",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.scriptRelPath, func(t *testing.T) {
			tool, toolName := generateScriptTool("/test", tc.scriptRelPath)
			assert.Equal(t, tc.expectedName, toolName)
			assert.NotNil(t, tool.Function)
		})
	}
}

// TestGenerateToolDefinitions_EmptySkill tests tool generation with minimal skill
func TestGenerateToolDefinitions_EmptySkill(t *testing.T) {
	skill := SkillPackage{
		Path: "/test/skill",
		Meta: SkillMeta{}, // No allowed tools specified
		Resources: SkillResources{
			Scripts: []string{}, // No scripts
		},
	}

	tools, scriptMap := GenerateToolDefinitions(&skill)

	// Should have all base tools
	assert.Greater(t, len(tools), 0)
	assert.Len(t, scriptMap, 0)

	// Verify tool structure
	for _, tool := range tools {
		assert.Equal(t, openai.ToolTypeFunction, tool.Type)
		assert.NotNil(t, tool.Function)
		assert.NotEmpty(t, tool.Function.Name)
	}
}

// TestGenerateScriptTool_ParametersStructure tests that parameters structure is correct
func TestGenerateScriptTool_ParametersStructure(t *testing.T) {
	skillPath := "/test/skill"
	scriptRelPath := "test.py"

	tool, _ := generateScriptTool(skillPath, scriptRelPath)

	// Check that parameters is a map
	assert.IsType(t, map[string]interface{}{}, tool.Function.Parameters)

	params := tool.Function.Parameters.(map[string]interface{})
	assert.Equal(t, "object", params["type"])
	assert.Contains(t, params, "properties")

	// Check properties structure
	properties := params["properties"].(map[string]interface{})
	assert.Contains(t, properties, "args")

	args := properties["args"].(map[string]interface{})
	assert.Equal(t, "array", args["type"])
	assert.Equal(t, "Arguments to pass to the script.", args["description"])
}
