package browser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chatgpt-element-recorder/pkg/config"
	"github.com/chatgpt-element-recorder/pkg/ui"
)

// CookieInfo represents a browser cookie
type CookieInfo struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	Domain   string  `json:"domain"`
	Path     string  `json:"path"`
	Expires  float64 `json:"expires,omitempty"`
	HTTPOnly bool    `json:"httpOnly,omitempty"`
	Secure   bool    `json:"secure,omitempty"`
	SameSite string  `json:"sameSite,omitempty"`
}

// CookieManager handles cookie operations
type CookieManager struct {
	cookiesPath string
}

// NewCookieManager creates a new cookie manager
func NewCookieManager() *CookieManager {
	cfg, err := config.LoadDynamicConfig()
	cookiesPath := "cookies/chatgpt.json" // default
	if err == nil {
		cookiesPath = cfg.Files.CookiesFile
	}
	
	return &CookieManager{
		cookiesPath: cookiesPath,
	}
}

// EnsureCookiesFile ensures the cookies file exists and is valid
func (cm *CookieManager) EnsureCookiesFile() error {
	// Ensure cookies directory exists
	cookiesDir := filepath.Dir(cm.cookiesPath)
	if err := os.MkdirAll(cookiesDir, 0755); err != nil {
		return fmt.Errorf("failed to create cookies directory: %v", err)
	}

	// Check if cookies file exists
	if _, err := os.Stat(cm.cookiesPath); os.IsNotExist(err) {
		ui.PrintInfo("Creating new cookies file...")
		return cm.createEmptyCookiesFile()
	}

	// Validate existing cookies file
	return cm.validateCookiesFile()
}

// createEmptyCookiesFile creates an empty but valid cookies file
func (cm *CookieManager) createEmptyCookiesFile() error {
	emptyCookies := []CookieInfo{}
	
	data, err := json.MarshalIndent(emptyCookies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal empty cookies: %v", err)
	}

	if err := os.WriteFile(cm.cookiesPath, data, 0644); err != nil {
		return fmt.Errorf("failed to create cookies file: %v", err)
	}

	ui.PrintSuccess(fmt.Sprintf("Created empty cookies file: %s", cm.cookiesPath))
	ui.PrintInfo("You may need to login manually to ChatGPT first")
	return nil
}

// validateCookiesFile validates the existing cookies file
func (cm *CookieManager) validateCookiesFile() error {
	data, err := os.ReadFile(cm.cookiesPath)
	if err != nil {
		return fmt.Errorf("failed to read cookies file: %v", err)
	}

	// Check if file is empty
	if len(data) == 0 {
		ui.PrintWarning("Cookies file is empty, creating default structure...")
		return cm.createEmptyCookiesFile()
	}

	// Try to parse JSON
	var cookies []CookieInfo
	if err := json.Unmarshal(data, &cookies); err != nil {
		ui.PrintWarning("Invalid cookies file format, backing up and recreating...")
		return cm.backupAndRecreate()
	}

	// Validate cookie content
	return cm.validateCookieContent(cookies)
}

// backupAndRecreate backs up invalid cookies file and creates new one
func (cm *CookieManager) backupAndRecreate() error {
	// Create backup filename with timestamp
	timestamp := time.Now().Format("20060102-150405")
	backupPath := fmt.Sprintf("%s.backup-%s", cm.cookiesPath, timestamp)
	
	// Copy current file to backup
	data, err := os.ReadFile(cm.cookiesPath)
	if err == nil {
		os.WriteFile(backupPath, data, 0644)
		ui.PrintInfo(fmt.Sprintf("Backed up invalid cookies to: %s", backupPath))
	}

	// Create new empty cookies file
	return cm.createEmptyCookiesFile()
}

