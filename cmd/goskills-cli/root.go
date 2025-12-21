package main

import (
	"os"

	"github.com/smallnest/goskills/log"
	"github.com/spf13/cobra"
)

// Version is the version of the tool, set at build time
var Version = "v0.1.0"

var rootCmd = &cobra.Command{
	Use:   "goskills-cli",
	Short: "A CLI tool for creating and managing Claude skills.",
	Long: `goskills-cli is a command-line interface to help you develop, parse,
and manage Claude Skill packages.`,
	Version: Version,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error("command execution failed: %v", err)
		os.Exit(1)
	}
}
