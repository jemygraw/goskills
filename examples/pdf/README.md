# PDF Skill Example

This example demonstrates how to use the GoSkills library to call the PDF skill for various PDF operations.

## Overview

The PDF skill (located in `../../testdata/skills/document-skills/pdf`) provides comprehensive PDF manipulation capabilities including:

- Extracting text and tables from PDFs
- Creating new PDF documents
- Merging and splitting PDFs
- Handling PDF forms
- Extracting metadata
- Adding watermarks
- Password protection

## Prerequisites

1. Go 1.21 or later
2. OpenAI API key (or compatible API)
3. Python environment with required packages (for PDF skill tools):
   ```bash
   pip install pypdf pdfplumber reportlab
   ```

## Setup

1. Set your OpenAI API key:
   ```bash
   export OPENAI_API_KEY="your-api-key-here"
   ```

2. (Optional) Set custom API base and model:
   ```bash
   export OPENAI_API_BASE="https://api.openai.com/v1"
   export OPENAI_MODEL="gpt-4o"
   ```

## Running the Example

From the `examples/pdf` directory:

```bash
# Run the example
go run main.go
```

Or from the project root:

```bash
# Build and run
go run ./examples/pdf/main.go
```

## What the Example Does

The example demonstrates three common PDF operations:

1. **Extract text from a PDF**: Reads the `Hello_World.pdf` file and extracts its text content
2. **Create a new PDF**: Creates a simple PDF file with custom text using reportlab
3. **Extract metadata**: Retrieves metadata information (title, author, subject) from a PDF

## Customization

You can modify the `main.go` file to try other PDF operations:

### Extract tables from PDF
```go
result, err := agent.Run(ctx, "Extract all tables from invoice.pdf and save to Excel")
```

### Merge PDFs
```go
result, err := agent.Run(ctx, "Merge doc1.pdf, doc2.pdf, and doc3.pdf into merged.pdf")
```

### Split PDF
```go
result, err := agent.Run(ctx, "Split document.pdf into separate pages")
```

### Add watermark
```go
result, err := agent.Run(ctx, "Add a watermark from watermark.pdf to document.pdf")
```

### Fill PDF form
```go
result, err := agent.Run(ctx, "Fill the PDF form template.pdf with data: name='John Doe', email='john@example.com'")
```

## PDF Skill Features

The PDF skill uses several Python libraries:

- **pypdf**: Basic PDF operations (merge, split, rotate, metadata)
- **pdfplumber**: Text and table extraction with layout preservation
- **reportlab**: PDF creation and generation
- **pytesseract** (optional): OCR for scanned PDFs
- **pdf2image** (optional): Convert PDF pages to images

For more details, see the skill documentation at `../../testdata/skills/document-skills/pdf/SKILL.md`.

## Troubleshooting

1. **"OPENAI_API_KEY environment variable is required"**
   - Make sure you've set the OPENAI_API_KEY environment variable

2. **Python package not found errors**
   - Install the required Python packages: `pip install pypdf pdfplumber reportlab`

3. **PDF file not found**
   - Make sure the PDF file exists in the current directory or provide the full path

4. **Tool execution failures**
   - Check that Python is installed and accessible
   - Verify that the Python packages are installed correctly
   - Enable verbose mode to see detailed error messages

## Advanced Usage

For more complex scenarios, you can:

1. Enable loop mode for interactive sessions:
   ```go
   cfg.Loop = true
   agent.RunLoop(ctx, "I want to work with PDFs")
   ```

2. Use custom allowed scripts:
   ```go
   cfg.AllowedScripts = []string{"*.py", "*.sh"}
   ```

3. Integrate with MCP (Model Context Protocol) servers:
   ```go
   mcpClient, err := mcp.NewClient(ctx, mcpConfig)
   agent, err := goskills.NewAgent(cfg, mcpClient)
   ```

## Learn More

- [GoSkills Documentation](../../README.md)
- [PDF Skill Reference](../../testdata/skills/document-skills/pdf/SKILL.md)
- [PDF Forms Guide](../../testdata/skills/document-skills/pdf/forms.md)
- [OpenAI API Documentation](https://platform.openai.com/docs)
