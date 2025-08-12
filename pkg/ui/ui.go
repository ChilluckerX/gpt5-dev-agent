package ui

import (
	"fmt"
	"strings"
)

// Colors
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

// PrintBanner prints the application banner
func PrintBanner() {
	// Clear screen for better presentation
	ClearScreen()
	
	// ASCII Art for CHATGPT-CLI with colors
	fmt.Println(Cyan + Bold + `
 ██████╗██╗  ██╗ █████╗ ████████╗ ██████╗ ██████╗ ████████╗      ██████╗██╗     ██╗
██╔════╝██║  ██║██╔══██╗╚══██╔══╝██╔════╝ ██╔══██╗╚══██╔══╝     ██╔════╝██║     ██║
██║     ███████║███████║   ██║   ██║  ███╗██████╔╝   ██║        ██║     ██║     ██║
██║     ██╔══██║██╔══██║   ██║   ██║   ██║██╔═══╝    ██║        ██║     ██║     ██║
╚██████╗██║  ██║██║  ██║   ██║   ╚██████╔╝██║        ██║███████╗╚██████╗███████╗██║
 ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═╝        ╚═╝╚══════╝ ╚═════╝╚══════╝╚═╝` + Reset)
	
	// Gradient effect with different colors
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