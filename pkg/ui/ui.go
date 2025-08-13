package ui

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/term"
)

// Colors & Styles (Original constants are kept for other UI elements)
const (
	Reset     = "\033[0m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Purple    = "\033[35m"
	Cyan      = "\033[36m"
	White     = "\033[37m"
	Bold      = "\033[1m"
	Dim       = "\033[2m"
	Italic    = "\033[3m"
	Underline = "\033[4m"
	Blink     = "\033[5m"
)

// rgb struct to hold color values
type rgb struct {
	r, g, b int
}

// Function to create a rainbow gradient effect on text
func createRainbowGradient(text string) string {
	// Define the key colors for our rainbow gradient
	rainbow := []rgb{
		{r: 255, g: 0, b: 0},   // Red
		{r: 255, g: 127, b: 0}, // Orange
		{r: 255, g: 255, b: 0}, // Yellow
		{r: 0, g: 255, b: 0},   // Green
		{r: 0, g: 0, b: 255},   // Blue
		{r: 75, g: 0, b: 130},  // Indigo
		{r: 148, g: 0, b: 211}, // Violet
	}

	var builder strings.Builder
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		if len(strings.TrimSpace(line)) == 0 {
			builder.WriteString("\n")
			continue
		}

		// We use a trick to find the first and last non-space character
		// to make the gradient tighter and more vibrant.
		startIdx := strings.IndexFunc(line, func(r rune) bool { return r != ' ' })
		endIdx := strings.LastIndexFunc(line, func(r rune) bool { return r != ' ' })

		if startIdx == -1 { // Handle lines with only spaces
			builder.WriteString(line + "\n")
			continue
		}

		for i, char := range line {
			if char == ' ' || i < startIdx || i > endIdx {
				builder.WriteRune(char)
				continue
			}

			// Calculate the character's position within the visible art (0.0 to 1.0)
			pos := float64(i-startIdx) / float64(endIdx-startIdx)

			// Determine which two colors to blend
			colorPos := pos * float64(len(rainbow)-1)
			idx1 := int(colorPos)
			idx2 := idx1 + 1
			if idx2 >= len(rainbow) {
				idx2 = len(rainbow) - 1
			}

			// Calculate the blend factor between the two colors
			blend := colorPos - float64(idx1)

			// Linear interpolation for each color component (R, G, B)
			r := int(float64(rainbow[idx1].r)*(1-blend) + float64(rainbow[idx2].r)*blend)
			g := int(float64(rainbow[idx1].g)*(1-blend) + float64(rainbow[idx2].g)*blend)
			b := int(float64(rainbow[idx1].b)*(1-blend) + float64(rainbow[idx2].b)*blend)

			// Write the True Color ANSI escape code and the character
			builder.WriteString(fmt.Sprintf("\033[38;2;%d;%d;%dm%c", r, g, b, char))
		}
		builder.WriteString(Reset + "\n") // Reset color at the end of each line
	}

	return builder.String()
}

// PrintBanner prints the application banner
func PrintBanner() {
	// Clear screen for better presentation
	ClearScreen()

	// ASCII Art for CHATGPT-CLI
	// We define it as a raw string literal first
	bannerArt := `
 ██████  ██████  ████████ ███████       ██████  ███████ ██    ██ 
██       ██   ██    ██    ██            ██   ██ ██      ██    ██ 
██   ███ ██████     ██    ███████ █████ ██   ██ █████   ██    ██ 
██    ██ ██         ██         ██       ██   ██ ██       ██  ██  
 ██████  ██         ██    ███████       ██████  ███████   ████`

	// Apply the rainbow gradient effect and print it
	fmt.Print(Bold + createRainbowGradient(bannerArt) + Reset)

	// Animated-style separator
	fmt.Println()
	fmt.Print(Red + "▓")
	fmt.Print(Yellow + "▓")
	fmt.Print(Green + "▓")
	fmt.Print(Cyan + "▓")
	fmt.Print(Blue + "▓")
	fmt.Print(Purple + "▓")
	fmt.Print(strings.Repeat(White+"▓", 60))
	fmt.Print(Purple + "▓")
	fmt.Print(Blue + "▓")
	fmt.Print(Cyan + "▓")
	fmt.Print(Green + "▓")
	fmt.Print(Yellow + "▓")
	fmt.Print(Red + "▓" + Reset)
	fmt.Println()
	fmt.Println()

	// Status message
	fmt.Println(Green + Bold + "✨ Initializing GPT5-DEV Agent CLI... ✨" + Reset)
	fmt.Println("Developer : @shahirul_aiman")
	fmt.Println()
}

