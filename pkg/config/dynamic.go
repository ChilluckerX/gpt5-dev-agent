package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// DynamicConfig represents the complete application configuration
type DynamicConfig struct {
	ChatGPT ChatGPTConfig `json:"chatgpt"`
	Browser BrowserConfig `json:"browser"`
	Files   FilesConfig   `json:"files"`
	UI      UIConfig      `json:"ui"`
	Agent   AgentConfig   `json:"agent"`
	mu      sync.RWMutex  `json:"-"`
}

// ChatGPTConfig contains ChatGPT-specific settings
type ChatGPTConfig struct {
	BaseURL       string `json:"base_url"`
	Timeout       int    `json:"timeout"`
	RetryAttempts int    `json:"retry_attempts"`
	WaitTimeout   int    `json:"wait_timeout"`
}

// BrowserConfig contains browser automation settings
type BrowserConfig struct {
	Headless           bool   `json:"headless"`
	WindowSize         string `json:"window_size"`
	UserAgent          string `json:"user_agent"`
	DisableAutomation  bool   `json:"disable_automation"`
	DisableExtensions  bool   `json:"disable_extensions"`
}

// FilesConfig contains file path settings
type FilesConfig struct {
	CookiesFile string `json:"cookies_file"`
	OutputDir   string `json:"output_dir"`
	ConfigDir   string `json:"config_dir"`
}

// UIConfig contains UI appearance settings
type UIConfig struct {
	SpinnerType  string            `json:"spinner_type"`
	TypingSpeed  int               `json:"typing_speed"`
	BorderSpeed  int               `json:"border_speed"`
	Colors       map[string]string `json:"colors"`
}

// AgentConfig contains agent behavior settings
type AgentConfig struct {
	Mode               string `json:"mode"`
	AutoContext        bool   `json:"auto_context"`
	ProjectAnalysis    bool   `json:"project_analysis"`
	SessionPersistence bool   `json:"session_persistence"`
}

// Selectors represents CSS selectors configuration
type Selectors struct {
	Input          SelectorGroup `json:"input"`
	SendButton     SelectorGroup `json:"send_button"`
	Response       SelectorGroup `json:"response"`
	ChatControls   SelectorMap   `json:"chat_controls"`
	PageElements   SelectorMap   `json:"page_elements"`
	Authentication SelectorMap   `json:"authentication"`
}

// SelectorGroup represents a primary selector with fallbacks
type SelectorGroup struct {
	Primary  string   `json:"primary"`
	Fallback []string `json:"fallback"`
}

// SelectorMap represents a map of named selectors
type SelectorMap map[string]string

// Prompts represents system prompts configuration
type Prompts struct {
	SystemPrompts    SystemPrompts              `json:"system_prompts"`
	ResponseFormats  map[string]interface{}     `json:"response_formats"`
	ProjectTemplates map[string]ProjectTemplate `json:"project_templates"`
}

// SystemPrompts contains various system prompt configurations
type SystemPrompts struct {
	DefaultAgent     AgentPrompt            `json:"default_agent"`
	ProjectContext   ProjectContextPrompt   `json:"project_context"`
	SpecializedModes map[string]string      `json:"specialized_modes"`
}

// AgentPrompt defines the agent's role and personality
type AgentPrompt struct {
	Role         string   `json:"role"`
	Personality  string   `json:"personality"`
	Capabilities []string `json:"capabilities"`
}

// ProjectContextPrompt defines how project context is presented
type ProjectContextPrompt struct {
	Template      string `json:"template"`
	GreetingStyle string `json:"greeting_style"`
}

// ProjectTemplate defines project-specific prompts
type ProjectTemplate struct {
	Greeting   string   `json:"greeting"`
	FocusAreas []string `json:"focus_areas"`
}

var (
	globalConfig    *DynamicConfig
	globalSelectors *Selectors
	globalPrompts   *Prompts
	configOnce      sync.Once
)

// LoadDynamicConfig loads configuration from JSON files
func LoadDynamicConfig() (*DynamicConfig, error) {
	var err error
	configOnce.Do(func() {
		globalConfig, err = loadConfigFromFile()
	})
	return globalConfig, err
}

// GetSelectors loads and returns CSS selectors
func GetSelectors() (*Selectors, error) {
	if globalSelectors == nil {
		selectors, err := loadSelectorsFromFile()
		if err != nil {
			return nil, err
		}
		globalSelectors = selectors
	}
	return globalSelectors, nil
}

// GetPrompts loads and returns system prompts
func GetPrompts() (*Prompts, error) {
	if globalPrompts == nil {
		prompts, err := loadPromptsFromFile()
		if err != nil {
			return nil, err
		}
		globalPrompts = prompts
	}
	return globalPrompts, nil
}

// loadConfigFromFile loads main configuration
func loadConfigFromFile() (*DynamicConfig, error) {
	configPath := "configs/config.json"
	data, err := os.ReadFile(configPath)
	if err != nil {
		return getDefaultConfig(), fmt.Errorf("failed to read config file: %v", err)
	}

	var config DynamicConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return getDefaultConfig(), fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

// loadSelectorsFromFile loads CSS selectors
func loadSelectorsFromFile() (*Selectors, error) {
	selectorsPath := "configs/selectors.json"
	data, err := os.ReadFile(selectorsPath)
	if err != nil {
		return getDefaultSelectors(), fmt.Errorf("failed to read selectors file: %v", err)
	}

	var selectors Selectors
	if err := json.Unmarshal(data, &selectors); err != nil {
		return getDefaultSelectors(), fmt.Errorf("failed to parse selectors file: %v", err)
	}

	return &selectors, nil
}

// loadPromptsFromFile loads system prompts
func loadPromptsFromFile() (*Prompts, error) {
	promptsPath := "configs/prompts.json"
	data, err := os.ReadFile(promptsPath)
	if err != nil {
		return getDefaultPrompts(), fmt.Errorf("failed to read prompts file: %v", err)
	}

	var prompts Prompts
	if err := json.Unmarshal(data, &prompts); err != nil {
		return getDefaultPrompts(), fmt.Errorf("failed to parse prompts file: %v", err)
	}

	return &prompts, nil
}

// SaveConfig saves the current configuration to file
func (c *DynamicConfig) SaveConfig() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	configPath := "configs/config.json"
	
	// Ensure config directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// GetString safely gets a string value with fallback
func (c *DynamicConfig) GetString(key, fallback string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	// This would be expanded to handle nested key access
	// For now, return fallback
	return fallback
}

// SetValue safely sets a configuration value
func (c *DynamicConfig) SetValue(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	// This would be expanded to handle nested key setting
	// For now, just save the config
	return c.SaveConfig()
}

// GetCookiesPath returns the full path to cookies file
func (c *DynamicConfig) GetCookiesPath() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Files.CookiesFile
}

// GetBaseURL returns the ChatGPT base URL
func (c *DynamicConfig) GetBaseURL() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.ChatGPT.BaseURL
}