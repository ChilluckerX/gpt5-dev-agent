package agent

import (
	"fmt"
	"strings"

	"github.com/chatgpt-element-recorder/pkg/chatgpt"
	"github.com/chatgpt-element-recorder/pkg/config"
	"github.com/chatgpt-element-recorder/pkg/ui"
)

// Agent represents the main agent system
type Agent struct {
	chatgpt   *chatgpt.ChatGPT
	config    *config.DynamicConfig
	mode      AgentMode
	context   *ProjectContext
	fileOps   *FileOperations
}

// AgentMode represents different operation modes
type AgentMode string

const (
	InteractiveMode AgentMode = "interactive"
	QueryMode       AgentMode = "query"
	AutoMode        AgentMode = "auto"
	ContextMode     AgentMode = "context"
)

// NewAgent creates a new agent instance
func NewAgent(chatgptClient *chatgpt.ChatGPT) (*Agent, error) {
	config, err := config.LoadDynamicConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	agent := &Agent{
		chatgpt: chatgptClient,
		config:  config,
		mode:    InteractiveMode,
		fileOps: NewFileOperations(),
	}

	// Initialize project context if enabled
	if config.Agent.ProjectAnalysis {
		agent.context = NewProjectContext()
	}

	return agent, nil
}

// SetMode changes the agent's operation mode
func (a *Agent) SetMode(mode AgentMode) {
	a.mode = mode
	ui.PrintInfo(fmt.Sprintf("Agent mode set to: %s", mode))
}

// GetMode returns the current agent mode
func (a *Agent) GetMode() AgentMode {
	return a.mode
}

// ProcessMessage processes a message based on the current mode
func (a *Agent) ProcessMessage(message string) (string, error) {
	switch a.mode {
	case InteractiveMode:
		return a.processInteractive(message)
	case QueryMode:
		return a.processQuery(message)
	case AutoMode:
		return a.processAuto(message)
	case ContextMode:
		return a.processWithContext(message)
	default:
		return a.processInteractive(message)
	}
}

// processInteractive handles interactive mode (default behavior)
func (a *Agent) processInteractive(message string) (string, error) {
	return a.chatgpt.SendMessage(message)
}

// processQuery handles single query mode
func (a *Agent) processQuery(message string) (string, error) {
	// For query mode, we might want to add specific formatting
	response, err := a.chatgpt.SendMessage(message)
	if err != nil {
		return "", err
	}
	
	// In query mode, we could format the response differently
	return response, nil
}

// processAuto handles autonomous mode
func (a *Agent) processAuto(message string) (string, error) {
	// Auto mode could include task breakdown, planning, etc.
	// For now, delegate to interactive mode
	return a.processInteractive(message)
}

// processWithContext handles context-aware processing
func (a *Agent) processWithContext(message string) (string, error) {
	if a.context != nil {
		// Enhance message with project context
		contextualMessage := a.context.EnhanceMessage(message)
		return a.chatgpt.SendMessage(contextualMessage)
	}
	return a.processInteractive(message)
}

// InitializeSession sets up the agent session with project context
func (a *Agent) InitializeSession() error {
	if !a.config.Agent.AutoContext {
		return nil
	}

	prompts, err := config.GetPrompts()
	if err != nil {
		return fmt.Errorf("failed to load prompts: %v", err)
	}

	// Generate system prompt based on project context
	systemPrompt := a.generateSystemPrompt(prompts)
	
	spinner := ui.NewSquareSpinner()
	spinner.Start("Analyzing project and setting up context...")
	
	// Send system prompt
	_, err = a.chatgpt.SendMessage(systemPrompt)
	spinner.Stop()
	
	if err != nil {
		ui.PrintWarning("Could not set up project context")
		return err
	}
	
	ui.PrintSuccess("Project context established! 🎯")
	return nil
}

