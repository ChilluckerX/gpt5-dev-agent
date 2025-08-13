package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/chatgpt-element-recorder/pkg/agent"
	"github.com/chatgpt-element-recorder/pkg/browser"
	"github.com/chatgpt-element-recorder/pkg/chatgpt"
	"github.com/chatgpt-element-recorder/pkg/config"
	"github.com/chatgpt-element-recorder/pkg/ui"
)

// CLI represents the command line interface
type CLI struct {
	chatgpt *chatgpt.ChatGPT
	scanner *bufio.Scanner
	agent   *agent.Agent // Agent system integration
	config  *config.DynamicConfig
}

// NewCLI creates a new CLI instance
func NewCLI(chatgptClient *chatgpt.ChatGPT) *CLI {
	// Load dynamic configuration
	config, err := config.LoadDynamicConfig()
	if err != nil {
		// Use default config if loading fails
		ui.PrintWarning("Could not load configuration, using defaults")
	}
	
	// Create agent instance
	agentInstance, err := agent.NewAgent(chatgptClient)
	if err != nil {
		// Continue without agent if creation fails
		ui.PrintWarning("Could not initialize agent system")
		agentInstance = nil
	}
	
	return &CLI{
		chatgpt: chatgptClient,
		scanner: bufio.NewScanner(os.Stdin),
		agent:   agentInstance,
		config:  config,
	}
}

// Start starts the CLI interface
func (cli *CLI) Start() error {
	cli.printWelcome()
	
	// Auto-send system prompt for initial context
	if err := cli.sendSystemPromptForNewChat(); err != nil {
		ui.PrintWarning("Could not establish initial project context")
	}

	for {
		fmt.Print("\n> ")

		if !cli.scanner.Scan() {
			break
		}

		input := strings.TrimSpace(cli.scanner.Text())
		if input == "" {
			continue
		}

		// Handle commands
		if strings.HasPrefix(input, "/") {
			if err := cli.handleCommand(input); err != nil {
				ui.PrintError(fmt.Sprintf("Error: %v", err))
			}
			continue
		}

		// Send message to ChatGPT with spinner
		spinner := ui.NewSpinner()
		spinner.Start("")

		response, err := cli.chatgpt.SendMessage(input)
		spinner.Stop()

		if err != nil {
			ui.PrintError(fmt.Sprintf("Error sending message: %v", err))
			continue
		}

		cli.printResponse(response)
	}

	return nil
}

// handleCommand handles CLI commands
func (cli *CLI) handleCommand(command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil
	}

	cmd := parts[0]

	switch cmd {
	case "/help", "/h":
		cli.printHelp()

	case "/new", "/n":
		spinner := ui.NewSquareSpinner()
		spinner.Start("Starting new chat...")
		err := cli.chatgpt.StartNewChat()
		spinner.Stop()
		
		if err != nil {
			return err
		}
		
		ui.PrintSuccess("New chat started")
		
		// Auto-send system prompt with project context
		return cli.sendSystemPromptForNewChat()

	case "/history", "/hist":
		return cli.showHistory()

	case "/open", "/o":
		if len(parts) < 2 {
			fmt.Println("❌ Usage: /open <chat_id_or_number>")
			return nil
		}
		return cli.openChat(parts[1])

	case "/quit", "/q", "/exit":
		ui.PrintSuccess("Goodbye!")
		os.Exit(0)

	case "/clear", "/cls":
		ui.ClearScreen()

	case "/cookies", "/c":
		if len(parts) < 2 {
			fmt.Println("❌ Usage: /cookies <validate|clean|status>")
			return nil
		}
		return cli.handleCookies(parts[1])

	default:
		fmt.Printf("❌ Unknown command: %s\n", cmd)
		fmt.Println("💡 Type /help for available commands")
	}

	return nil
}

