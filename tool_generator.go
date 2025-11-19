package goskills

import (
	"fmt"
	"path/filepath"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"github.com/smallnest/goskills/tool"
)

// GenerateToolDefinitions generates the list of OpenAI tools for a given skill.
// It returns the tool definitions and a map of tool names to script paths for execution.
func GenerateToolDefinitions(skill SkillPackage) ([]openai.Tool, map[string]string) {
	var tools []openai.Tool
	scriptMap := make(map[string]string)

	// 1. Base Tools
	baseTools := tool.GetBaseTools()

	if len(skill.Meta.AllowedTools) > 0 {
		allowedMap := make(map[string]bool)
		for _, t := range skill.Meta.AllowedTools {
			allowedMap[t] = true
		}

		for _, t := range baseTools {
			if allowedMap[t.Function.Name] {
				tools = append(tools, t)
			}
		}
	} else {
		tools = append(tools, baseTools...)
	}

	// 2. Script Tools
	for _, scriptRelPath := range skill.Resources.Scripts {
		toolDef, toolName := generateScriptTool(skill.Path, scriptRelPath)
		tools = append(tools, toolDef)
		scriptMap[toolName] = filepath.Join(skill.Path, scriptRelPath)
	}

	return tools, scriptMap
}

func generateScriptTool(skillPath, scriptRelPath string) (openai.Tool, string) {
	// Normalize name: replace non-alphanumeric with underscore
	safeName := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, scriptRelPath)
	toolName := "run_" + safeName

	// Determine type based on extension
	ext := filepath.Ext(scriptRelPath)
	var description string
	if ext == ".py" {
		description = fmt.Sprintf("Executes the python script '%s'.", scriptRelPath)
	} else {
		description = fmt.Sprintf("Executes the shell script '%s'.", scriptRelPath)
	}

	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        toolName,
			Description: description,
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"args": map[string]interface{}{
						"type":        "array",
						"description": "Arguments to pass to the script.",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
				},
			},
		},
	}, toolName
}
