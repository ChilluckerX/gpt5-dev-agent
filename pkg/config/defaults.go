package config

// getDefaultConfig returns default configuration when config file is not available
func getDefaultConfig() *DynamicConfig {
	return &DynamicConfig{
		ChatGPT: ChatGPTConfig{
			BaseURL:       "https://chatgpt.com",
			Timeout:       300,
			RetryAttempts: 3,
			WaitTimeout:   30,
		},
		Browser: BrowserConfig{
			Headless:          false,
			WindowSize:        "1920,1080",
			UserAgent:         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			DisableAutomation: true,
			DisableExtensions: false,
		},
		Files: FilesConfig{
			CookiesFile: "cookies/chatgpt.json",
			OutputDir:   "output",
			ConfigDir:   "configs",
		},
		UI: UIConfig{
			SpinnerType: "square",
			TypingSpeed: 30,
			BorderSpeed: 10,
			Colors: map[string]string{
				"success": "\033[32m",
				"error":   "\033[31m",
				"warning": "\033[33m",
				"info":    "\033[36m",
				"dim":     "\033[2m",
				"reset":   "\033[0m",
			},
		},
		Agent: AgentConfig{
			Mode:               "interactive",
			AutoContext:        true,
			ProjectAnalysis:    true,
			SessionPersistence: true,
		},
	}
}

// getDefaultSelectors returns default CSS selectors when selectors file is not available
func getDefaultSelectors() *Selectors {
	return &Selectors{
		Input: SelectorGroup{
			Primary: "#prompt-textarea",
			Fallback: []string{
				"textarea[placeholder*='Message']",
				"textarea[data-id='root']",
				"[contenteditable='true'][data-testid*='textbox']",
			},
		},
		SendButton: SelectorGroup{
			Primary: "[data-testid='send-button']",
			Fallback: []string{
				"button[aria-label*='Send']",
				"button:has(svg[data-icon='send'])",
				"[data-testid='fruitjuice-send-button']",
			},
		},
		Response: SelectorGroup{
			Primary: "[data-message-author-role='assistant']",
			Fallback: []string{
				".group\\/conversation-turn .markdown",
				"[data-testid*='conversation-turn-'] .markdown",
			},
		},
		ChatControls: SelectorMap{
			"new_chat":       "a[href='/']",
			"stop_generating": "[aria-label*='Stop']",
			"regenerate":     "[aria-label*='Regenerate']",
		},
		PageElements: SelectorMap{
			"chat_list":         "[data-testid='conversation-turn-']",
			"sidebar":           "[data-testid='sidebar']",
			"main_content":      "main",
			"loading_indicator": "[data-testid*='loading']",
		},
		Authentication: SelectorMap{
			"login_button":  "[data-testid='login-button']",
			"signup_button": "[data-testid='signup-button']",
			"user_menu":     "[data-testid='user-menu']",
		},
	}
}

// getDefaultPrompts returns default system prompts when prompts file is not available
func getDefaultPrompts() *Prompts {
	return &Prompts{
		SystemPrompts: SystemPrompts{
			DefaultAgent: AgentPrompt{
				Role:        "You are RovoDev, a friendly and expert software development assistant.",
				Personality: "Be conversational and friendly, like a senior developer colleague. Ask intelligent follow-up questions about their work.",
				Capabilities: []string{
					"Act as a knowledgeable coding assistant and mentor",
					"Provide helpful suggestions based on the project structure you see",
					"Offer specific help based on the technologies and files you observe",
					"Analyze code patterns and suggest improvements",
				},
			},
			ProjectContext: ProjectContextPrompt{
				Template:      "You are helping a developer who is currently working in the directory: {current_dir}\n\nProject Analysis:\n{project_info}\n\n{role_description}\n\nPlease greet the user by acknowledging what you see in their project and ask how you can help them today. Be specific about what you notice in their codebase.",
				GreetingStyle: "professional_friendly",
			},
			SpecializedModes: map[string]string{
				"code_review":  "Focus on code quality, best practices, and potential improvements.",
				"debugging":    "Help identify and solve bugs, errors, and issues in the code.",
				"architecture": "Provide guidance on system design, architecture patterns, and scalability.",
				"learning":     "Explain concepts, provide tutorials, and help with learning new technologies.",
			},
		},
		ResponseFormats: map[string]interface{}{
			"code_block": map[string]string{
				"prefix": "```{language}\n",
				"suffix": "\n```",
			},
			"explanation": map[string][]string{
				"structure": {"overview", "details", "examples", "next_steps"},
			},
		},
		ProjectTemplates: map[string]ProjectTemplate{
			"go": {
				Greeting:   "I can see you're working on a Go project. I notice {project_details}. How can I help you with your Go development today?",
				FocusAreas: []string{"performance", "concurrency", "modules", "testing"},
			},
			"python": {
				Greeting:   "I see you're working on a Python project. I notice {project_details}. What would you like to work on?",
				FocusAreas: []string{"packages", "virtual_env", "testing", "performance"},
			},
			"javascript": {
				Greeting:   "I can see you're working on a JavaScript/Node.js project. I notice {project_details}. How can I assist you?",
				FocusAreas: []string{"npm", "dependencies", "async", "testing"},
			},
			"generic": {
				Greeting:   "I can see you're working on a {project_type} project. I notice {project_details}. How can I help you today?",
				FocusAreas: []string{"structure", "dependencies", "best_practices"},
			},
		},
	}
}

// Legacy compatibility functions to maintain existing API

// GetLegacyBaseURL returns base URL for legacy compatibility
func GetLegacyBaseURL() string {
	config, err := LoadDynamicConfig()
	if err != nil {
		return "https://chatgpt.com" // fallback to hardcoded value
	}
	return config.GetBaseURL()
}

// GetLegacyCookiesFile returns cookies file path for legacy compatibility
func GetLegacyCookiesFile() string {
	config, err := LoadDynamicConfig()
	if err != nil {
		return "cookies/chatgpt.json" // fallback to hardcoded value
	}
	return config.GetCookiesPath()
}