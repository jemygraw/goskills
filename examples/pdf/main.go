package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/smallnest/goskills"
)

func main() {
	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Configure the agent
	cfg := goskills.RunnerConfig{
		APIKey:           apiKey,
		SkillsDir:        "../../testdata/skills", // Path to skills directory
		Verbose:          true,
		AutoApproveTools: true,
		Loop:             false,
	}

	// Get API base if provided
	if apiBase := os.Getenv("OPENAI_API_BASE"); apiBase != "" {
		cfg.APIBase = apiBase
	}

	// Get model if provided
	if model := os.Getenv("OPENAI_MODEL"); model != "" {
		cfg.Model = model
	}

	ctx := context.Background()

	// Create agent
	agent, err := goskills.NewAgent(cfg, nil)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// Example 1: Extract text from PDF
	fmt.Println("=== Example 1: Extract text from a PDF ===")
	result1, err := agent.Run(ctx, "Extract text from the PDF file Hello_World.pdf")
	if err != nil {
		log.Printf("Error in Example 1: %v", err)
	} else {
		fmt.Printf("Result:\n%s\n\n", result1)
	}

	// Example 2: Create a simple PDF
	fmt.Println("=== Example 2: Create a simple PDF ===")
	result2, err := agent.Run(ctx, "Create a PDF file named 'welcome.pdf' with the text 'Welcome to PDF Skills!' using reportlab")
	if err != nil {
		log.Printf("Error in Example 2: %v", err)
	} else {
		fmt.Printf("Result:\n%s\n\n", result2)
	}

	// Example 3: Merge PDFs (if you have multiple PDFs)
	fmt.Println("=== Example 3: Get PDF metadata ===")
	result3, err := agent.Run(ctx, "Extract metadata (title, author, subject) from Hello_World.pdf")
	if err != nil {
		log.Printf("Error in Example 3: %v", err)
	} else {
		fmt.Printf("Result:\n%s\n\n", result3)
	}

	fmt.Println("\n=== All examples completed ===")
}
