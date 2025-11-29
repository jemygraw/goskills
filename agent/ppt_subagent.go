package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

// PPTSubagent generates a modern HTML presentation from content.
type PPTSubagent struct {
	client             *openai.Client
	model              string
	verbose            bool
	interactionHandler InteractionHandler
	outputDir          string
}

// NewPPTSubagent creates a new PPTSubagent.
func NewPPTSubagent(client *openai.Client, model string, verbose bool, interactionHandler InteractionHandler, outputDir string) *PPTSubagent {
	return &PPTSubagent{
		client:             client,
		model:              model,
		verbose:            verbose,
		interactionHandler: interactionHandler,
		outputDir:          outputDir,
	}
}

// Type returns the task type this subagent handles.
func (p *PPTSubagent) Type() TaskType {
	return TaskTypePPT
}

// Slide represents a single slide in the presentation.
type Slide struct {
	Title   string   `json:"title"`
	Content []string `json:"content"`          // Bullet points or paragraphs
	Image   string   `json:"image,omitempty"`  // Image description or URL
	Layout  string   `json:"layout,omitempty"` // e.g., "title-center", "split-image-right", "bullets"
}

// Execute generates a PPT from the input content.
func (p *PPTSubagent) Execute(ctx context.Context, task Task) (Result, error) {
	if p.verbose {
		fmt.Println("📊 PPT Subagent")
	}
	if p.interactionHandler != nil {
		p.interactionHandler.Log(fmt.Sprintf("> PPT Subagent: %s", task.Description))
	}

	// Ensure output directory exists
	if err := os.MkdirAll(p.outputDir, 0755); err != nil {
		return Result{
			TaskType: TaskTypePPT,
			Success:  false,
			Error:    fmt.Sprintf("failed to create output directory: %v", err),
		}, err
	}

	// Get content from parameters or description
	content, ok := task.Parameters["content"].(string)
	if !ok || content == "Use the content from the previous REPORT task." {
		// Try to get from context (passed from previous task)
		if ctxContent, ok := task.Parameters["context"].([]string); ok && len(ctxContent) > 0 {
			// Try to find the output from the REPORT task
			var foundReport bool
			for i := len(ctxContent) - 1; i >= 0; i-- {
				if strings.Contains(ctxContent[i], "Output from REPORT task:") {
					content = ctxContent[i]
					// Extract the content after the header
					if idx := strings.Index(content, "\n"); idx != -1 {
						content = content[idx+1:]
					}
					foundReport = true
					break
				}
			}

			if !foundReport {
				// If no REPORT output found, use the last task's output
				content = ctxContent[len(ctxContent)-1]
				// Extract the content after the header if present
				if idx := strings.Index(content, "Output from "); idx != -1 {
					if newlineIdx := strings.Index(content[idx:], "\n"); newlineIdx != -1 {
						content = content[idx+newlineIdx+1:]
					}
				}
			}
			content = strings.TrimSpace(content)
		} else if !ok {
			content = task.Description
		}
	}

	if p.verbose {
		fmt.Println("  Generating slide structure...")
	}

	// 1. Generate Slide Structure
	slides, err := p.generateSlides(ctx, content)
	if err != nil {
		return Result{
			TaskType: TaskTypePPT,
			Success:  false,
			Error:    fmt.Sprintf("failed to generate slides: %v", err),
		}, err
	}

	if p.verbose {
		fmt.Printf("  ✓ Generated %d slides\n", len(slides))
	}

	// 2. Generate and Build
	url, err := p.GenerateAndBuild(slides)
	if err != nil {
		// Log detailed error to terminal/logs
		if p.verbose {
			fmt.Printf("❌ PPT Build Failed: %v\n", err)
		}
		if p.interactionHandler != nil {
			// We might want to log it here too, but maybe truncated or full?
			// The user said "print to terminal logs", implying server side.
			// But interactionHandler.Log sends to the UI.
			// Let's log a simplified message to UI and keep full detail in terminal.
			p.interactionHandler.Log("❌ PPT Build Failed. Check server logs for details.")
		}

		// Return a generic error message in the Result so the UI doesn't get cluttered
		return Result{
			TaskType: TaskTypePPT,
			Success:  false,
			Error:    "Presentation build failed. Please check the server logs for details.",
		}, nil
	}

	return Result{
		TaskType: TaskTypePPT,
		Success:  true,
		Output:   fmt.Sprintf("Presentation generated successfully. View it at: %s", url),
		Metadata: map[string]interface{}{
			"ppt_url": url,
			"slides":  slides,
		},
	}, nil
}