// validateCookieContent validates the content of cookies
func (cm *CookieManager) validateCookieContent(cookies []CookieInfo) error {
	if len(cookies) == 0 {
		ui.PrintWarning("No cookies found - you may need to login to ChatGPT")
		return nil
	}

	// Check for essential ChatGPT cookies
	hasSessionCookie := false
	hasAuthCookie := false
	validCookieCount := 0
	expiredCookieCount := 0
	
	currentTime := float64(time.Now().Unix())
	
	for _, cookie := range cookies {
		// Check if cookie is for ChatGPT domain
		if !cm.isChatGPTCookie(cookie) {
			continue
		}
		
		validCookieCount++
		
		// Check for session-related cookies
		if cm.isSessionCookie(cookie) {
			hasSessionCookie = true
		}
		
		// Check for auth-related cookies
		if cm.isAuthCookie(cookie) {
			hasAuthCookie = true
		}
		
		// Check if cookie is expired
		if cookie.Expires > 0 && cookie.Expires < currentTime {
			expiredCookieCount++
		}
	}

	// Report validation results
	ui.PrintInfo(fmt.Sprintf("Found %d ChatGPT cookies", validCookieCount))
	
	if expiredCookieCount > 0 {
		ui.PrintWarning(fmt.Sprintf("%d cookies have expired", expiredCookieCount))
	}
	
	if !hasSessionCookie && !hasAuthCookie {
		ui.PrintWarning("No authentication cookies found - manual login may be required")
	} else {
		ui.PrintSuccess("Authentication cookies detected")
	}

	return nil
}

// isChatGPTCookie checks if cookie belongs to ChatGPT
func (cm *CookieManager) isChatGPTCookie(cookie CookieInfo) bool {
	chatgptDomains := []string{
		"chatgpt.com",
		".chatgpt.com", 
		"chat.openai.com",
		".openai.com",
	}
	
	for _, domain := range chatgptDomains {
		if cookie.Domain == domain {
			return true
		}
	}
	return false
}

// isSessionCookie checks if cookie is session-related
func (cm *CookieManager) isSessionCookie(cookie CookieInfo) bool {
	sessionCookieNames := []string{
		"__Secure-next-auth.session-token",
		"next-auth.session-token", 
		"session",
		"sessionid",
		"JSESSIONID",
	}
	
	for _, name := range sessionCookieNames {
		if cookie.Name == name {
			return true
		}
	}
	return false
}

// isAuthCookie checks if cookie is auth-related
func (cm *CookieManager) isAuthCookie(cookie CookieInfo) bool {
	authCookieNames := []string{
		"__Secure-next-auth.csrf-token",
		"next-auth.csrf-token",
		"csrf-token",
		"auth-token",
		"access_token",
		"_auth",
	}
	
	for _, name := range authCookieNames {
		if cookie.Name == name {
			return true
		}
	}
	return false
}

// GetCookiesPath returns the cookies file path
func (cm *CookieManager) GetCookiesPath() string {
	return cm.cookiesPath
}

// LoadCookies loads and validates cookies
func (cm *CookieManager) LoadCookies() ([]CookieInfo, error) {
	// Ensure cookies file exists and is valid
	if err := cm.EnsureCookiesFile(); err != nil {
		return nil, err
	}

	// Read cookies file
	data, err := os.ReadFile(cm.cookiesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cookies file: %v", err)
	}

	// Parse cookies
	var cookies []CookieInfo
	if err := json.Unmarshal(data, &cookies); err != nil {
		return nil, fmt.Errorf("failed to parse cookies: %v", err)
	}

	return cookies, nil
}

// SaveCookies saves cookies to file
func (cm *CookieManager) SaveCookies(cookies []CookieInfo) error {
	// Ensure cookies directory exists
	cookiesDir := filepath.Dir(cm.cookiesPath)
	if err := os.MkdirAll(cookiesDir, 0755); err != nil {
		return fmt.Errorf("failed to create cookies directory: %v", err)
	}

	// Marshal cookies to JSON
	data, err := json.MarshalIndent(cookies, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cookies: %v", err)
	}

	// Write to file
	if err := os.WriteFile(cm.cookiesPath, data, 0644); err != nil {
		return fmt.Errorf("failed to save cookies: %v", err)
	}

	ui.PrintSuccess(fmt.Sprintf("Saved %d cookies to %s", len(cookies), cm.cookiesPath))
	return nil
}

// CleanExpiredCookies removes expired cookies
func (cm *CookieManager) CleanExpiredCookies() error {
	cookies, err := cm.LoadCookies()
	if err != nil {
		return err
	}

	currentTime := float64(time.Now().Unix())
	var validCookies []CookieInfo
	removedCount := 0

	for _, cookie := range cookies {
		// Keep cookies that are not expired or have no expiry
		if cookie.Expires == 0 || cookie.Expires > currentTime {
			validCookies = append(validCookies, cookie)
		} else {
			removedCount++
		}
	}

	if removedCount > 0 {
		ui.PrintInfo(fmt.Sprintf("Removed %d expired cookies", removedCount))
		return cm.SaveCookies(validCookies)
	}

	return nil
}