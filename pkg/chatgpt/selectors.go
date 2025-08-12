package chatgpt

// Selectors are hardcoded for stability and simplicity.
const (
	InputElement     = `#prompt-textarea`
	SubmitButton     = `button[data-testid="send-button"]`
	StopButton       = `button[data-testid="stop-button"]`
	LastResponse     = `div[data-message-author-role="assistant"]:last-child .markdown`
	NewChatButton    = `a[href="/"]`
	HistoryLink      = `a[href^="/c/"]`
	AssistantMessage = `div[data-message-author-role="assistant"]`
)