// GenerateAndBuild generates the markdown and builds the Slidev project.
func (p *PPTSubagent) GenerateAndBuild(slides []Slide) (string, error) {
	timestamp := time.Now().Unix()
	dirName := fmt.Sprintf("ppt_%d", timestamp)
	projectDir := filepath.Join(p.outputDir, dirName)

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create project directory: %v", err)
	}

	markdown := p.generateSlidevMarkdown(slides)
	if err := os.WriteFile(filepath.Join(projectDir, "slides.md"), []byte(markdown), 0644); err != nil {
		return "", fmt.Errorf("failed to write slides.md: %v", err)
	}

	if p.verbose {
		fmt.Printf("  ✓ Generated slides.md in %s\n", projectDir)
	}

	// Build with Slidev
	basePath := fmt.Sprintf("/generated/%s/dist/", dirName)

	// Create a simple package.json
	packageJson := `{
  "name": "slidev-project",
  "private": true,
  "scripts": {
    "build": "slidev build --out dist --base "
  },
  "dependencies": {
    "@slidev/cli": "^0.48.0",
    "@slidev/theme-default": "latest",
    "vue": "^3.4.0"
  }
}`
	packageJson = strings.Replace(packageJson, "--base ", "--base "+basePath, 1)

	if err := os.WriteFile(filepath.Join(projectDir, "package.json"), []byte(packageJson), 0644); err != nil {
		return "", fmt.Errorf("failed to write package.json: %v", err)
	}

	// Run npm install
	if p.verbose {
		fmt.Println("  Installing dependencies (npm install)...")
	}
	if p.interactionHandler != nil {
		p.interactionHandler.Log("Installing dependencies...")
	}

	installCmd := exec.Command("npm", "install")
	installCmd.Dir = projectDir
	if output, err := installCmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("npm install failed: %v\nOutput: %s", err, string(output))
	}

	// Run npm run build
	if p.verbose {
		fmt.Println("  Building Slidev project (npm run build)...")
	}
	if p.interactionHandler != nil {
		p.interactionHandler.Log("Building presentation...")
	}

	buildCmd := exec.Command("npm", "run", "build")
	buildCmd.Dir = projectDir
	if output, err := buildCmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("slidev build failed: %v\nOutput: %s", err, string(output))
	}

	if p.verbose {
		fmt.Println("  ✓ Build complete")
	}
	if p.interactionHandler != nil {
		p.interactionHandler.Log("✓ Presentation built successfully")
	}

	return fmt.Sprintf("%sindex.html", basePath), nil
}

func (p *PPTSubagent) generateSlides(ctx context.Context, content string) ([]Slide, error) {
	systemPrompt := `You are a professional presentation designer. Your goal is to convert the provided text into a structured slide deck (5-20 slides).
The design should be modern, concise, and engaging.

Output ONLY a JSON array of objects, where each object represents a slide with:
- "title": The slide title.
- "content": An array of strings (bullet points or short paragraphs).
- "image": A description of an image that would fit this slide (for future generation) or a placeholder URL.
- "layout": Suggested layout ("title-center", "split-image-right", "bullets", "quote").

Ensure the first slide is a Title slide and the last is a Thank You/Conclusion slide.
Keep text concise. Use bullet points where possible.

Example:
[
  {"title": "The Future of AI", "content": ["AI is evolving rapidly", "Impact on all industries"], "layout": "title-center"},
  {"title": "Key Trends", "content": ["Generative Models", "Agentic Workflows"], "layout": "bullets"}
]`

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("Create a slide deck from this content (Language: Chinese):\n\n%s", content),
		},
	}

	req := openai.ChatCompletionRequest{
		Model:       p.model,
		Messages:    messages,
		Temperature: 0.7,
	}

	resp, err := p.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	jsonContent := resp.Choices[0].Message.Content

	// Clean up markdown code blocks if present
	if idx := strings.Index(jsonContent, "```json"); idx != -1 {
		jsonContent = jsonContent[idx+7:]
	} else if idx := strings.Index(jsonContent, "```"); idx != -1 {
		jsonContent = jsonContent[idx+3:]
	}
	if idx := strings.LastIndex(jsonContent, "```"); idx != -1 {
		jsonContent = jsonContent[:idx]
	}
	jsonContent = strings.TrimSpace(jsonContent)

	var slides []Slide
	if err := json.Unmarshal([]byte(jsonContent), &slides); err != nil {
		return nil, fmt.Errorf("failed to parse slides JSON: %w", err)
	}

	return slides, nil
}

