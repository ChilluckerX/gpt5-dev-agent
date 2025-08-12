package ui

import (
	"fmt"
	"os"
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
		{r: 255, g: 0, b: 0},     // Red
		{r: 255, g: 127, b: 0},    // Orange
		{r: 255, g: 255, b: 0},    // Yellow
		{r: 0, g: 255, b: 0},     // Green
		{r: 0, g: 0, b: 255},     // Blue
		{r: 75, g: 0, b: 130},    // Indigo
		{r: 148, g: 0, b: 211},   // Violet
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
 ██████╗██╗  ██╗ █████╗ ████████╗ ██████╗ ██████╗ ████████╗      ██████╗██╗     ██╗
██╔════╝██║  ██║██╔══██╗╚══██╔══╝██╔════╝ ██╔══██╗╚══██╔══╝     ██╔════╝██║     ██║
██║     ███████║███████║   ██║   ██║  ███╗██████╔╝   ██║        ██║     ██║     ██║
██║     ██╔══██║██╔══██║   ██║   ██║   ██║██╔═══╝    ██║        ██║     ██║     ██║
╚██████╗██║  ██║██║  ██║   ██║   ╚██████╔╝██║        ██║███████╗╚██████╗███████╗██║
 ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═╝        ╚═╝╚══════╝ ╚═════╝╚══════╝╚═╝`

	// Apply the rainbow gradient effect and print it
	fmt.Print(Bold + createRainbowGradient(bannerArt) + Reset)
	
	// The rest of your banner remains the same
	fmt.Println(Purple + Bold + `
                    ╔═══════════════════════════════════════════════════════╗` + Reset)
	fmt.Println(Purple + `                    ║` + Yellow + Bold + `          🤖 AI-Powered Terminal Interface          ` + Purple + `║` + Reset)
	fmt.Println(Purple + `                    ║` + Green + `             🚀 Go Edition - High Performance        ` + Purple + `║` + Reset)
	fmt.Println(Purple + `                    ║` + Cyan + `              🇲🇾 Advanced Scraper Technology         ` + Purple + `║` + Reset)
	fmt.Println(Purple + Bold + `                    ╚═══════════════════════════════════════════════════════╝` + Reset)
	
	// Animated-style separator
	fmt.Println()
	fmt.Print(Red + "▓")
	fmt.Print(Yellow + "▓")
	fmt.Print(Green + "▓")
	fmt.Print(Cyan + "▓")
	fmt.Print(Blue + "▓")
	fmt.Print(Purple + "▓")
	fmt.Print(strings.Repeat(White + "▓", 60))
	fmt.Print(Purple + "▓")
	fmt.Print(Blue + "▓")
	fmt.Print(Cyan + "▓")
	fmt.Print(Green + "▓")
	fmt.Print(Yellow + "▓")
	fmt.Print(Red + "▓" + Reset)
	fmt.Println()
	fmt.Println()
	
	// Status message
	fmt.Println(Green + Bold + "                           ✨ Initializing ChatGPT CLI... ✨" + Reset)
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
	NavyBlue = "\033[48;5;17m"  // Navy blue background
	CodeText = "\033[97m"       // Bright white text for code
)

// ProcessResponseWithCodeHighlight processes response text and applies code highlighting
func ProcessResponseWithCodeHighlight(text string) []ResponseLine {
	lines := strings.Split(text, "\n")
	var result []ResponseLine
	
	inCodeBlock := false
	codeLanguage := ""
	
	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		
		// Check if this line starts a code block
		if isCodeLanguageLine(trimmedLine) {
			codeLanguage = trimmedLine
			inCodeBlock = true
			// Skip the language line, Copy, and Edit lines
			continue
		}
		
		// Skip "Copy" and "Edit" lines after language declaration
		if inCodeBlock && (trimmedLine == "Copy" || trimmedLine == "Edit") {
			continue
		}
		
		// Check if we should end the code block
		if inCodeBlock && shouldEndCodeBlock(line, lines, i) {
			inCodeBlock = false
			codeLanguage = ""
		}
		
		// Add the line with appropriate formatting
		result = append(result, ResponseLine{
			Text:     line,
			IsCode:   inCodeBlock,
			Language: codeLanguage,
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

// isCodeLanguageLine checks if a line indicates the start of a code block
func isCodeLanguageLine(line string) bool {
	commonLanguages := []string{
		"python", "javascript", "java", "go", "rust", "c++", "c", "php", 
		"ruby", "swift", "kotlin", "typescript", "html", "css", "sql",
		"bash", "shell", "powershell", "json", "xml", "yaml", "dockerfile",
	}
	
	line = strings.ToLower(strings.TrimSpace(line))
	for _, lang := range commonLanguages {
		if line == lang {
			return true
		}
	}
	return false
}

// shouldEndCodeBlock determines if we should end the current code block
func shouldEndCodeBlock(currentLine string, allLines []string, currentIndex int) bool {
	trimmed := strings.TrimSpace(currentLine)
	
	// Empty line followed by non-indented text usually ends code
	if trimmed == "" && currentIndex+1 < len(allLines) {
		nextLine := strings.TrimSpace(allLines[currentIndex+1])
		if nextLine != "" && !strings.HasPrefix(allLines[currentIndex+1], "    ") && !strings.HasPrefix(allLines[currentIndex+1], "\t") {
			// Check if next line looks like explanation text
			if isExplanationText(nextLine) {
				return true
			}
		}
	}
	
	return false
}

// isExplanationText checks if a line looks like explanation rather than code
func isExplanationText(line string) bool {
	// Common patterns that indicate explanation text
	explanationStarters := []string{
		"kalau", "jika", "untuk", "ini akan", "kod ini", "awak boleh", 
		"saya", "anda", "bila", "apabila", "contoh", "example",
		"this will", "this code", "you can", "if you", "when you",
	}
	
	lowerLine := strings.ToLower(line)
	for _, starter := range explanationStarters {
		if strings.Contains(lowerLine, starter) {
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
	fmt.Println(Purple + "💡 Commands:" + Reset)
	fmt.Println("  " + Cyan + "/help" + Reset + "    - Show help")
	fmt.Println("  " + Cyan + "/new" + Reset + "     - Start new chat")
	fmt.Println("  " + Cyan + "/history" + Reset + " - Show chat history")
	fmt.Println("  " + Cyan + "/open <id>" + Reset + " - Open specific chat")
	fmt.Println("  " + Cyan + "/quit" + Reset + "    - Exit")
	fmt.Println()
	fmt.Println(Green + "💬 Just type your message to chat with ChatGPT!" + Reset)
}