// generateSystemPrompt creates a system prompt based on configuration
func (a *Agent) generateSystemPrompt(prompts *config.Prompts) string {
	var systemPrompt strings.Builder
	
	// Add role and personality
	defaultAgent := prompts.SystemPrompts.DefaultAgent
	systemPrompt.WriteString(defaultAgent.Role + "\n\n")
	systemPrompt.WriteString(defaultAgent.Personality + "\n\n")
	
	// Add capabilities
	systemPrompt.WriteString("Your capabilities:\n")
	for _, capability := range defaultAgent.Capabilities {
		systemPrompt.WriteString("- " + capability + "\n")
	}
	systemPrompt.WriteString("\n")
	
	// Add project context if available
	if a.context != nil {
		projectInfo := a.context.GetProjectInfo()
		contextTemplate := prompts.SystemPrompts.ProjectContext.Template
		
		// Replace placeholders
		contextPrompt := strings.ReplaceAll(contextTemplate, "{current_dir}", a.context.GetCurrentDir())
		contextPrompt = strings.ReplaceAll(contextPrompt, "{project_info}", projectInfo)
		contextPrompt = strings.ReplaceAll(contextPrompt, "{role_description}", defaultAgent.Role)
		
		systemPrompt.WriteString(contextPrompt)
	}
	
	return systemPrompt.String()
}

// StartNewChat starts a new chat session
func (a *Agent) StartNewChat() error {
	err := a.chatgpt.StartNewChat()
	if err != nil {
		return err
	}
	
	// Re-initialize session with context
	return a.InitializeSession()
}

// GetConfig returns the agent's configuration
func (a *Agent) GetConfig() *config.DynamicConfig {
	return a.config
}

// UpdateConfig updates the agent's configuration
func (a *Agent) UpdateConfig(newConfig *config.DynamicConfig) error {
	a.config = newConfig
	return newConfig.SaveConfig()
}

// GetProjectContext returns the current project context
func (a *Agent) GetProjectContext() *ProjectContext {
	return a.context
}

// RefreshProjectContext refreshes the project analysis
func (a *Agent) RefreshProjectContext() error {
	if a.context != nil {
		return a.context.Refresh()
	}
	return nil
}

// File Access Methods

// ReadFile reads a specific file and returns its content
func (a *Agent) ReadFile(filename string) (string, error) {
	return a.fileOps.ReadFile(filename)
}

// ListFiles lists all files in the current directory or specified path
func (a *Agent) ListFiles(path string) ([]FileInfo, error) {
	return a.fileOps.ListFiles(path)
}

// SearchFiles searches for files matching a pattern
func (a *Agent) SearchFiles(pattern string) ([]FileInfo, error) {
	return a.fileOps.SearchFiles(pattern)
}

// ReadMultipleFiles reads multiple files and returns their content
func (a *Agent) ReadMultipleFiles(filenames []string) (map[string]string, error) {
	return a.fileOps.ReadMultipleFiles(filenames)
}

// GetFileTree returns a tree structure of the project
func (a *Agent) GetFileTree(maxDepth int) (string, error) {
	return a.fileOps.GetFileTree(maxDepth)
}

// ProcessFileQuery processes queries related to file operations
func (a *Agent) ProcessFileQuery(query string) (string, error) {
	// Detect file-related queries and provide appropriate responses
	lowerQuery := strings.ToLower(query)
	
	// Check for file reading requests
	if strings.Contains(lowerQuery, "read file") || strings.Contains(lowerQuery, "show me") {
		return a.handleFileReadRequest(query)
	}
	
	// Check for file listing requests
	if strings.Contains(lowerQuery, "list files") || strings.Contains(lowerQuery, "show files") {
		return a.handleFileListRequest(query)
	}
	
	// Check for file tree requests
	if strings.Contains(lowerQuery, "file tree") || strings.Contains(lowerQuery, "project structure") {
		return a.handleFileTreeRequest(query)
	}
	
	// Check for file search requests
	if strings.Contains(lowerQuery, "find file") || strings.Contains(lowerQuery, "search") {
		return a.handleFileSearchRequest(query)
	}
	
	// Default: process as normal message
	return a.ProcessMessage(query)
}