// showHistory shows chat history
func (cli *CLI) showHistory() error {
	spinner := ui.NewSquareSpinner()
	spinner.Start("Loading chat history...")

	history, err := cli.chatgpt.GetChatHistory()
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to get history: %v", err)
	}

	if len(history) == 0 {
		ui.PrintWarning("No chat history found")
		return nil
	}

	fmt.Println("\n📜 Recent Chat History:")
	ui.PrintSeparator()

	for i, item := range history {
		fmt.Printf("%d. %s\n", i+1, item.Title)
		fmt.Printf("   ID: %s\n", item.ID)
		fmt.Println()
	}

	ui.PrintInfo("Use '/open <number>' or '/open <chat_id>' to open a chat")
	return nil
}

// openChat opens a specific chat
func (cli *CLI) openChat(identifier string) error {
	// Check if it's a number (history index)
	if num, err := strconv.Atoi(identifier); err == nil {
		// Get history and open by index
		history, err := cli.chatgpt.GetChatHistory()
		if err != nil {
			return fmt.Errorf("failed to get history: %v", err)
		}

		if num < 1 || num > len(history) {
			return fmt.Errorf("invalid history number: %d (available: 1-%d)", num, len(history))
		}

		chatID := history[num-1].ID
		fmt.Printf("📂 Opening chat: %s\n", history[num-1].Title)
		return cli.chatgpt.OpenChat(chatID)
	}

	// Otherwise treat as chat ID
	fmt.Printf("📂 Opening chat ID: %s\n", identifier)
	return cli.chatgpt.OpenChat(identifier)
}

// printWelcome prints welcome message
func (cli *CLI) printWelcome() {
	ui.PrintWelcome()
}

// printHelp prints help information
func (cli *CLI) printHelp() {
	fmt.Println("\n📖 ChatGPT CLI Help")
	fmt.Println("=" + strings.Repeat("=", 30))
	fmt.Println()
	fmt.Println("🔧 Commands:")
	fmt.Println("  /help, /h           - Show this help")
	fmt.Println("  /new, /n            - Start a new chat")
	fmt.Println("  /history, /hist     - Show recent chat history")
	fmt.Println("  /open <id>, /o <id> - Open chat by ID or number")
	fmt.Println("  /clear, /cls        - Clear screen")
	fmt.Println("  /quit, /q, /exit    - Exit the CLI")
	fmt.Println()
	fmt.Println("💬 Usage:")
	fmt.Println("  - Type any message to send to ChatGPT")
	fmt.Println("  - Use /new to start fresh conversation")
	fmt.Println("  - Use /history to see previous chats")
	fmt.Println("  - Use /open 1 to open first chat from history")
	fmt.Println("  - Use /open <chat-id> to open specific chat")
	fmt.Println()
	fmt.Println("🎯 Examples:")
	fmt.Println("  Hello, how are you?")
	fmt.Println("  /new")
	fmt.Println("  /history")
	fmt.Println("  /open 1")
	fmt.Println("  /open 689916e6-3df0-8331-8eb6-e6f0c648cea4")
}

// printResponse prints ChatGPT response with formatting and typing effect
func (cli *CLI) printResponse(response string) {
	// Simple clean formatting without aggressive code detection
	response = strings.TrimSpace(response)

	// Remove "Thought for Xs" prefix if present
	if strings.HasPrefix(response, "Thought for") {
		lines := strings.Split(response, "\n")
		if len(lines) > 1 {
			response = strings.Join(lines[1:], "\n")
		}
	}

	fmt.Println()

	// Calculate responsive box width based on terminal size
	boxWidth := ui.GetTerminalWidth()
	headerText := "  Response   "
	headerLine := headerText + strings.Repeat("─", boxWidth-len(headerText)-2)

	// Print the header line immediately (no typing effect for border)
	fmt.Print("\033[92m╭" + headerLine + "╮\033[0m\n")

	// Process response with code highlighting
	responseLines := ui.ProcessResponseWithCodeHighlight(response)

	for _, responseLine := range responseLines {
		// Print border immediately
		fmt.Print("\033[92m│   \033[0m")

		// Apply code highlighting if this is a code line
		if responseLine.IsCode {
			// Navy blue background with white text for code
			fmt.Print(ui.NavyBlue + ui.CodeText)
			ui.TypeText(responseLine.Text, 20*time.Millisecond) // Slightly faster for code
			fmt.Print("\033[0m")                                // Reset colors
		} else {
			// Normal text with typing effect
			ui.TypeText(responseLine.Text, 30*time.Millisecond)
		}

		// Calculate padding to fill the line
		padding := boxWidth - len(responseLine.Text) - 5 // 5 = "│   " + "│"
		if padding > 0 {
			if responseLine.IsCode {
				// Continue navy blue background for padding
				fmt.Print(ui.NavyBlue + strings.Repeat(" ", padding) + "\033[0m")
			} else {
				fmt.Print(strings.Repeat(" ", padding))
			}
		}
		fmt.Print("\033[92m│\033[0m\n")
	}

	// Print the bottom border immediately (no typing effect)
	fmt.Print("\033[92m╰" + strings.Repeat("─", boxWidth-2) + "╯\033[0m\n")
}

