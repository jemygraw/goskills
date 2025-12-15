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

	// Example 1: Create a new Word document
	fmt.Println("=== Example 1: Create a new Word document ===")
	result1, err := agent.Run(ctx, `Create a new Word document named 'report.docx' with:
- Title: "Monthly Report"
- Heading: "Executive Summary"
- A paragraph with text: "This is our monthly progress report for the current quarter."
- A bullet list with items: "Revenue increased by 15%", "New clients: 50", "Team size: 25"`)
	if err != nil {
		log.Printf("Error in Example 1: %v", err)
	} else {
		fmt.Printf("Result:\n%s\n\n", result1)
	}

	// Example 2: Extract text from a Word document (if exists)
	// First, check if report.docx exists
	if _, err := os.Stat("report.docx"); err == nil {
		fmt.Println("=== Example 2: Extract text from Word document ===")
		result2, err := agent.Run(ctx, "Extract all text content from report.docx and display it")
		if err != nil {
			log.Printf("Error in Example 2: %v", err)
		} else {
			fmt.Printf("Result:\n%s\n\n", result2)
		}
	} else {
		fmt.Println("=== Example 2: Skipped (report.docx not found) ===\n")
	}

	// Example 3: Create a professional document with formatting
	fmt.Println("=== Example 3: Create a formatted business letter ===")
	result3, err := agent.Run(ctx, `Create a Word document named 'letter.docx' with a professional business letter:
- Date: December 15, 2025
- To: John Doe, ABC Corporation
- From: Jane Smith, XYZ Company
- Subject: Partnership Proposal
- Body: A formal paragraph proposing a business partnership
- Closing: "Sincerely" with signature line`)
	if err != nil {
		log.Printf("Error in Example 3: %v", err)
	} else {
		fmt.Printf("Result:\n%s\n\n", result3)
	}

	// Example 4: Create a document with table
	fmt.Println("=== Example 4: Create a document with a table ===")
	result4, err := agent.Run(ctx, `Create a Word document named 'table.docx' with a table containing:
- Headers: Product, Quantity, Price
- Row 1: Laptop, 10, $1200
- Row 2: Mouse, 50, $25
- Row 3: Keyboard, 30, $75`)
	if err != nil {
		log.Printf("Error in Example 4: %v", err)
	} else {
		fmt.Printf("Result:\n%s\n\n", result4)
	}

	fmt.Println("\n=== All examples completed ===")
}
