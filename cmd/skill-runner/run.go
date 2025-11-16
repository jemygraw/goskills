package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/smallnest/goskills"
	openai "github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
)

var (
	modelName string
	apiBase   string
)

var runCmd = &cobra.Command{
	Use:   "run [prompt]",
	Short: "Processes a user request by selecting and executing a skill.",
	Long: `Processes a user request by simulating the Claude skill-use workflow with an OpenAI-compatible model.
	
This command first discovers all available skills, then asks the LLM to select the most appropriate one.
Finally, it executes the selected skill by feeding its content to the LLM as a system prompt.

Requires the OPENAI_API_KEY environment variable to be set.
You can specify a custom model and API base URL using flags.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userPrompt := strings.Join(args, " ")
		skillsPath := "./examples/skills" // Hardcoded for now

		// --- PRE-FLIGHT CHECK ---
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			return errors.New("OPENAI_API_KEY environment variable is not set")
		}

		config := openai.DefaultConfig(apiKey)
		if apiBase != "" {
			config.BaseURL = apiBase
		}
		client := openai.NewClientWithConfig(config)
		ctx := context.Background()

		// --- STEP 1: SKILL DISCOVERY ---
		fmt.Println("🔎 Discovering available skills...")
		availableSkills, err := discoverSkills(skillsPath)
		if err != nil {
			return fmt.Errorf("failed to discover skills: %w", err)
		}
		if len(availableSkills) == 0 {
			return errors.New("no valid skills found")
		}
		fmt.Printf("✅ Found %d skills.\n\n", len(availableSkills))

		// --- STEP 2: SKILL SELECTION ---
		fmt.Println("🧠 Asking LLM to select the best skill...")
		selectedSkillName, err := selectSkill(ctx, client, userPrompt, availableSkills)
		if err != nil {
			return fmt.Errorf("failed during skill selection: %w", err)
		}

		selectedSkill, ok := availableSkills[selectedSkillName]
		if !ok {
			fmt.Printf("⚠️ LLM selected a non-existent skill '%s'. Aborting.\n", selectedSkillName)
			return nil
		}
		fmt.Printf("✅ LLM selected skill: %s\n\n", selectedSkillName)

		// --- STEP 3: SKILL EXECUTION ---
		fmt.Println("🚀 Executing skill...")
		fmt.Println(strings.Repeat("-", 40))

		err = executeSkill(ctx, client, userPrompt, selectedSkill)
		if err != nil {
			return fmt.Errorf("failed during skill execution: %w", err)
		}

		return nil
	},
}

func discoverSkills(skillsRoot string) (map[string]goskills.SkillPackage, error) {
	skills := make(map[string]goskills.SkillPackage)
	entries, err := os.ReadDir(skillsRoot)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		skillPath := filepath.Join(skillsRoot, entry.Name())
		skillPackage, err := goskills.ParseSkillPackage(skillPath)
		if err != nil {
			continue // Not a valid skill
		}
		skills[skillPackage.Meta.Name] = *skillPackage
	}
	return skills, nil
}

func selectSkill(ctx context.Context, client *openai.Client, userPrompt string, skills map[string]goskills.SkillPackage) (string, error) {
	var sb strings.Builder
	sb.WriteString("User Request: \"" + userPrompt + "\"\n\n")
	sb.WriteString("Available Skills:\n")
	for name, skill := range skills {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", name, skill.Meta.Description))
	}
	sb.WriteString("\nBased on the user request, which single skill is the most appropriate to use? Respond with only the name of the skill.")

	req := openai.ChatCompletionRequest{
		Model: modelName, // Use configurable model name
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an expert assistant that selects the most appropriate skill to handle a user's request. Your response must be only the exact name of the chosen skill, with no other text or explanation.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: sb.String(),
			},
		},
		Temperature: 0,
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	// Clean up the response to get only the skill name
	skillName := strings.TrimSpace(resp.Choices[0].Message.Content)
	skillName = strings.Trim(skillName, `"'`) // Trim quotes and backticks

	return skillName, nil
}

func executeSkill(ctx context.Context, client *openai.Client, userPrompt string, skill goskills.SkillPackage) error {
	// Reconstruct the skill body from structured parts
	var skillBody strings.Builder
	for _, part := range skill.Body {
		switch p := part.(type) {
		case goskills.TitlePart:
			skillBody.WriteString(fmt.Sprintf("\n# %s\n", p.Text))
		case goskills.SectionPart:
			skillBody.WriteString(fmt.Sprintf("\n## %s\n%s\n", p.Title, p.Content))
		case goskills.MarkdownPart:
			skillBody.WriteString(p.Content + "\n")
		case goskills.ImplementationPart:
			skillBody.WriteString(fmt.Sprintf("\nImplementation in %s:\n", p.Language))
			skillBody.WriteString("```" + p.Language + "\n")
			skillBody.WriteString(p.Code)
			skillBody.WriteString("```\n")
		}
	}

	req := openai.ChatCompletionRequest{
		Model: modelName, // Use configurable model name
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: skillBody.String(),
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: userPrompt,
			},
		},
		Stream: true,
	}

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return fmt.Errorf("ChatCompletionStream error: %w", err)
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println()
			return nil
		}
		if err != nil {
			return fmt.Errorf("stream error: %w", err)
		}
		fmt.Print(response.Choices[0].Delta.Content)
	}
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&modelName, "model", "m", "gpt-4o", "OpenAI-compatible model name to use")
	runCmd.Flags().StringVarP(&apiBase, "api-base", "b", "", "OpenAI-compatible API base URL (e.g., https://api.openai.com/v1)")
}