# MarkItDown Skill Example

This example demonstrates how to use the GoSkills library to call the MarkItDown skill for converting various file formats to Markdown.

## Overview

The MarkItDown skill (located in `../../testdata/skills/markitdown`) is a comprehensive document conversion tool that converts 20+ file formats into clean, LLM-optimized Markdown. It supports:

- **Document Formats**: PDF, DOCX, PPTX, XLSX, EPUB
- **Web Content**: HTML, web pages, YouTube transcripts, RSS feeds
- **Structured Data**: CSV, JSON, XML
- **Media Files**: Images (with OCR), Audio (with transcription)
- **Archives**: ZIP files (batch processing)

## Prerequisites

1. Go 1.21 or later
2. OpenAI API key (or compatible API)
3. Python 3.10 or higher
4. MarkItDown Python package:
   ```bash
   # Full installation (all features)
   pip install 'markitdown[all]'

   # Or install specific features only
   pip install 'markitdown[pdf]'      # PDF support
   pip install 'markitdown[docx]'     # Word support
   pip install 'markitdown[pptx]'     # PowerPoint support
   pip install 'markitdown[xlsx]'     # Excel support
   pip install 'markitdown[audio]'    # Audio transcription
   pip install 'markitdown[youtube]'  # YouTube transcripts
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

From the `examples/markitdown` directory:

```bash
# Run the example
go run main.go
```

Or from the project root:

```bash
# Build and run
go run ./examples/markitdown/main.go
```

## What the Example Does

The example demonstrates five different conversion scenarios:

1. **Web page to Markdown**: Converts the MarkItDown GitHub page to Markdown
2. **PDF to Markdown**: Extracts text from a PDF file
3. **CSV to Markdown table**: Converts CSV data to a formatted Markdown table
4. **JSON to Markdown**: Converts structured JSON to readable Markdown
5. **HTML to Markdown**: Converts HTML content to clean Markdown

## Sample Files

The example creates three sample files for demonstration:
- `sample.csv` - Employee data in CSV format
- `sample.json` - Company information in JSON format
- `sample.html` - A simple HTML page

These files are automatically created when you run the example.

## Customization

You can modify the `main.go` file to try other MarkItDown features:

### Convert Office Documents
```go
result, err := agent.Run(ctx, "Use markitdown to convert presentation.pptx to markdown")
```

### Extract YouTube Transcript
```go
result, err := agent.Run(ctx, "Use markitdown to get the transcript from https://youtube.com/watch?v=VIDEO_ID")
```

### OCR on Images
```go
result, err := agent.Run(ctx, "Use markitdown to extract text from image.jpg using OCR")
```

### Transcribe Audio
```go
result, err := agent.Run(ctx, "Use markitdown to transcribe audio.wav to text")
```

### Batch Convert Multiple Files
```go
result, err := agent.Run(ctx, "Use markitdown to convert all PDF files in the documents/ directory to markdown")
```

### Convert ZIP Archive
```go
result, err := agent.Run(ctx, "Use markitdown to extract and convert all files in archive.zip to markdown")
```

### Convert EPUB Book
```go
result, err := agent.Run(ctx, "Use markitdown to convert ebook.epub to markdown")
```

### Enhanced PDF with Azure Document Intelligence
```go
result, err := agent.Run(ctx, "Use markitdown with Azure Document Intelligence to convert complex.pdf with better table extraction")
```

## MarkItDown Features

### Document Conversion
- **PDF**: Text extraction with optional Azure Document Intelligence for better table handling
- **DOCX**: Word documents with structure preservation
- **PPTX**: PowerPoint presentations with slide content
- **XLSX**: Excel spreadsheets as Markdown tables

### Media Processing
- **Images**: EXIF metadata extraction + OCR text recognition
- **Audio**: Speech-to-text transcription (requires speech_recognition)

### Web Content
- **HTML**: Clean conversion preserving structure
- **YouTube**: Video transcript extraction
- **EPUB**: E-book conversion
- **RSS**: Feed content extraction

### Structured Data
- **CSV**: Convert to Markdown tables
- **JSON**: Pretty-printed, readable format
- **XML**: Structured text extraction

### Advanced Features
- **LLM Integration**: Use GPT-4o for image descriptions in presentations
- **Batch Processing**: Process entire directories or ZIP archives
- **Plugin System**: Extensible architecture for custom converters

## Common Use Cases

### 1. Preparing Documents for RAG (Retrieval-Augmented Generation)
```go
result, err := agent.Run(ctx, `Use markitdown to convert these documents for RAG:
- manual.pdf
- guide.docx
- faq.html
Save each as markdown in the knowledge_base/ directory`)
```

### 2. Document Analysis Pipeline
```go
result, err := agent.Run(ctx, "Use markitdown to convert all documents in input/ to markdown in output/")
```

### 3. Extract Meeting Notes from Audio
```go
result, err := agent.Run(ctx, "Use markitdown to transcribe meeting.wav and save as meeting_notes.md")
```

### 4. Archive Research Papers
```go
result, err := agent.Run(ctx, "Use markitdown to convert research_papers.zip to markdown files")
```

### 5. Create Knowledge Base from Various Sources
```go
result, err := agent.Run(ctx, `Use markitdown to create a knowledge base from:
- Company policies (PDF)
- Product documentation (DOCX)
- FAQs (HTML)
- Tutorial videos (YouTube URLs)`)
```

## Output Format

MarkItDown produces clean, token-efficient Markdown optimized for LLM consumption:
- ✅ Preserves document structure (headings, lists, tables)
- ✅ Maintains hyperlinks and formatting
- ✅ Includes relevant metadata
- ✅ No temporary files (streaming approach)
- ✅ Consistent output format across different input types

## Troubleshooting

1. **"OPENAI_API_KEY environment variable is required"**
   - Set the OPENAI_API_KEY environment variable

2. **"markitdown command not found" or ImportError**
   - Install MarkItDown: `pip install 'markitdown[all]'`

3. **PDF conversion fails**
   - Install PDF support: `pip install 'markitdown[pdf]'`

4. **OCR not working on images**
   - Make sure tesseract-ocr is installed
   - Linux: `sudo apt-get install tesseract-ocr`
   - macOS: `brew install tesseract`

5. **Audio transcription fails**
   - Install audio support: `pip install 'markitdown[audio]'`
   - Requires speech_recognition package

6. **YouTube transcript extraction fails**
   - Install YouTube support: `pip install 'markitdown[youtube]'`
   - Requires youtube-transcript-api package

7. **"Python 3.10 or higher required"**
   - Upgrade Python to version 3.10 or later

## Advanced Usage

### Using with Azure Document Intelligence
For enhanced PDF processing with better table extraction:

```bash
# Set Azure credentials
export AZURE_DOCINTEL_ENDPOINT="your-endpoint"
export AZURE_DOCINTEL_KEY="your-key"
```

```go
result, err := agent.Run(ctx, "Use markitdown with Azure Document Intelligence to convert complex_tables.pdf")
```

### Using with LLM for Image Descriptions
For AI-powered image descriptions in presentations:

```go
result, err := agent.Run(ctx, "Use markitdown with GPT-4o to convert presentation.pptx with detailed image descriptions")
```

### Batch Processing Script
Use the included batch processing utility:

```bash
python ../../testdata/skills/markitdown/scripts/batch_convert.py ./input ./output
```

### Command-Line Usage
MarkItDown can also be used directly from command line:

```bash
# Convert single file
markitdown document.pdf -o output.md

