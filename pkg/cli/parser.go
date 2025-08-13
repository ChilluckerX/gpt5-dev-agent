package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/chatgpt-element-recorder/pkg/agent"
)

// CLIArgs represents parsed command line arguments
type CLIArgs struct {
	Mode        string
	Query       string
	Interactive bool
	Config      string
	Help        bool
	Version     bool
	Debug       bool
	NoContext   bool
	OutputFile  string
}

// ParseArgs parses command line arguments similar to sengpt
func ParseArgs() (*CLIArgs, error) {
	args := &CLIArgs{}
	
	// Define flags
	flag.StringVar(&args.Mode, "mode", "interactive", "Operation mode: interactive, query, auto, context")
	flag.StringVar(&args.Mode, "m", "interactive", "Operation mode (short)")
	flag.StringVar(&args.Query, "query", "", "Single query to execute (for query mode)")
	flag.StringVar(&args.Query, "q", "", "Single query (short)")
	flag.BoolVar(&args.Interactive, "interactive", false, "Force interactive mode")
	flag.BoolVar(&args.Interactive, "i", false, "Force interactive mode (short)")
	flag.StringVar(&args.Config, "config", "", "Path to config file")
	flag.StringVar(&args.Config, "c", "", "Path to config file (short)")
	flag.BoolVar(&args.Help, "help", false, "Show help message")
	flag.BoolVar(&args.Help, "h", false, "Show help (short)")
	flag.BoolVar(&args.Version, "version", false, "Show version information")
	flag.BoolVar(&args.Version, "v", false, "Show version (short)")
	flag.BoolVar(&args.Debug, "debug", false, "Enable debug mode")
	flag.BoolVar(&args.Debug, "d", false, "Enable debug mode (short)")
	flag.BoolVar(&args.NoContext, "no-context", false, "Disable project context analysis")
	flag.StringVar(&args.OutputFile, "output", "", "Output file for responses")
	flag.StringVar(&args.OutputFile, "o", "", "Output file (short)")
	
	// Custom usage function
	flag.Usage = func() {
		printUsage()
	}
	
	flag.Parse()
	
	// Handle remaining arguments as query if no -q flag
	if args.Query == "" && len(flag.Args()) > 0 {
		args.Query = strings.Join(flag.Args(), " ")
	}
	
	// Validate arguments
	if err := validateArgs(args); err != nil {
		return nil, err
	}
	
	return args, nil
}

// validateArgs validates the parsed arguments
func validateArgs(args *CLIArgs) error {
	// Validate mode
	validModes := []string{"interactive", "query", "auto", "context"}
	isValidMode := false
	for _, mode := range validModes {
		if args.Mode == mode {
			isValidMode = true
			break
		}
	}
	if !isValidMode {
		return fmt.Errorf("invalid mode: %s. Valid modes: %s", args.Mode, strings.Join(validModes, ", "))
	}
	
	// Query mode requires a query
	if args.Mode == "query" && args.Query == "" {
		return fmt.Errorf("query mode requires a query (-q or --query)")
	}
	
	return nil
}

