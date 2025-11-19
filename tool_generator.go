package goskills

import (
	"fmt"
	"path/filepath"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// GenerateToolDefinitions generates the list of OpenAI tools for a given skill.
// It returns the tool definitions and a map of tool names to script paths for execution.
func GenerateToolDefinitions(skill SkillPackage) ([]openai.Tool, map[string]string) {
	var tools []openai.Tool
	scriptMap := make(map[string]string)

	// 1. Base Tools
	baseTools := getBaseTools()

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

func getBaseTools() []openai.Tool {
	return []openai.Tool{
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "run_shell_script",
				Description: "Executes a shell script and returns its combined stdout and stderr. Use this for general shell commands.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"scriptPath": map[string]interface{}{
							"type":        "string",
							"description": "The path to the shell script to execute.",
						},
						"args": map[string]interface{}{
							"type":        "array",
							"description": "A list of string arguments to pass to the script.",
							"items": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"required": []string{"scriptPath"},
				},
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "run_python_script",
				Description: "Executes a Python script and returns its combined stdout and stderr.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"scriptPath": map[string]interface{}{
							"type":        "string",
							"description": "The path to the Python script to execute.",
						},
						"args": map[string]interface{}{
							"type":        "array",
							"description": "A list of string arguments to pass to the script.",
							"items": map[string]interface{}{
								"type": "string",
							},
						},
					},
					"required": []string{"scriptPath"},
				},
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "read_file",
				Description: "Reads the content of a file and returns it as a string.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"filePath": map[string]interface{}{
							"type":        "string",
							"description": "The path to the file to read.",
						},
					},
					"required": []string{"filePath"},
				},
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "write_file",
				Description: "Writes the given content to a file. If the file does not exist, it will be created. If it exists, its content will be truncated.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"filePath": map[string]interface{}{
							"type":        "string",
							"description": "The path to the file to write.",
						},
						"content": map[string]interface{}{
							"type":        "string",
							"description": "The content to write to the file.",
						},
					},
					"required": []string{"filePath", "content"},
				},
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "duckduckgo_search",
				Description: "Performs a DuckDuckGo search for the given query and returns a summary or related topics.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "The search query.",
						},
					},
					"required": []string{"query"},
				},
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "wikipedia_search",
				Description: "Performs a search on Wikipedia for the given query and returns a summary of the relevant entry.",
				Parameters: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query": map[string]interface{}{
							"type":        "string",
							"description": "The search query for Wikipedia.",
						},
					},
					"required": []string{"query"},
				},
			},
		},
	}
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