# Convert with specific options
markitdown --help
```

## Performance Tips

1. **Batch Processing**: Process multiple files together for efficiency
2. **Modular Installation**: Install only the features you need to reduce dependencies
3. **Streaming Approach**: MarkItDown doesn't create temporary files, making it fast
4. **Azure Document Intelligence**: Use for complex PDFs with tables (requires Azure subscription)

## Clean Up

After running the examples, you can remove the generated sample files:

```bash
rm sample.csv sample.json sample.html
```

## Learn More

- [GoSkills Documentation](../../README.md)
- [MarkItDown Skill Reference](../../testdata/skills/markitdown/SKILL.md)
- [Document Conversion Guide](../../testdata/skills/markitdown/references/document_conversion.md)
- [Media Processing Guide](../../testdata/skills/markitdown/references/media_processing.md)
- [Web Content Guide](../../testdata/skills/markitdown/references/web_content.md)
- [Structured Data Guide](../../testdata/skills/markitdown/references/structured_data.md)
- [Advanced Integrations](../../testdata/skills/markitdown/references/advanced_integrations.md)
- [MarkItDown GitHub Repository](https://github.com/microsoft/markitdown)

## Related Examples

- [PDF Example](../pdf/) - Working with PDF documents
- [DOCX Example](../docx/) - Working with Word documents

## Example Output

When you run the example, you'll see conversions of:
- A GitHub web page converted to clean Markdown
- PDF text extraction
- CSV data formatted as a Markdown table
- JSON data in readable Markdown format
- HTML content converted to Markdown

All optimized for LLM processing and analysis!