// printUsage prints the usage information
func printUsage() {
	fmt.Fprintf(os.Stderr, `ChatGPT CLI Agent - Intelligent development assistant

Usage:
  %s [OPTIONS] [QUERY]

Modes:
  interactive    Interactive chat mode (default)
  query         Single query mode
  auto          Autonomous task execution mode
  context       Context-aware assistance mode

Options:
  -m, --mode MODE        Operation mode (interactive, query, auto, context)
  -q, --query QUERY      Single query to execute
  -i, --interactive      Force interactive mode
  -c, --config FILE      Path to config file
  -o, --output FILE      Output file for responses
  --no-context          Disable project context analysis
  -d, --debug           Enable debug mode
  -h, --help            Show this help message
  -v, --version         Show version information

Examples:
  %s                                    # Start interactive mode
  %s -q "explain this code"             # Single query
  %s -m context "help with Go project" # Context-aware mode
  %s -i --no-context                   # Interactive without context
  %s -o output.txt -q "generate docs"  # Save response to file

For more information, visit: https://github.com/your-repo/chatgpt-cli
`, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
}

// ExecuteWithArgs executes the CLI with parsed arguments
func ExecuteWithArgs(args *CLIArgs, cliInstance *CLI) error {
	// Handle special flags first
	if args.Help {
		printUsage()
		return nil
	}
	
	if args.Version {
		printVersion()
		return nil
	}
	
	// Load custom config if specified
	if args.Config != "" {
		if err := loadCustomConfig(args.Config); err != nil {
			return fmt.Errorf("failed to load config: %v", err)
		}
	}
	
	// Create agent and set mode
	agentInstance, err := agent.NewAgent(cliInstance.chatgpt)
	if err != nil {
		return fmt.Errorf("failed to create agent: %v", err)
	}
	
	// Set agent mode
	switch args.Mode {
	case "interactive":
		agentInstance.SetMode(agent.InteractiveMode)
	case "query":
		agentInstance.SetMode(agent.QueryMode)
	case "auto":
		agentInstance.SetMode(agent.AutoMode)
	case "context":
		agentInstance.SetMode(agent.ContextMode)
	}
	
	// Initialize session unless disabled
	if !args.NoContext {
		if err := agentInstance.InitializeSession(); err != nil {
			// Don't fail, just warn
			fmt.Printf("Warning: Could not initialize project context: %v\n", err)
		}
	}
	
	// Execute based on mode
	switch args.Mode {
	case "query":
		return executeQueryMode(agentInstance, args)
	case "interactive":
		return executeInteractiveMode(cliInstance, agentInstance, args)
	case "auto":
		return executeAutoMode(agentInstance, args)
	case "context":
		return executeContextMode(agentInstance, args)
	default:
		return executeInteractiveMode(cliInstance, agentInstance, args)
	}
}

// executeQueryMode executes a single query
func executeQueryMode(agent *agent.Agent, args *CLIArgs) error {
	response, err := agent.ProcessMessage(args.Query)
	if err != nil {
		return fmt.Errorf("query failed: %v", err)
	}
	
	// Output response
	if args.OutputFile != "" {
		return writeToFile(args.OutputFile, response)
	}
	
	fmt.Println(response)
	return nil
}

// executeInteractiveMode starts interactive mode
func executeInteractiveMode(cliInstance *CLI, agentInstance *agent.Agent, args *CLIArgs) error {
	// Set the agent in CLI instance
	cliInstance.agent = agentInstance
	
	// Start interactive mode
	return cliInstance.Start()
}

// executeAutoMode executes autonomous mode
func executeAutoMode(agent *agent.Agent, args *CLIArgs) error {
	// Auto mode implementation would go here
	// For now, fall back to query mode
	if args.Query != "" {
		return executeQueryMode(agent, args)
	}
	
	fmt.Println("Auto mode: Please specify a task with -q or --query")
	return nil
}

// executeContextMode executes context-aware mode
func executeContextMode(agent *agent.Agent, args *CLIArgs) error {
	// Context mode could provide enhanced project analysis
	if args.Query != "" {
		return executeQueryMode(agent, args)
	}
	
	// Show project context
	context := agent.GetProjectContext()
	if context != nil {
		fmt.Println("Project Context:")
		fmt.Println(context.GetProjectInfo())
	}
	
	return nil
}

// printVersion prints version information
func printVersion() {
	fmt.Println("ChatGPT CLI Agent v1.0.0")
	fmt.Println("Intelligent development assistant")
	fmt.Println("Built with Go")
}

// loadCustomConfig loads a custom configuration file
func loadCustomConfig(configPath string) error {
	// This would load a custom config file
	// For now, just validate the file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", configPath)
	}
	return nil
}

// writeToFile writes content to a file
func writeToFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

// GetModeFromString converts string to AgentMode
func GetModeFromString(mode string) agent.AgentMode {
	switch strings.ToLower(mode) {
	case "interactive":
		return agent.InteractiveMode
	case "query":
		return agent.QueryMode
	case "auto":
		return agent.AutoMode
	case "context":
		return agent.ContextMode
	default:
		return agent.InteractiveMode
	}
}