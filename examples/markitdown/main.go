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

	// Example 1: Convert a web page to Markdown
	fmt.Println("=== Example 1: Convert web page to Markdown ===")
	result1, err := agent.Run(ctx, "Use markitdown to convert the webpage https://github.com/microsoft/markitdown to markdown format and show me the content")
	if err != nil {
		log.Printf("Error in Example 1: %v", err)
	} else {
		fmt.Printf("Result:\n%s\n\n", result1)
	}

	// Example 2: Convert PDF to Markdown (if sample.pdf exists)
	fmt.Println("=== Example 2: Convert PDF to Markdown ===")
	if _, err := os.Stat("../pdf/Hello_World.pdf"); err == nil {
		result2, err := agent.Run(ctx, "Use markitdown to convert the PDF file ../pdf/Hello_World.pdf to markdown and show the extracted text")
		if err != nil {
			log.Printf("Error in Example 2: %v", err)
		} else {
			fmt.Printf("Result:\n%s\n\n", result2)
		}
	} else {
		fmt.Println("Skipped: ../pdf/Hello_World.pdf not found")
	}

	// Example 3: Create a sample CSV and convert it to Markdown table
	fmt.Println("=== Example 3: Convert CSV to Markdown table ===")

	// Create a sample CSV file
	csvContent := `Name,Age,Department,Salary
John Doe,30,Engineering,85000
Jane Smith,28,Marketing,72000
Bob Johnson,35,Sales,78000
Alice Brown,32,HR,68000`

	err = os.WriteFile("sample.csv", []byte(csvContent), 0644)
	if err != nil {
		log.Printf("Error creating sample CSV: %v", err)
	} else {
		result3, err := agent.Run(ctx, "Use markitdown to convert the CSV file sample.csv to a markdown table and display it")
		if err != nil {
			log.Printf("Error in Example 3: %v", err)
		} else {
			fmt.Printf("Result:\n%s\n\n", result3)
		}
	}

	// Example 4: Convert JSON to Markdown
	fmt.Println("=== Example 4: Convert JSON to Markdown ===")

	// Create a sample JSON file
	jsonContent := `{
  "company": "Tech Corp",
  "employees": [
    {
      "name": "John Doe",
      "role": "Software Engineer",
      "skills": ["Python", "Go", "JavaScript"]
    },
    {
      "name": "Jane Smith",
      "role": "Product Manager",
      "skills": ["Strategy", "Communication", "Analytics"]
    }
  ],
  "founded": 2020,
  "locations": ["San Francisco", "New York", "London"]
}`

	err = os.WriteFile("sample.json", []byte(jsonContent), 0644)
	if err != nil {
		log.Printf("Error creating sample JSON: %v", err)
	} else {
		result4, err := agent.Run(ctx, "Use markitdown to convert the JSON file sample.json to readable markdown format")
		if err != nil {
			log.Printf("Error in Example 4: %v", err)
		} else {
			fmt.Printf("Result:\n%s\n\n", result4)
		}
	}

	// Example 5: Convert HTML to Markdown
	fmt.Println("=== Example 5: Convert HTML to Markdown ===")

	htmlContent := `<!DOCTYPE html>
<html>
<head>
    <title>Sample Page</title>
</head>
<body>
    <h1>Welcome to MarkItDown</h1>
    <p>This is a <strong>sample HTML</strong> page with various elements:</p>
    <ul>
        <li>Unordered list item 1</li>
        <li>Unordered list item 2</li>
    </ul>
    <h2>Features</h2>
    <ol>
        <li>Convert documents to Markdown</li>
        <li>Preserve structure and formatting</li>
        <li>Support multiple file formats</li>
    </ol>
    <p>Visit <a href="https://github.com/microsoft/markitdown">MarkItDown on GitHub</a> for more info.</p>
</body>
</html>`

	err = os.WriteFile("sample.html", []byte(htmlContent), 0644)
	if err != nil {
		log.Printf("Error creating sample HTML: %v", err)
	} else {
		result5, err := agent.Run(ctx, "Use markitdown to convert the HTML file sample.html to clean markdown format")
		if err != nil {
			log.Printf("Error in Example 5: %v", err)
		} else {
			fmt.Printf("Result:\n%s\n\n", result5)
		}
	}

	fmt.Println("\n=== All examples completed ===")
	fmt.Println("\nGenerated files:")
	fmt.Println("  - sample.csv (sample CSV data)")
	fmt.Println("  - sample.json (sample JSON data)")
	fmt.Println("  - sample.html (sample HTML page)")
	fmt.Println("\nYou can clean up with: rm sample.csv sample.json sample.html")
}
