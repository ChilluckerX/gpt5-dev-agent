package main

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/chatgpt-element-recorder/pkg/browser"
	"github.com/chatgpt-element-recorder/pkg/chatgpt"
	"github.com/chatgpt-element-recorder/pkg/cli"
	"github.com/chatgpt-element-recorder/pkg/config"
	"github.com/chatgpt-element-recorder/pkg/ui"
)

func main() {
	// Print banner
	ui.PrintBanner()

	// --- Base setup for the browser (Go scraper style) ---
	spinner := ui.NewSquareSpinner()
	spinner.Start("Setting up browser...")
	
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // Keep visible for CLI interaction
		chromedp.Flag("enable-automation", false), // Critical!
		chromedp.Flag("disable-extensions", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"), // Critical!
		chromedp.Flag("window-size", "1920,1080"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36`),
	)
	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	// --- Create a single, long-lived browser context (Go scraper approach) ---
	ctx, ctxCancel := chromedp.NewContext(allocCtx)
	defer ctxCancel()
	
	spinner.Stop()
	ui.PrintSuccess("Browser setup complete")

	// --- 1. Initial Setup: Load existing ChatGPT cookies ---
	spinner = ui.NewDotSpinner()
	spinner.Start("Loading ChatGPT cookies...")
	
	if err := chromedp.Run(ctx, browser.LoadCookiesAction()); err != nil {
		spinner.Stop()
		ui.PrintWarning("Could not load cookies - you may need to login manually")
	} else {
		spinner.Stop()
		ui.PrintSuccess("Cookies loaded successfully")
	}

	// --- 2. Navigate to ChatGPT (Go scraper technique) ---
	spinner = ui.NewSquareSpinner()
	spinner.Start("Navigating to ChatGPT...")
	
	targetURL := config.BaseURL
	if err := chromedp.Run(ctx, chromedp.Navigate(targetURL)); err != nil {
		spinner.Stop()
		ui.PrintError("Failed to navigate to ChatGPT")
		log.Fatalf("Navigation error: %v", err)
	}

	// --- Add a reload step to handle potential blank page on first load (Go scraper trick!) ---
	spinner.Update("Applying Go scraper reload technique...")
	time.Sleep(3 * time.Second) // Wait a moment for the initial (potentially blank) page
	if err := chromedp.Run(ctx, chromedp.Reload()); err != nil {
		spinner.Stop()
		ui.PrintError("Failed to reload page")
		log.Fatalf("Reload error: %v", err)
	}

	// --- 3. Wait for ChatGPT to load ---
	spinner.Update("Waiting for ChatGPT to load...")
	if err := chromedp.Run(ctx, browser.WaitForChatGPTLoad()); err != nil {
		spinner.Stop()
		ui.PrintWarning("Could not confirm ChatGPT loaded properly")
		ui.PrintInfo("Please make sure you're logged in and can see the chat interface")
	} else {
		spinner.Stop()
		ui.PrintSuccess("ChatGPT loaded successfully")
	}

	// --- 4. Create ChatGPT client and CLI ---
	chatgptClient := chatgpt.NewChatGPT(ctx)
	
	// Wait for page to be ready
	spinner = ui.NewSpinner()
	spinner.Start("Checking page readiness...")
	if err := chatgptClient.WaitForPageLoad(); err != nil {
		spinner.Stop()
		ui.PrintWarning("Page load check failed - continuing anyway")
	} else {
		spinner.Stop()
		ui.PrintSuccess("Page is ready")
	}

	// Create and start CLI
	cliApp := cli.NewCLI(chatgptClient)
	
	ui.PrintSuccess("ChatGPT CLI is ready!")
	ui.PrintInfo("Browser window will stay open for interaction")
	
	// Start the CLI interface
	if err := cliApp.Start(); err != nil {
		ui.PrintError("CLI error occurred")
		log.Fatalf("CLI error: %v", err)
	}
}