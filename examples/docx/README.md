# DOCX Skill Example

This example demonstrates how to use the GoSkills library to call the DOCX skill for various Word document operations.

## Overview

The DOCX skill (located in `../../testdata/skills/document-skills/docx`) provides comprehensive Word document manipulation capabilities including:

- Creating new Word documents with rich formatting
- Extracting text and content from existing documents
- Editing and modifying existing documents
- Working with tracked changes (redlining)
- Adding and managing comments
- Tables, lists, and complex formatting
- Converting documents to other formats

## Prerequisites

1. Go 1.21 or later
2. OpenAI API key (or compatible API)
3. Node.js and npm (for docx-js library):
   ```bash
   npm install -g docx
   ```
4. Python environment with required packages:
   ```bash
   pip install defusedxml
   ```
5. (Optional) Additional tools for advanced features:
   ```bash
   # For text extraction and conversion
   sudo apt-get install pandoc

   # For document to PDF conversion
   sudo apt-get install libreoffice

   # For PDF to image conversion
   sudo apt-get install poppler-utils
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

From the `examples/docx` directory:

```bash
# Run the example
go run main.go
```

Or from the project root:

```bash
# Build and run
go run ./examples/docx/main.go
```

## What the Example Does

The example demonstrates four common DOCX operations:

1. **Create a new document**: Creates a Word document with title, heading, paragraph, and bullet list
2. **Extract text**: Reads and extracts text content from the created document
3. **Create a business letter**: Generates a professionally formatted business letter
4. **Create a document with table**: Creates a document containing a formatted table

## Customization

You can modify the `main.go` file to try other DOCX operations:

### Create a document with images
```go
result, err := agent.Run(ctx, "Create a Word document with a title and embed the image logo.png")
```

### Extract content from existing document
```go
result, err := agent.Run(ctx, "Convert document.docx to markdown format and show the content")
```

### Edit existing document with tracked changes
```go
result, err := agent.Run(ctx, "Edit contract.docx and change '30 days' to '60 days' using tracked changes")
```

### Add comments to document
```go
result, err := agent.Run(ctx, "Add a comment to the first paragraph of document.docx saying 'Please review this section'")
```

### Create a complex formatted document
```go
result, err := agent.Run(ctx, `Create a document named 'proposal.docx' with:
- Cover page with title and author
- Table of contents
- Multiple sections with headings
- Numbered and bullet lists
- Tables with formatted data
- Footer with page numbers`)
```

### Extract and analyze document structure
```go
result, err := agent.Run(ctx, "Analyze the structure of document.docx and list all headings and sections")
```

### Convert document to PDF
```go
result, err := agent.Run(ctx, "Convert document.docx to PDF format")
```

## DOCX Skill Features

The DOCX skill supports multiple tools and workflows:

### 1. Creating New Documents (docx-js)
- Rich text formatting (bold, italic, underline, colors)
- Paragraphs and headings
- Bullet and numbered lists
- Tables with styling
- Images and embedded media
- Headers and footers
- Page numbers and sections

### 2. Reading/Analyzing Documents
- Text extraction using pandoc
- Markdown conversion
- Structure analysis
- Raw XML access for detailed inspection

### 3. Editing Existing Documents (Document Library)
- Content modification
- Tracked changes (redlining)
- Comments and annotations
- Format preservation
- Complex OOXML manipulation

### 4. Document Conversion
- DOCX to PDF
- DOCX to Markdown
- DOCX to Images

## Document Workflows

### Creating a New Document
The skill uses **docx-js** (JavaScript/TypeScript library) to create Word documents from scratch with rich formatting options.

### Editing an Existing Document
The skill uses the **Document library** (Python) for OOXML manipulation, providing both high-level methods and direct XML access.

### Redlining Workflow (Tracked Changes)
For professional document review with tracked changes:
1. Convert document to markdown to review current content
2. Identify changes needed
3. Unpack the DOCX file (it's a ZIP containing XML)
4. Apply tracked changes using Python scripts
5. Pack the modified files back into DOCX format
6. Verify all changes were applied correctly

## Troubleshooting

1. **"OPENAI_API_KEY environment variable is required"**
   - Make sure you've set the OPENAI_API_KEY environment variable

2. **"docx command not found"**
   - Install the docx-js library: `npm install -g docx`

3. **"pandoc command not found"**
   - Install pandoc: `sudo apt-get install pandoc` (Linux) or `brew install pandoc` (macOS)

4. **Python package errors**
   - Install required packages: `pip install defusedxml`

5. **Document creation fails**
   - Ensure Node.js and npm are installed
   - Verify docx-js is properly installed globally

6. **Tool execution failures**
   - Check that Python is installed and accessible
   - Verify that Node.js is installed for docx-js
   - Enable verbose mode to see detailed error messages

## Advanced Usage

### Interactive Mode with Loop
For complex document operations, enable loop mode for an interactive session:
```go
cfg.Loop = true
agent.RunLoop(ctx, "I want to create and edit Word documents")
```

### Working with Multiple Documents
```go
result, err := agent.Run(ctx, `Merge the content from doc1.docx and doc2.docx into a new document merged.docx`)
```

### Batch Operations
```go
result, err := agent.Run(ctx, `Convert all .docx files in the current directory to markdown format`)
```

### Custom Styling
```go
result, err := agent.Run(ctx, `Create a document with custom styles:
- Heading 1: Arial 18pt, bold, blue
- Body: Times New Roman 12pt
- Highlight important text in yellow`)
```

## File Structure

After running the examples, you'll see:
- `report.docx` - Monthly report with bullet points
- `letter.docx` - Professional business letter
- `table.docx` - Document with formatted table

## Learn More

- [GoSkills Documentation](../../README.md)
- [DOCX Skill Reference](../../testdata/skills/document-skills/docx/SKILL.md)
- [DOCX-JS Documentation](../../testdata/skills/document-skills/docx/docx-js.md)
- [OOXML Manipulation Guide](../../testdata/skills/document-skills/docx/ooxml.md)
- [OpenAI API Documentation](https://platform.openai.com/docs)

## Related Examples

- [PDF Example](../pdf/) - Working with PDF documents
- PPTX Example - Working with PowerPoint presentations (if available)

## Tips

1. **Start Simple**: Begin with basic document creation before trying complex formatting
2. **Use Tracked Changes**: For professional document review, always use the redlining workflow
3. **Check Dependencies**: Ensure all required tools (pandoc, docx-js, Python packages) are installed
4. **Verbose Mode**: Enable verbose mode (`Verbose: true`) to see detailed execution steps
5. **Auto-Approve**: Use `AutoApproveTools: true` for automated workflows, or set to `false` for manual review