func (p *PPTSubagent) generateSlidevMarkdown(slides []Slide) string {
	var sb strings.Builder

	// 1. Global Frontmatter
	sb.WriteString("---\n")
	sb.WriteString("theme: default\n")
	sb.WriteString("highlighter: shiki\n")
	sb.WriteString("lineNumbers: false\n")
	sb.WriteString("info: | \n")
	sb.WriteString("  Generated by GoSkills Agent\n")
	sb.WriteString("drawings:\n")
	sb.WriteString("  enabled: false\n")
	sb.WriteString("transition: slide-left\n")
	sb.WriteString("mdc: true\n")
	// Dark theme background
	sb.WriteString("background: https://picsum.photos/1920/1080?blur=4\n")
	// sb.WriteString("class: text-white\n") // Removed global class to avoid duplicates

	// Inject first slide layout
	if len(slides) > 0 {
		s0 := slides[0]
		if s0.Layout == "split-image-right" {
			sb.WriteString("layout: image-right\n")
			img := s0.Image
			if img == "" || !strings.HasPrefix(img, "http") || strings.Contains(img, "source.unsplash.com") {
				img = "https://picsum.photos/800/600?random=0"
			}
			sb.WriteString(fmt.Sprintf("image: %s\n", img))
			sb.WriteString("class: text-white\n")
		} else if s0.Layout == "title-center" {
			sb.WriteString("layout: center\n")
			sb.WriteString("class: text-center text-white\n")
		} else if s0.Layout == "two-cols" {
			sb.WriteString("layout: two-cols\n")
			sb.WriteString("class: text-white\n")
		} else {
			sb.WriteString("layout: default\n")
			sb.WriteString("class: text-white\n")
		}
	} else {
		// Fallback if no slides
		sb.WriteString("class: text-white\n")
	}
	sb.WriteString("---\n\n")

	// 2. Generate Slides
	for i, slide := range slides {
		if i > 0 {
			sb.WriteString("\n---\n")

			if slide.Layout == "split-image-right" {
				sb.WriteString("layout: image-right\n")
				img := slide.Image
				if img == "" || !strings.HasPrefix(img, "http") || strings.Contains(img, "source.unsplash.com") {
					img = fmt.Sprintf("https://picsum.photos/800/600?random=%d", i)
				}
				sb.WriteString(fmt.Sprintf("image: %s\n", img))
				sb.WriteString("class: text-white\n")
			} else if slide.Layout == "title-center" {
				sb.WriteString("layout: center\n")
				sb.WriteString("class: text-center text-white\n")
			} else if slide.Layout == "two-cols" {
				sb.WriteString("layout: two-cols\n")
				sb.WriteString("class: text-white\n")
			} else {
				sb.WriteString("layout: default\n")
				sb.WriteString("class: text-white\n")
			}
			sb.WriteString("---\n\n")
		}

		// Title with Gradient
		sb.WriteString(fmt.Sprintf("# <span class=\"bg-gradient-to-r from-cyan-400 to-purple-500 bg-clip-text text-transparent\">%s</span>\n\n", slide.Title))

		// Content Wrapper with Glassmorphism and Animation
		sb.WriteString("<div class=\"bg-black/40 backdrop-blur-md p-6 rounded-xl border border-white/10 shadow-2xl mt-4\" v-motion :initial=\"{ y: 30, opacity: 0 }\" :enter=\"{ y: 0, opacity: 1, transition: { duration: 500 } }\">\n\n")

		if slide.Layout == "two-cols" && len(slide.Content) > 1 {
			half := len(slide.Content) / 2

			sb.WriteString("<v-clicks>\n\n")
			for _, item := range slide.Content[:half] {
				sb.WriteString(fmt.Sprintf("- %s\n", item))
			}
			sb.WriteString("\n</v-clicks>\n\n")

			sb.WriteString("</div>\n") // Close left wrapper
			sb.WriteString("::right::\n")
			sb.WriteString("<div class=\"bg-black/40 backdrop-blur-md p-6 rounded-xl border border-white/10 shadow-2xl mt-4\" v-motion :initial=\"{ y: 30, opacity: 0 }\" :enter=\"{ y: 0, opacity: 1, transition: { duration: 500, delay: 200 } }\">\n\n")

			sb.WriteString("<v-clicks>\n\n")
			for _, item := range slide.Content[half:] {
				sb.WriteString(fmt.Sprintf("- %s\n", item))
			}
			sb.WriteString("\n</v-clicks>\n")
		} else {
			if len(slide.Content) > 0 {
				sb.WriteString("<v-clicks>\n\n")
				for _, item := range slide.Content {
					sb.WriteString(fmt.Sprintf("- %s\n", item))
				}
				sb.WriteString("\n</v-clicks>\n")
			}
		}

		sb.WriteString("\n</div>\n") // Close main wrapper

		// Presenter Notes
		sb.WriteString("\n<!--\n")
		sb.WriteString(fmt.Sprintf("Presenter note for slide %d: %s\n", i+1, slide.Title))
		sb.WriteString("-->\n")
	}

	return sb.String()
}

// Unused but kept for interface compatibility if needed
func (p *PPTSubagent) generateHTML(slides []Slide, filepath string) error {
	return nil
}