// clearScreen clears the terminal screen (deprecated - use ui.ClearScreen)
func (cli *CLI) clearScreen() {
	ui.ClearScreen()
}

// generateSystemPrompt creates a system prompt with project context
func (cli *CLI) generateSystemPrompt() string {
	currentDir, _ := os.Getwd()
	
	// Analyze project structure
	projectInfo := cli.analyzeProjectStructure()
	
	systemPrompt := fmt.Sprintf(`You are GPT5-DEV, a friendly and expert software development assistant. You're helping a developer who is currently working in the directory: %s

Project Analysis:
%s

Your role:
- Act as a knowledgeable coding assistant and mentor
- Provide helpful suggestions based on the project structure you see
- Be conversational and friendly, like a senior developer colleague
- Ask intelligent follow-up questions about their work
- Offer specific help based on the technologies and files you observe

Please greet the user by acknowledging what you see in their project and ask how you can help them today. Be specific about what you notice in their codebase.`, currentDir, projectInfo)

	return systemPrompt
}

// analyzeProjectStructure analyzes the current directory and returns project info
func (cli *CLI) analyzeProjectStructure() string {
	var analysis strings.Builder
	currentDir, _ := os.Getwd()
	
	// Get directory name
	projectName := filepath.Base(currentDir)
	analysis.WriteString(fmt.Sprintf("Project: %s\n", projectName))
	
	// Analyze files and folders
	entries, err := os.ReadDir(".")
	if err != nil {
		analysis.WriteString("Unable to read directory structure\n")
		return analysis.String()
	}
	
	var files []string
	var folders []string
	var configFiles []string
	var codeFiles []string
	
	for _, entry := range entries {
		name := entry.Name()
		
		// Skip hidden files and common ignore patterns
		if strings.HasPrefix(name, ".") && name != ".env" && name != ".gitignore" {
			continue
		}
		
		if entry.IsDir() {
			folders = append(folders, name)
		} else {
			files = append(files, name)
			
			// Categorize files
			ext := strings.ToLower(filepath.Ext(name))
			switch {
			case name == "go.mod" || name == "package.json" || name == "requirements.txt" || name == "Cargo.toml" || name == "pom.xml":
				configFiles = append(configFiles, name)
			case ext == ".go" || ext == ".py" || ext == ".js" || ext == ".ts" || ext == ".java" || ext == ".rs" || ext == ".cpp" || ext == ".c":
				codeFiles = append(codeFiles, name)
			case name == "README.md" || name == "Dockerfile" || name == ".gitignore":
				configFiles = append(configFiles, name)
			}
		}
	}
	
	// Build analysis
	if len(configFiles) > 0 {
		analysis.WriteString(fmt.Sprintf("Config files: %s\n", strings.Join(configFiles, ", ")))
	}
	
	if len(codeFiles) > 0 {
		analysis.WriteString(fmt.Sprintf("Code files: %s\n", strings.Join(codeFiles, ", ")))
	}
	
	if len(folders) > 0 {
		analysis.WriteString(fmt.Sprintf("Directories: %s\n", strings.Join(folders, ", ")))
	}
	
	// Detect project type
	projectType := cli.detectProjectType(configFiles, codeFiles)
	if projectType != "" {
		analysis.WriteString(fmt.Sprintf("Detected: %s project\n", projectType))
	}
	
	return analysis.String()
}

