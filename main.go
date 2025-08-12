package main

import (
	"context"
	"log"
	"time"

	"github.com/chatgpt-element-recorder/pkg/browser"
	"github.com/chatgpt-element-recorder/pkg/chatgpt"
	"github.com/chatgpt-element-recorder/pkg/cli"
	"github.com/chatgpt-element-recorder/pkg/config"
	"github.com/chatgpt-element-recorder/pkg/ui"
	"github.com/chromedp/chromedp"
)

func main() {
	// Print banner
	ui.PrintBanner()

	// --- Unified startup process with single progress indicator ---
	spinner := ui.NewSquareSpinner()
	spinner.Start("Initializing ChatGPT CLI...")

	// Browser setup
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true), // 
		chromedp.Flag("enable-automation", false), // Critical!
		chromedp.Flag("disable-extensions", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"), // Critical!
		chromedp.Flag("window-size", "1920,1080"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36`),
	)
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	ctx, ctxCancel := chromedp.NewContext(allocCtx)
	defer ctxCancel()

	// Load cookies
	spinner.Update("Loading saved session...")
	time.Sleep(500 * time.Millisecond) // Brief pause for smooth transition
	if err := chromedp.Run(ctx, browser.LoadCookiesAction()); err != nil {
		// Continue silently - cookies not critical
	}

	// Navigate to ChatGPT
	spinner.Update("Connecting to ChatGPT...")
	time.Sleep(300 * time.Millisecond) // Brief pause for smooth transition
	targetURL := config.BaseURL
	if err := chromedp.Run(ctx, chromedp.Navigate(targetURL)); err != nil {
		spinner.Stop()
		ui.PrintError("Failed to connect to ChatGPT")
		log.Fatalf("Navigation error: %v", err)
	}

	// Reload technique for stability
	spinner.Update("Optimizing connection...")
	time.Sleep(3 * time.Second)
	if err := chromedp.Run(ctx, chromedp.Reload()); err != nil {
		spinner.Stop()
		ui.PrintError("Connection optimization failed")
		log.Fatalf("Reload error: %v", err)
	}

	// Wait for ChatGPT to load
	spinner.Update("Verifying interface...")
	time.Sleep(300 * time.Millisecond) // Brief pause for smooth transition
	if err := chromedp.Run(ctx, browser.WaitForChatGPTLoad()); err != nil {
		spinner.Stop()
		ui.PrintWarning("Interface verification incomplete - please ensure you're logged in")
		ui.PrintInfo("You may need to login manually in the browser window")
		return
	}

	// Create ChatGPT client and final checks
	chatgptClient := chatgpt.NewChatGPT(ctx)
	spinner.Update("Finalizing setup...")
	time.Sleep(300 * time.Millisecond) // Brief pause for smooth transition
	if err := chatgptClient.WaitForPageLoad(); err != nil {
		// Continue anyway - not critical
	}

	spinner.Stop()
	ui.PrintSuccess("GPT5-DEV Agent CLI ready! 🚀")

	// Create and start CLI
	cliApp := cli.NewCLI(chatgptClient)

	// Start the CLI interface
	if err := cliApp.Start(); err != nil {
		ui.PrintError("CLI error occurred")
		log.Fatalf("CLI error: %v", err)
	}
}
