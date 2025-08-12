package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/chatgpt-element-recorder/pkg/chatgpt"
	"github.com/chatgpt-element-recorder/pkg/ui"
)

// CLI represents the command line interface
type CLI struct {
	chatgpt *chatgpt.ChatGPT
	scanner *bufio.Scanner
}

// NewCLI creates a new CLI instance
func NewCLI(chatgptClient *chatgpt.ChatGPT) *CLI {
	return &CLI{
		chatgpt: chatgptClient,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// Start starts the CLI interface
func (cli *CLI) Start() error {
	cli.printWelcome()

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
		spinner.Start("Sending message and waiting for response...")

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
		if err == nil {
			ui.PrintSuccess("New chat started")
		}
		return err

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
			fmt.Print("\033[0m") // Reset colors
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
