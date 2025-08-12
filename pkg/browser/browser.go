package browser

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chatgpt-element-recorder/pkg/config"
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
		cookiesData, err := os.ReadFile(config.CookiesFile)
		if os.IsNotExist(err) {
			// No cookies file - continue silently
			return nil
		} else if err != nil {
			return err
		}

		var chatgptCookies []ChatGPTCookie
		if err := json.Unmarshal(cookiesData, &chatgptCookies); err != nil {
			return err
		}

		var cookies []*network.CookieParam
		for _, cookie := range chatgptCookies {
			cookieParam := &network.CookieParam{
				Name:     cookie.Name,
				Value:    cookie.Value,
				Domain:   cookie.Domain,
				Path:     cookie.Path,
				Secure:   cookie.Secure,
				HTTPOnly: cookie.HTTPOnly,
			}

			// Skip expiry handling to avoid compatibility issues
			// Cookies will still work without explicit expiry dates

			switch strings.ToLower(cookie.SameSite) {
			case "strict":
				cookieParam.SameSite = network.CookieSameSiteStrict
			case "lax":
				cookieParam.SameSite = network.CookieSameSiteLax
			case "none", "no_restriction":
				cookieParam.SameSite = network.CookieSameSiteNone
			}

			cookies = append(cookies, cookieParam)
		}

		// Cookies loaded silently for clean UI
		return network.SetCookies(cookies).Do(ctx)
	})
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
