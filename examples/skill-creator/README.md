# Skill Creator Example

This example demonstrates how to use the GoSkills library to call the Skill Creator skill for creating, validating, and packaging new skills.

## Overview

The Skill Creator skill (located in `../../testdata/skills/skill-creator`) is a meta-skill that helps you create new skills for Claude. It provides:

- **Skill Initialization**: Generate properly structured skill templates
- **Skill Validation**: Check if skills meet all requirements
- **Skill Packaging**: Package skills into distributable zip files
- **Best Practices Guidance**: Learn how to create effective skills

## What are Skills?

Skills are modular, self-contained packages that extend Claude's capabilities by providing:

1. **Specialized workflows** - Multi-step procedures for specific domains
2. **Tool integrations** - Instructions for working with specific file formats or APIs
3. **Domain expertise** - Company-specific knowledge, schemas, business logic
4. **Bundled resources** - Scripts, references, and assets for complex tasks

## Skill Structure

Every skill consists of:

```
skill-name/
├── SKILL.md (required)
│   ├── YAML frontmatter metadata (name, description)
│   └── Markdown instructions
└── Bundled Resources (optional)
    ├── scripts/          - Executable code (Python/Bash/etc.)
    ├── references/       - Documentation loaded as needed
    └── assets/           - Files used in output (templates, images, etc.)
```

## Prerequisites

1. Go 1.21 or later
2. OpenAI API key (or compatible API)
3. Python 3.7 or higher (for skill creation scripts)

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

From the `examples/skill-creator` directory:

```bash
# Run the example
go run main.go
```

Or from the project root:

```bash
# Build and run
go run ./examples/skill-creator/main.go
```

## What the Example Does

The example demonstrates five key skill creation operations:

1. **Initialize a skill**: Creates an `image-optimizer` skill with proper structure
2. **Validate a skill**: Checks if the created skill meets all requirements
3. **Package a skill**: Creates a distributable zip file of the skill
4. **Create another skill**: Generates a `postgres-helper` database skill
5. **Get best practices**: Learns about effective skill creation

## Customization

You can modify the `main.go` file to create different types of skills:

### Create a PDF Processing Skill
```go
result, err := agent.Run(ctx, `Use skill-creator to create a 'pdf-editor' skill that helps with:
- Rotating PDF pages
- Merging multiple PDFs
- Splitting PDFs into separate pages
- Extracting text from PDFs
Include a script for PDF rotation in the scripts/ directory`)
```

### Create a Company-Specific Skill
```go
result, err := agent.Run(ctx, `Use skill-creator to create a 'company-policies' skill that:
- Provides company HR policies and procedures
- Includes references for employee handbook
- Has templates for common forms
- Contains assets for company branding`)
```

### Create a Frontend Development Skill
```go
result, err := agent.Run(ctx, `Use skill-creator to create a 'react-builder' skill for building React apps:
- Include boilerplate React project structure in assets/
- Add reference documentation for React best practices
- Include scripts for project setup and build`)
```

### Create a Data Analysis Skill
```go
result, err := agent.Run(ctx, `Use skill-creator to create a 'data-analyst' skill that:
- Helps with data cleaning and transformation
- Provides reference for pandas and numpy best practices
- Includes scripts for common data operations
- Has templates for visualization`)
```

### Create an API Integration Skill
```go
result, err := agent.Run(ctx, `Use skill-creator to create a 'stripe-integration' skill:
- Document Stripe API endpoints and authentication
- Include reference for payment flow patterns
- Add scripts for common payment operations
- Provide examples of webhook handling`)
```

### Validate an Existing Skill
```go
result, err := agent.Run(ctx, "Use skill-creator to validate the skill at ./my-custom-skill")
```

### Package a Skill for Distribution
```go
result, err := agent.Run(ctx, "Use skill-creator to package the skill at ./my-custom-skill into ./dist")
```

## Skill Creation Process

The Skill Creator follows a structured 6-step process:

### Step 1: Understanding with Concrete Examples
- Identify what functionality the skill should support
- Gather examples of how the skill will be used
- Understand trigger patterns

### Step 2: Planning Reusable Contents
- Determine what scripts are needed
- Identify reference documentation required
- Plan assets and templates

### Step 3: Initializing the Skill
- Use `init_skill.py` to generate template
- Creates proper directory structure
- Generates SKILL.md with placeholders

### Step 4: Editing the Skill
- Fill in SKILL.md with instructions
- Add scripts, references, and assets
- Use imperative/infinitive form for instructions

### Step 5: Packaging the Skill
- Validate skill structure and content
- Create distributable zip file
- Fix any validation errors

### Step 6: Iterate
- Test the skill on real tasks
- Identify improvements
- Update and retest

## Bundled Resource Types

### Scripts (`scripts/`)
Executable code for deterministic tasks:
- **When to use**: Repeatedly rewritten code or need reliability
- **Example**: `rotate_pdf.py`, `optimize_image.py`
- **Benefits**: Token efficient, deterministic execution