// detectProjectType tries to determine the project type based on files
func (cli *CLI) detectProjectType(configFiles, codeFiles []string) string {
	// Check config files first
	for _, file := range configFiles {
		switch file {
		case "go.mod":
			return "Go"
		case "package.json":
			return "Node.js/JavaScript"
		case "requirements.txt", "setup.py":
			return "Python"
		case "Cargo.toml":
			return "Rust"
		case "pom.xml":
			return "Java/Maven"
		case "Dockerfile":
			return "Docker"
		}
	}
	
	// Check code files
	for _, file := range codeFiles {
		ext := strings.ToLower(filepath.Ext(file))
		switch ext {
		case ".go":
			return "Go"
		case ".py":
			return "Python"
		case ".js", ".ts":
			return "JavaScript/TypeScript"
		case ".java":
			return "Java"
		case ".rs":
			return "Rust"
		case ".cpp", ".c":
			return "C/C++"
		}
	}
	
	return ""
}

// sendSystemPromptForNewChat sends system prompt when starting new chat
func (cli *CLI) sendSystemPromptForNewChat() error {
	systemPrompt := cli.generateSystemPrompt()
	
	spinner := ui.NewSquareSpinner()
	spinner.Start("Analyzing project and setting up context...")
	
	// Send system prompt
	_, err := cli.chatgpt.SendMessage(systemPrompt)
	spinner.Stop()
	
	if err != nil {
		ui.PrintWarning("Could not set up project context")
		return err
	}
	
	ui.PrintSuccess("Project context established! 🎯")
	return nil
}


// handleCookies handles cookie management commands
func (cli *CLI) handleCookies(action string) error {
	cookieManager := browser.NewCookieManager()
	
	switch strings.ToLower(action) {
	case "validate", "v":
		spinner := ui.NewSquareSpinner()
		spinner.Start("Validating cookies...")
		err := cookieManager.EnsureCookiesFile()
		spinner.Stop()
		if err != nil {
			ui.PrintError(fmt.Sprintf("Cookie validation failed: %v", err))
		} else {
			ui.PrintSuccess("Cookies validation completed!")
		}
		return nil
		
	case "clean", "c":
		spinner := ui.NewSquareSpinner()
		spinner.Start("Cleaning expired cookies...")
		err := cookieManager.CleanExpiredCookies()
		spinner.Stop()
		if err != nil {
			ui.PrintError(fmt.Sprintf("Failed to clean cookies: %v", err))
		} else {
			ui.PrintSuccess("Cookie cleanup completed!")
		}
		return nil
		
	case "status", "s":
		fmt.Println("\n🍪 Cookie Status:")
		ui.PrintSeparator()
		fmt.Printf("📁 Cookies file: %s\n", cookieManager.GetCookiesPath())
		
		if _, err := os.Stat(cookieManager.GetCookiesPath()); os.IsNotExist(err) {
			fmt.Println("❌ Cookies file does not exist")
			fmt.Println("💡 Run \"/cookies validate\" to create it")
		} else {
			cookies, err := cookieManager.LoadCookies()
			if err != nil {
				ui.PrintError(fmt.Sprintf("Failed to load cookies: %v", err))
			} else if len(cookies) == 0 {
				fmt.Println("❌ No cookies found")
				fmt.Println("💡 You may need to login to ChatGPT manually")
			} else {
				fmt.Printf("📊 Total cookies: %d\n", len(cookies))
				fmt.Println("✅ Cookies file is valid")
			}
		}
		ui.PrintSeparator()
		return nil
		
	default:
		fmt.Printf("❌ Unknown cookie action: %s\n", action)
		fmt.Println("💡 Available actions: validate, clean, status")
		return nil
	}
}