// PrintSuccess prints a success message
func PrintSuccess(message string) {
	fmt.Println(Green + "✅ " + message + Reset)
}

// PrintError prints an error message
func PrintError(message string) {
	fmt.Println(Red + "❌ " + message + Reset)
}

// PrintWarning prints a warning message
func PrintWarning(message string) {
	fmt.Println(Yellow + "⚠️  " + message + Reset)
}

// PrintInfo prints an info message
func PrintInfo(message string) {
	fmt.Println(Blue + "💡 " + message + Reset)
}

// PrintLoading prints a loading message
func PrintLoading(message string) {
	fmt.Println(Cyan + "⏳ " + message + Reset)
}

// ClearScreen clears the terminal screen
func ClearScreen() {
	fmt.Print("\033[2J\033[H")
}

// TypeText simulates typing effect for text output
func TypeText(text string, delay time.Duration) {
	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(delay)
	}
}

// DebugResponse prints raw response content for debugging
func DebugResponse(response string) {
	fmt.Println("\n" + Yellow + "🔍 DEBUG: Raw Response Content" + Reset)
	fmt.Println(Blue + "=" + strings.Repeat("=", 50) + Reset)

	lines := strings.Split(response, "\n")
	for i, line := range lines {
		// Show line numbers and raw content
		fmt.Printf(Dim+"%3d: "+Reset, i+1)

		// Show special characters
		displayLine := strings.ReplaceAll(line, "\t", Cyan+"[TAB]"+Reset)
		displayLine = strings.ReplaceAll(displayLine, "    ", Cyan+"[4SP]"+Reset)

		// Highlight potential code markers
		if strings.Contains(line, "```") {
			displayLine = strings.ReplaceAll(displayLine, "```", Red+"```"+Reset)
		}
		if strings.Contains(line, "Copy") {
			displayLine = strings.ReplaceAll(displayLine, "Copy", Yellow+"Copy"+Reset)
		}
		if strings.Contains(line, "Edit") {
			displayLine = strings.ReplaceAll(displayLine, "Edit", Yellow+"Edit"+Reset)
		}

		fmt.Printf("%s\n", displayLine)
	}

	fmt.Println(Blue + "=" + strings.Repeat("=", 50) + Reset)
	fmt.Printf(Green+"Total lines: %d"+Reset+"\n\n", len(lines))
}

// GetTerminalWidth gets the current terminal width
func GetTerminalWidth() int {
	// Try to get terminal width from stdout
	if width, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil && width > 20 {
		// Ensure minimum width of 40 and maximum of 120 for readability
		if width < 40 {
			return 40
		}
		if width > 120 {
			return 120
		}
		return width
	}

	// Fallback to 80 if unable to detect
	return 80
}

// Code highlighting colors
const (
	NavyBlue = "\033[48;5;17m" // Navy blue background
	CodeText = "\033[97m"      // Bright white text for code
)

// Regex patterns for fence detection
var (
	fenceStart = regexp.MustCompile(`^\s*(` + "```" + `|~~~)\s*([A-Za-z0-9+#._-]*)\s*$`)
	fenceEnd   = regexp.MustCompile(`^\s*(` + "```" + `|~~~)\s*$`)
)

// ProcessResponseWithCodeHighlight processes response text and applies code highlighting
func ProcessResponseWithCodeHighlight(text string) []ResponseLine {
	lines := strings.Split(text, "\n")
	var result []ResponseLine

	inCodeBlock := false
	codeLang := ""
	skipNext := 0

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trim := strings.TrimSpace(line)

		// Skip lines if we're in skip mode
		if skipNext > 0 {
			skipNext--
			continue
		}

		// Check if this is a language declaration line
		if isLanguageDeclaration(trim) {
			codeLang = trim
			inCodeBlock = true

			// Skip the language line and check for Copy/Edit lines
			skipCount := 1 // Skip language line

			// Check next lines for Copy/Edit and skip them too
			for j := i + 1; j < len(lines) && j < i+3; j++ {
				nextTrim := strings.TrimSpace(lines[j])
				if nextTrim == "Copy" || nextTrim == "Edit" {
					skipCount++
				} else {
					break
				}
			}

			skipNext = skipCount - 1 // -1 because we'll increment i at end of loop
			continue
		}

		// Check if we should end the code block
		if inCodeBlock {
			// End code block if we hit empty line followed by non-code content
			if trim == "" && i+1 < len(lines) {
				nextLine := strings.TrimSpace(lines[i+1])
				if nextLine != "" && !isIndentedCodeLine(lines[i+1]) && !isLanguageDeclaration(nextLine) {
					// Check if next line looks like explanation
					if isExplanationLine(nextLine) {
						inCodeBlock = false
						codeLang = ""
					}
				}
			}
		}

		// Add line with appropriate formatting
		result = append(result, ResponseLine{
			Text:     line,
			IsCode:   inCodeBlock,
			Language: codeLang,
		})
	}

	return result
}