### References (`references/`)
Documentation loaded into context as needed:
- **When to use**: Domain knowledge Claude should reference
- **Example**: `database_schema.md`, `api_docs.md`
- **Benefits**: Keeps SKILL.md lean, loaded on demand

### Assets (`assets/`)
Files used in output, not loaded into context:
- **When to use**: Templates, images, boilerplate code
- **Example**: `logo.png`, `template.html`, `project-boilerplate/`
- **Benefits**: Separate output resources from documentation

## Progressive Disclosure

Skills use a three-level loading system for efficiency:

1. **Metadata** (name + description) - Always in context (~100 words)
2. **SKILL.md body** - When skill triggers (<5k words)
3. **Bundled resources** - As needed by Claude (Unlimited)

## Writing Style Guidelines

When creating skills:
- ✅ Use imperative/infinitive form ("To accomplish X, do Y")
- ✅ Be objective and instructional
- ❌ Avoid second person ("You should...")
- ❌ Avoid conversational tone

## Skill Creation Scripts

The skill-creator provides three helper scripts:

### 1. init_skill.py
Initialize a new skill with proper structure:
```bash
python scripts/init_skill.py <skill-name> --path <output-directory>
```

Creates:
- Skill directory structure
- SKILL.md template with frontmatter
- Example scripts/, references/, and assets/ directories

### 2. quick_validate.py
Validate skill structure and content:
```bash
python scripts/quick_validate.py <path-to-skill>
```

Checks:
- YAML frontmatter format
- Required fields (name, description)
- Naming conventions
- File organization

### 3. package_skill.py
Package skill for distribution:
```bash
python scripts/package_skill.py <path-to-skill> [output-directory]
```

Performs:
- Automatic validation
- Creates zip file named after skill
- Maintains proper directory structure

## Example Skills to Create

Here are some practical skills you might want to create:

### Business Skills
- **invoice-generator**: Generate professional invoices
- **meeting-notes**: Structure and format meeting notes
- **email-templates**: Company email templates and style

### Development Skills
- **code-reviewer**: Code review best practices and checklists
- **test-writer**: Generate unit and integration tests
- **docker-helper**: Docker configuration and best practices

### Data Skills
- **excel-analyzer**: Excel data analysis and formulas
- **sql-optimizer**: Database query optimization
- **data-visualizer**: Data visualization best practices

### Content Skills
- **blog-writer**: Blog post structure and SEO
- **social-media**: Social media content and scheduling
- **documentation**: Technical documentation templates

### Domain-Specific Skills
- **legal-contract**: Legal contract templates and clauses
- **medical-notes**: Medical documentation standards
- **financial-analysis**: Financial modeling and analysis

## Troubleshooting

1. **"OPENAI_API_KEY environment variable is required"**
   - Set the OPENAI_API_KEY environment variable

2. **"Python not found" errors**
   - Ensure Python 3.7+ is installed and in PATH

3. **Validation errors**
   - Check SKILL.md has proper YAML frontmatter
   - Ensure name and description fields are present
   - Verify directory structure matches requirements

4. **Packaging fails**
   - Run validation first to identify issues
   - Fix reported errors before packaging
   - Ensure all referenced files exist

## Best Practices

### Metadata Quality
- Be specific in skill descriptions
- Use third-person form ("This skill should be used when...")
- Clearly define when the skill should trigger

### SKILL.md Content
- Keep under 5k words
- Focus on procedural knowledge
- Reference bundled resources
- Use clear, imperative instructions

### Bundled Resources
- Don't duplicate information
- Use references for detailed documentation
- Keep scripts focused and reusable
- Organize assets by purpose

### Progressive Disclosure
- Move large reference material to separate files
- Keep SKILL.md lean with essential instructions
- Use grep patterns for large reference files

## Output Structure

After running the examples, you'll see:

```
examples/skill-creator/
├── my-skills/
│   ├── image-optimizer/
│   │   ├── SKILL.md
│   │   ├── scripts/
│   │   ├── references/
│   │   └── assets/
│   └── postgres-helper/
│       ├── SKILL.md
│       ├── scripts/
│       ├── references/
│       └── assets/
└── dist/
    ├── image-optimizer.zip
    └── postgres-helper.zip
```

## Learn More

- [GoSkills Documentation](../../README.md)
- [Skill Creator Reference](../../testdata/skills/skill-creator/SKILL.md)
- [Agent Skills Specification](../../testdata/skills/agent_skills_spec.md)
- [Claude Skills Documentation](https://docs.claude.com/en/docs/agents-and-tools/agent-skills/)

## Related Examples

- [PDF Example](../pdf/) - Working with PDF documents
- [DOCX Example](../docx/) - Working with Word documents
- [MarkItDown Example](../markitdown/) - Converting various formats to Markdown

## Next Steps

After creating skills:
1. Test them with real use cases
2. Iterate based on feedback
3. Share with your team or community
4. Document usage examples
5. Maintain and update as needs evolve

## Clean Up

To remove generated files:

```bash
# Remove created skills
rm -rf ./my-skills

# Remove packaged files
rm -rf ./dist
```
