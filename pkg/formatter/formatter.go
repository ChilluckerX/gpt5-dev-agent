package formatter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/chatgpt-element-recorder/pkg/ui"
)

// FormatResponse formats ChatGPT responses with proper code highlighting and structure
func FormatResponse(response string) string {
	// Clean up the response first
	response = strings.TrimSpace(response)
	
	// Remove "Thought for Xs" prefix if present
	thoughtRegex := regexp.MustCompile(`^Thought for \d+s\s*`)
	response = thoughtRegex.ReplaceAllString(response, "")
	
	// Only detect code blocks if they have VERY clear indicators
	if strings.Contains(response, "```") || 
	   (strings.Contains(response, "python") && (strings.Contains(response, "def ") || strings.Contains(response, "import ") || strings.Contains(response, "print("))) ||
	   (strings.Contains(response, "javascript") && strings.Contains(response, "function")) {
		response = formatCodeBlocks(response)
	}
	
	// Skip inline code formatting for now to avoid false positives
	// response = formatInlineCode(response)
	
	// Add proper line breaks for readability
	response = formatParagraphs(response)
	
	return response
}

// formatCodeBlocks detects and formats multi-line code blocks
func formatCodeBlocks(text string) string {
	// Pattern for code blocks (python, javascript, etc.)
	codeBlockRegex := regexp.MustCompile(`(?i)(python|javascript|js|go|java|c\+\+|cpp|html|css|sql|bash|shell|json|yaml|xml)(?:Copy|Edit)?\s*([\s\S]+?)(?:\n\n|\n[A-Z]|$)`)
	
	return codeBlockRegex.ReplaceAllStringFunc(text, func(match string) string {
		// Extract language and code
		parts := codeBlockRegex.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}
		
		language := strings.ToLower(parts[1])
		code := strings.TrimSpace(parts[2])
		
		// Format the code block
		formatted := "\n" + ui.Blue + "📄 " + strings.ToUpper(language) + " Code:" + ui.Reset + "\n"
		formatted += ui.Cyan + "```" + language + ui.Reset + "\n"
		formatted += formatCodeContent(code) + "\n"
		formatted += ui.Cyan + "```" + ui.Reset + "\n"
		
		return formatted
	})
}

// formatCodeContent formats the actual code content
func formatCodeContent(code string) string {
	lines := strings.Split(code, "\n")
	var formatted []string
	
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			formatted = append(formatted, "")
			continue
		}
		
		// Add line numbers for longer code blocks
		if len(lines) > 5 {
			lineNum := ui.Yellow + sprintf("%2d", i+1) + ui.Reset + " │ "
			formatted = append(formatted, lineNum+ui.Green+line+ui.Reset)
		} else {
			formatted = append(formatted, ui.Green+line+ui.Reset)
		}
	}
	
	return strings.Join(formatted, "\n")
}

// formatInlineCode formats inline code snippets
func formatInlineCode(text string) string {
	// Pattern for inline code (words that look like code)
	inlineCodeRegex := regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*\.[a-zA-Z_][a-zA-Z0-9_]*|[a-zA-Z_][a-zA-Z0-9_]*\(\)|[A-Z_][A-Z0-9_]{2,})\b`)
	
	return inlineCodeRegex.ReplaceAllStringFunc(text, func(match string) string {
		// Don't format if it's already in a code block
		return ui.Cyan + "`" + match + "`" + ui.Reset
	})
}

// formatParagraphs adds proper spacing between paragraphs
func formatParagraphs(text string) string {
	// Split into paragraphs
	paragraphs := strings.Split(text, "\n\n")
	var formatted []string
	
	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}
		
		// Check if this is a list item
		if strings.HasPrefix(para, "-") || strings.HasPrefix(para, "*") || strings.HasPrefix(para, "•") {
			formatted = append(formatted, formatList(para))
		} else if strings.HasPrefix(para, "#") {
			// Format headers
			formatted = append(formatted, formatHeader(para))
		} else {
			// Regular paragraph
			formatted = append(formatted, para)
		}
	}
	
	return strings.Join(formatted, "\n\n")
}

// formatList formats list items
func formatList(text string) string {
	lines := strings.Split(text, "\n")
	var formatted []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") || strings.HasPrefix(line, "•") {
			// Format list item
			content := strings.TrimSpace(line[1:])
			formatted = append(formatted, ui.Yellow+"  •"+ui.Reset+" "+content)
		} else {
			formatted = append(formatted, "    "+line)
		}
	}
	
	return strings.Join(formatted, "\n")
}

// formatHeader formats markdown-style headers
func formatHeader(text string) string {
	if strings.HasPrefix(text, "###") {
		return ui.Blue + ui.Bold + text + ui.Reset
	} else if strings.HasPrefix(text, "##") {
		return ui.Purple + ui.Bold + text + ui.Reset
	} else if strings.HasPrefix(text, "#") {
		return ui.Cyan + ui.Bold + text + ui.Reset
	}
	return text
}

// Helper function for sprintf (simple implementation)
func sprintf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}