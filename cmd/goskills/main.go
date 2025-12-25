package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/smallnest/goskills"
	"github.com/smallnest/goskills/log"
	goskills_mcp "github.com/smallnest/goskills/mcp"
	"github.com/spf13/cobra"
)

// Version is the version of the tool, set at build time
var Version = "v0.5.3"

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

	rootCmd.AddCommand(downloadCmd)

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
			Debug:            cfg.Debug,
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

var forceDownload bool

var downloadCmd = &cobra.Command{
	Use:   "download <github_url>",
	Short: "Downloads a skill from GitHub to ~/.goskills/skills",
	Long: `Downloads a skill package from a GitHub URL.
Examples:
  goskills download https://github.com/owner/repo
  goskills download https://github.com/ComposioHQ/awesome-claude-skills/tree/master/meeting-insights-analyzer`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		githubURL := args[0]

		// Parse GitHub URL to get owner, repo, branch, and path
		owner, repo, branch, dirPath, err := parseGitHubURL(githubURL)
		if err != nil {
			return fmt.Errorf("failed to parse GitHub URL: %w", err)
		}

		// Get home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		// Create ~/.goskills/skills directory if it doesn't exist
		skillsDir := filepath.Join(homeDir, ".goskills", "skills")
		if err := os.MkdirAll(skillsDir, 0755); err != nil {
			return fmt.Errorf("failed to create skills directory: %w", err)
		}

		// Extract skill name from directory path
		skillName := filepath.Base(dirPath)
		if skillName == "." || skillName == "" {
			skillName = repo
		}
		targetDir := filepath.Join(skillsDir, skillName)

		// Check if skill already exists
		if _, err := os.Stat(targetDir); err == nil {
			if forceDownload {
				log.Info("Skill '%s' already exists, removing due to -f flag...", skillName)
				if err := os.RemoveAll(targetDir); err != nil {
					return fmt.Errorf("failed to remove existing directory: %w", err)
				}
			} else {
				return fmt.Errorf("skill '%s' already exists in %s (use -f to force overwrite)", skillName, targetDir)
			}
		}

		log.Info("Downloading skill '%s' from GitHub...", skillName)

		// Download files from GitHub
		if err := downloadGitHubDirectory(owner, repo, branch, dirPath, targetDir); err != nil {
			return fmt.Errorf("failed to download skill: %w", err)
		}

		log.Info("Successfully downloaded skill to: %s", targetDir)
		return nil
	},
}

func init() {
	downloadCmd.Flags().BoolVarP(&forceDownload, "force", "f", false, "Force remove existing directory before downloading")
}

// parseGitHubURL parses a GitHub URL and extracts owner, repo, branch, and directory path
// Supports formats:
// - https://github.com/{owner}/{repo} (defaults to master branch, root path)
// - https://github.com/{owner}/{repo}/tree/{branch}/{path}
func parseGitHubURL(url string) (owner, repo, branch, path string, err error) {
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "github.com/")

	parts := strings.Split(url, "/")

	// Handle simple repo URL: owner/repo
	if len(parts) == 2 {
		owner = parts[0]
		repo = parts[1]
		branch = "master"
		path = ""
		return owner, repo, branch, path, nil
	}

	// Handle full URL with tree: owner/repo/tree/branch/path
	if len(parts) < 5 || parts[2] != "tree" {
		return "", "", "", "", fmt.Errorf("invalid GitHub URL format. Expected: https://github.com/owner/repo or https://github.com/owner/repo/tree/branch/path")
	}

	owner = parts[0]
	repo = parts[1]
	branch = parts[3]
	path = strings.Join(parts[4:], "/")

	return owner, repo, branch, path, nil
}

// GitHubContent represents a file or directory from GitHub API
type GitHubContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
	URL         string `json:"url"`
}

// downloadGitHubDirectory downloads all files from a GitHub directory recursively
func downloadGitHubDirectory(owner, repo, branch, dirPath, targetDir string) error {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", owner, repo, dirPath, branch)

	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to fetch directory contents: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var contents []GitHubContent
	if err := json.NewDecoder(resp.Body).Decode(&contents); err != nil {
		return fmt.Errorf("failed to decode GitHub API response: %w", err)
	}

	// Create target directory
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
	}

	for _, item := range contents {
		itemPath := filepath.Join(targetDir, item.Name)

		if item.Type == "file" {
			log.Info("Downloading file: %s", item.Name)
			if err := downloadFile(item.DownloadURL, itemPath); err != nil {
				return fmt.Errorf("failed to download file %s: %w", item.Name, err)
			}
		} else if item.Type == "dir" {
			log.Info("Downloading directory: %s", item.Name)
			// Recursively download subdirectory
			if err := downloadGitHubDirectory(owner, repo, branch, item.Path, itemPath); err != nil {
				return fmt.Errorf("failed to download directory %s: %w", item.Name, err)
			}
		}
	}

	return nil
}

// downloadFile downloads a file from a URL and saves it to the specified path
func downloadFile(url, filepath string) error {
	// Handle data URLs (base64 encoded content)
	if strings.HasPrefix(url, "data:") {
		return downloadDataURL(url, filepath)
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// downloadDataURL handles data URLs with base64 encoded content
func downloadDataURL(dataURL, filepath string) error {
	// Parse data URL format: data:[<mediatype>][;base64],<data>
	parts := strings.SplitN(dataURL, ",", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid data URL format")
	}

	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return fmt.Errorf("failed to decode base64 data: %w", err)
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
