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

	// Example 1: Initialize a new skill
	fmt.Println("=== Example 1: Initialize a new skill ===")
	result1, err := agent.Run(ctx, `Use skill-creator to initialize a new skill called 'image-optimizer' in the ./my-skills directory.
The skill should help optimize images for web use by:
- Resizing images to specific dimensions
- Compressing images to reduce file size
- Converting between formats (PNG, JPEG, WebP)`)
	if err != nil {
		log.Printf("Error in Example 1: %v", err)
	} else {
		fmt.Printf("Result:\n%s\n\n", result1)
	}

	// Example 2: Validate a skill (if created)
	if _, err := os.Stat("./my-skills/image-optimizer"); err == nil {
		fmt.Println("=== Example 2: Validate the created skill ===")
		result2, err := agent.Run(ctx, "Use skill-creator to validate the skill at ./my-skills/image-optimizer and check if it meets all requirements")
		if err != nil {
			log.Printf("Error in Example 2: %v", err)
		} else {
			fmt.Printf("Result:\n%s\n\n", result2)
		}

		// Example 3: Package the skill
		fmt.Println("=== Example 3: Package the skill for distribution ===")
		result3, err := agent.Run(ctx, "Use skill-creator to package the skill at ./my-skills/image-optimizer into a distributable zip file in the ./dist directory")
		if err != nil {
			log.Printf("Error in Example 3: %v", err)
		} else {
			fmt.Printf("Result:\n%s\n\n", result3)
		}
	} else {
		fmt.Println("=== Example 2 & 3: Skipped (skill not created) ===")
	}

	// Example 4: Create a database skill
	fmt.Println("=== Example 4: Create a database query skill ===")
	result4, err := agent.Run(ctx, `Use skill-creator to create a new skill called 'postgres-helper' in ./my-skills directory.
This skill should help with PostgreSQL database operations:
- Writing optimized SQL queries
- Explaining query execution plans
- Suggesting indexes for performance
- Include a reference file with common PostgreSQL functions and best practices`)
	if err != nil {
		log.Printf("Error in Example 4: %v", err)
	} else {
		fmt.Printf("Result:\n%s\n\n", result4)
	}

	// Example 5: Get guidance on skill creation best practices
	fmt.Println("=== Example 5: Get skill creation best practices ===")
	result5, err := agent.Run(ctx, `What are the best practices for creating an effective skill?
Specifically, I want to understand:
- When to use scripts vs references vs assets
- How to write good skill descriptions
- How to organize skill content`)
	if err != nil {
		log.Printf("Error in Example 5: %v", err)
	} else {
		fmt.Printf("Result:\n%s\n\n", result5)
	}

	fmt.Println("\n=== All examples completed ===")
	fmt.Println("\nCreated directories:")
	fmt.Println("  - ./my-skills/image-optimizer (if successful)")
	fmt.Println("  - ./my-skills/postgres-helper (if successful)")
	fmt.Println("  - ./dist (for packaged skills)")
	fmt.Println("\nYou can explore these directories to see the generated skill structure")
}