// ResponseLine represents a line in the response with formatting info
type ResponseLine struct {
	Text     string
	IsCode   bool
	Language string
}

// parseFenceStart checks if a line starts a code fence and returns language
func parseFenceStart(line string) (ok bool, lang string) {
	m := fenceStart.FindStringSubmatch(line)
	if m == nil {
		return false, ""
	}
	lang = strings.ToLower(strings.TrimSpace(m[2]))
	return true, lang
}

// isFenceEnd checks if a line ends a code fence
func isFenceEnd(line string) bool {
	return fenceEnd.MatchString(line)
}

// isIndentedCodeLine checks if a line is indented (4 spaces or tab)
func isIndentedCodeLine(line string) bool {
	return strings.HasPrefix(line, "    ") || strings.HasPrefix(line, "\t")
}

// shouldStopIndentedBlock determines if we should end the indented code block
func shouldStopIndentedBlock(next string) bool {
	trim := strings.TrimSpace(next)
	if trim == "" {
		return true // empty line
	}
	if isIndentedCodeLine(next) {
		return false // still indented
	}
	// next line is not indented => end code block
	return true
}

// isLanguageDeclaration checks if a line is a programming language declaration
func isLanguageDeclaration(line string) bool {
	commonLanguages := []string{
		"python", "javascript", "java", "go", "rust", "c++", "c", "php",
		"ruby", "swift", "kotlin", "typescript", "html", "css", "sql",
		"bash", "shell", "powershell", "json", "xml", "yaml", "dockerfile",
		"markdown", "text", "plaintext", "output",
	}

	line = strings.ToLower(strings.TrimSpace(line))
	for _, lang := range commonLanguages {
		if line == lang {
			return true
		}
	}
	return false
}

// isExplanationLine checks if a line looks like explanation text
func isExplanationLine(line string) bool {
	// Common patterns that indicate explanation text
	explanationPatterns := []string{
		"output:", "hasil:", "contoh:", "example:", "note:", "catatan:",
		"kalau", "jika", "untuk", "ini akan", "kod ini", "awak boleh",
		"saya", "anda", "bila", "apabila", "nak saya", "boleh juga",
		"this will", "this code", "you can", "if you", "when you",
		"1.", "2.", "3.", "4.", "5.", // numbered lists
	}

	lowerLine := strings.ToLower(line)
	for _, pattern := range explanationPatterns {
		if strings.Contains(lowerLine, pattern) {
			return true
		}
	}

	return false
}

// PrintSeparator prints a separator line
func PrintSeparator() {
	fmt.Println(Blue + "─" + strings.Repeat("─", 50) + Reset)
}

// PrintWelcome prints the welcome message
func PrintWelcome() {
	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = "Unknown"
	}

	fmt.Println(Purple + "💡 Commands:" + Reset)
	fmt.Println("  " + Cyan + "/help" + Reset + "    - Show help")
	fmt.Println("  " + Cyan + "/new" + Reset + "     - Start new chat")
	fmt.Println("  " + Cyan + "/history" + Reset + " - Show chat history")
	fmt.Println("  " + Cyan + "/open <id>" + Reset + " - Open specific chat")
	fmt.Println("  " + Cyan + "/quit" + Reset + "    - Exit")
	fmt.Println()
	fmt.Println(Green + "💬 Just type your message to chat with ChatGPT!" + Reset)
	fmt.Println("Model: " + Cyan + "GPT5" + Reset)
	fmt.Println(Dim + "📁 Working in: " + currentDir + Reset)
}
