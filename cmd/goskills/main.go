package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/smallnest/goskills"
	"github.com/smallnest/goskills/log"
	goskills_mcp "github.com/smallnest/goskills/mcp"
	"github.com/spf13/cobra"
)

// Version is the version of the tool, set at build time
var Version = "v0.4.3"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goskills",
	Short: "A CLI tool for running Claude skills with LLMs.",
	Long: `goskills is a command-line interface to help you execute Claude Skill packages
using Large Language Models (LLMs) like OpenAI's GPT models.`,
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// Disable the default completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		log.Error("command execution failed: %v", err)
		os.Exit(1)
	}
}

func main() {
	rootCmd.AddCommand(runCmd)
	setupFlags(runCmd)

	Execute()
}

var runCmd = &cobra.Command{
	Use:   "run [prompt]",
	Short: "Processes a user request by selecting and executing a skill.",
	Long: `Processes a user request by simulating the Claude skill-use workflow with an OpenAI-compatible model.
	
This command first discovers all available skills, then asks the LLM to select the most appropriate one.
Finally, it executes the selected skill by feeding its content to the LLM as a system prompt.
If the LLM decides to call a tool, the tool will be executed and its output fed back to the LLM.

Requires the OPENAI_API_KEY environment variable to be set.
You can specify a custom model and API base URL using flags.`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		userPrompt := strings.Join(args, " ")
		if len(args) == 0 {
			userPromptBytes, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read from stdin: %w", err)
			}
			userPrompt = strings.TrimSpace(string(userPromptBytes))
		}

		if userPrompt == "" {
			return cmd.Help()
		}

		cfg, err := loadConfig(cmd)
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		runnerCfg := goskills.RunnerConfig{
			APIKey:           cfg.APIKey,
			APIBase:          cfg.APIBase,
			Model:            cfg.Model,
			SkillsDir:        cfg.SkillsDir,
			Verbose:          cfg.Verbose,
			AutoApproveTools: cfg.AutoApproveTools,
			AllowedScripts:   cfg.AllowedScripts,
			Loop:             cfg.Loop,
		}

		ctx := context.Background()

		// Initialize MCP Client
		var mcpClient *goskills_mcp.Client
		var mcpConfigPath string

		if cfg.McpConfig != "" {
			mcpConfigPath = cfg.McpConfig
		} else {
			// Check local mcp.json
			if _, err := os.Stat("mcp.json"); err == nil {
				mcpConfigPath = "mcp.json"
			}
			// TODO: Check ~/.claude.json if needed in future
		}

		if mcpConfigPath != "" {
			if cfg.Verbose {
				log.Info("loading mcp config from: %s", mcpConfigPath)
			}
			mcpConfig, err := goskills_mcp.LoadConfig(mcpConfigPath)
			if err != nil {
				log.Warn("failed to load mcp config: %v", err)
			} else {
				mcpClient, err = goskills_mcp.NewClient(ctx, mcpConfig)
				if err != nil {
					log.Warn("failed to create mcp client: %v", err)
				} else {
					defer mcpClient.Close()
					if cfg.Verbose {
						log.Info("mcp client initialized")
					}
				}
			}
		}

		agent, err := goskills.NewAgent(runnerCfg, mcpClient)
		if err != nil {
			return fmt.Errorf("failed to create agent: %w", err)
		}

		if runnerCfg.Loop {
			return agent.RunLoop(ctx, userPrompt)
		}

		result, err := agent.Run(ctx, userPrompt)
		if err != nil {
			return err
		}

		fmt.Println(result)
		return nil
	},
}
