package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chatgpt-element-recorder/pkg/config"
	"github.com/chatgpt-element-recorder/pkg/ui"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// ChatGPTCookie represents a cookie from the JSON file
type ChatGPTCookie struct {
	Domain         string  `json:"domain"`
	ExpirationDate float64 `json:"expirationDate,omitempty"`
	HostOnly       bool    `json:"hostOnly"`
	HTTPOnly       bool    `json:"httpOnly"`
	Name           string  `json:"name"`
	Path           string  `json:"path"`
	SameSite       string  `json:"sameSite,omitempty"`
	Secure         bool    `json:"secure"`
	Session        bool    `json:"session"`
	StoreID        *string `json:"storeId"`
	Value          string  `json:"value"`
}

func LoadCookiesAction() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// Create cookie manager
		cookieManager := NewCookieManager()
		
		// Ensure cookies file exists and is valid
		if err := cookieManager.EnsureCookiesFile(); err != nil {
			ui.PrintWarning(fmt.Sprintf("Cookie validation failed: %v", err))
			return nil // Continue without cookies
		}
		
		// Load validated cookies using legacy format for compatibility
		cookiesData, err := os.ReadFile(cookieManager.GetCookiesPath())
		if os.IsNotExist(err) {
			ui.PrintInfo("No cookies file found - manual login required")
			return nil
		} else if err != nil {
			ui.PrintWarning(fmt.Sprintf("Failed to read cookies: %v", err))
			return nil
		}

		// Try to parse as legacy ChatGPTCookie format first
		var chatgptCookies []ChatGPTCookie
		if err := json.Unmarshal(cookiesData, &chatgptCookies); err != nil {
			ui.PrintWarning("Invalid cookie format - manual login required")
			return nil
		}

		if len(chatgptCookies) == 0 {
			ui.PrintInfo("No cookies to load - manual login required")
			return nil
		}

		// Convert to network.CookieParam and validate
		var cookies []*network.CookieParam
		validCookieCount := 0
		expiredCookieCount := 0
		currentTime := float64(time.Now().Unix())

		for _, cookie := range chatgptCookies {
			// Check if cookie is expired
			if cookie.ExpirationDate > 0 && cookie.ExpirationDate < currentTime {
				expiredCookieCount++
				continue // Skip expired cookies
			}

			// Only load ChatGPT-related cookies
			if !isChatGPTDomain(cookie.Domain) {
				continue
			}

			cookieParam := &network.CookieParam{
				Name:     cookie.Name,
				Value:    cookie.Value,
				Domain:   cookie.Domain,
				Path:     cookie.Path,
				Secure:   cookie.Secure,
				HTTPOnly: cookie.HTTPOnly,
			}

			// Set expiry if available
			if cookie.ExpirationDate > 0 {
				expires := cdp.TimeSinceEpoch(time.Unix(int64(cookie.ExpirationDate), 0))
				cookieParam.Expires = &expires
			}

			// Set SameSite attribute
			switch strings.ToLower(cookie.SameSite) {
			case "strict":
				cookieParam.SameSite = network.CookieSameSiteStrict
			case "lax":
				cookieParam.SameSite = network.CookieSameSiteLax
			case "none", "no_restriction":
				cookieParam.SameSite = network.CookieSameSiteNone
			}

			cookies = append(cookies, cookieParam)
			validCookieCount++
		}

		// Report cookie loading status
		if expiredCookieCount > 0 {
			ui.PrintWarning(fmt.Sprintf("Skipped %d expired cookies", expiredCookieCount))
		}

		if validCookieCount == 0 {
			ui.PrintWarning("No valid ChatGPT cookies found - manual login required")
			return nil
		}

		// Load cookies into browser
		if err := network.SetCookies(cookies).Do(ctx); err != nil {
			ui.PrintWarning(fmt.Sprintf("Failed to load cookies: %v", err))
			return nil
		}

		ui.PrintSuccess(fmt.Sprintf("Loaded %d ChatGPT cookies", validCookieCount))
		return nil
	})
}

// isChatGPTDomain checks if domain belongs to ChatGPT
func isChatGPTDomain(domain string) bool {
	chatgptDomains := []string{
		"chatgpt.com",
		".chatgpt.com", 
		"chat.openai.com",
		".openai.com",
	}
	
	for _, chatgptDomain := range chatgptDomains {
		if domain == chatgptDomain {
			return true
		}
	}
	return false
}

// SaveCookiesAction retrieves cookies from the browser and saves them to a file.
func SaveCookiesAction() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		log.Println("Saving cookies to", config.CookiesFile)
		cookies, err := network.GetCookies().Do(ctx)
		if err != nil {
			return err
		}

		cookiesData, err := json.MarshalIndent(cookies, "", "  ")
		if err != nil {
			return err
		}

		return os.WriteFile(config.CookiesFile, cookiesData, 0644)
	})
}

// WaitForUserInteraction waits for user to perform an action and provides instructions
func WaitForUserInteraction(instruction string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		log.Println("---")
		log.Printf("ACTION REQUIRED: %s", instruction)
		log.Println("Please perform the action in the browser window.")
		log.Println("Press ENTER in this terminal when you're done...")
		log.Println("---")

		// Wait for user input
		var input string
		_, err := fmt.Scanln(&input)
		return err
	})
}

// WaitForChatGPTLoad waits for ChatGPT to fully load
func WaitForChatGPTLoad() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// Just wait for basic page elements - no noisy logging
		err := chromedp.WaitVisible(`main`, chromedp.ByQuery).Do(ctx)
		if err != nil {
			// Fallback - wait for body
			chromedp.WaitVisible(`body`, chromedp.ByQuery).Do(ctx)
		}
		return nil
	})
}