// handleFileReadRequest handles requests to read specific files
func (a *Agent) handleFileReadRequest(query string) (string, error) {
	// Extract filename from query (simple implementation)
	words := strings.Fields(query)
	var filename string
	
	for i, word := range words {
		// Look for file extensions or common filenames
		if strings.Contains(word, ".") || word == "main.go" || word == "README.md" {
			filename = word
			break
		}
		// Look for patterns like "read file main.go"
		if (word == "file" || word == "File") && i+1 < len(words) {
			filename = words[i+1]
			break
		}
	}
	
	if filename == "" {
		return "Please specify which file you'd like me to read. For example: 'read file main.go'", nil
	}
	
	content, err := a.ReadFile(filename)
	if err != nil {
		return fmt.Sprintf("Sorry, I couldn't read the file '%s': %v", filename, err), nil
	}
	
	// Send file content to ChatGPT with context
	contextualQuery := fmt.Sprintf("Here's the content of %s:\n\n```\n%s\n```\n\nPlease analyze this file and provide insights about the code structure, functionality, and any suggestions for improvement.", filename, content)
	
	return a.chatgpt.SendMessage(contextualQuery)
}

// handleFileListRequest handles requests to list files
func (a *Agent) handleFileListRequest(query string) (string, error) {
	files, err := a.ListFiles("")
	if err != nil {
		return fmt.Sprintf("Sorry, I couldn't list the files: %v", err), nil
	}
	
	var response strings.Builder
	response.WriteString("Here are the files in your project:\n\n")
	
	// Group files by category
	categories := make(map[FileCategory][]FileInfo)
	for _, file := range files {
		categories[file.Category] = append(categories[file.Category], file)
	}
	
	// Display by category
	categoryNames := map[FileCategory]string{
		CodeFile:     "📄 Code Files",
		ConfigFile:   "⚙️ Configuration Files",
		DocumentFile: "📚 Documentation",
		TestFile:     "🧪 Test Files",
		BuildFile:    "🔨 Build Files",
		UnknownFile:  "📁 Other Files",
	}
	
	for category, categoryFiles := range categories {
		if len(categoryFiles) > 0 {
			response.WriteString(fmt.Sprintf("\n%s:\n", categoryNames[category]))
			for _, file := range categoryFiles {
				response.WriteString(fmt.Sprintf("  - %s\n", file.Path))
			}
		}
	}
	
	// Send to ChatGPT for analysis
	contextualQuery := fmt.Sprintf("%s\n\nPlease analyze this project structure and provide insights about the codebase organization.", response.String())
	
	return a.chatgpt.SendMessage(contextualQuery)
}

// handleFileTreeRequest handles requests for file tree
func (a *Agent) handleFileTreeRequest(query string) (string, error) {
	tree, err := a.GetFileTree(3) // Max depth of 3
	if err != nil {
		return fmt.Sprintf("Sorry, I couldn't generate the file tree: %v", err), nil
	}
	
	contextualQuery := fmt.Sprintf("Here's the project file tree structure:\n\n```\n%s\n```\n\nPlease analyze this project structure and provide insights about the organization and architecture.", tree)
	
	return a.chatgpt.SendMessage(contextualQuery)
}

// handleFileSearchRequest handles file search requests
func (a *Agent) handleFileSearchRequest(query string) (string, error) {
	// Extract search pattern from query
	words := strings.Fields(query)
	var pattern string
	
	for i, word := range words {
		if (word == "find" || word == "search") && i+1 < len(words) {
			pattern = words[i+1]
			break
		}
	}
	
	if pattern == "" {
		return "Please specify what file you're looking for. For example: 'find file main' or 'search config'", nil
	}
	
	files, err := a.SearchFiles(pattern)
	if err != nil {
		return fmt.Sprintf("Sorry, I couldn't search for files: %v", err), nil
	}
	
	if len(files) == 0 {
		return fmt.Sprintf("No files found matching '%s'", pattern), nil
	}
	
	var response strings.Builder
	response.WriteString(fmt.Sprintf("Found %d file(s) matching '%s':\n\n", len(files), pattern))
	
	for _, file := range files {
		response.WriteString(fmt.Sprintf("  - %s (%s)\n", file.Path, file.Category))
	}
	
	return response.String(), nil
}