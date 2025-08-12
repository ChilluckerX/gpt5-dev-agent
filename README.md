# ChatGPT CLI - Go Edition

🚀 **Command-line interface untuk ChatGPT menggunakan Go scraper techniques**

Menggunakan selectors yang manually fetched untuk reliable automation!

## 🇲🇾 Features

- ✅ **Go scraper bypass** - Cloudflare detection avoidance
- ✅ **Real CLI interface** - Interactive command-line experience  
- ✅ **Chat history** - Browse and open previous conversations
- ✅ **New chat** - Start fresh conversations
- ✅ **Manual selectors** - Uses proven selectors from manual_fetch.txt

## 🛠️ Setup

### 1. Install Go
```bash
# Download from https://golang.org/dl/
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Ensure Chrome Browser & ChatGPT Cookies
- Chrome browser installed
- Valid ChatGPT cookies in `cookies/chatgpt.json`

## 🚀 Usage

### Run the CLI:
```bash
go run main.go
```

### CLI Commands:

| Command | Description |
|---------|-------------|
| `/help`, `/h` | Show help |
| `/new`, `/n` | Start new chat |
| `/history`, `/hist` | Show chat history |
| `/open <id>`, `/o <id>` | Open specific chat |
| `/clear`, `/cls` | Clear screen |
| `/quit`, `/q`, `/exit` | Exit CLI |

### Examples:

```bash
💬 ChatGPT CLI> Hello, how are you?
💬 ChatGPT CLI> /new
💬 ChatGPT CLI> /history
💬 ChatGPT CLI> /open 1
💬 ChatGPT CLI> /open 689916e6-3df0-8331-8eb6-e6f0c648cea4
💬 ChatGPT CLI> /quit
```

## 📁 Project Structure

```
├── main.go                     # Main CLI application
├── go.mod                     # Go modules
├── cookies/
│   └── chatgpt.json          # ChatGPT session cookies
├── pkg/
│   ├── browser/              # Browser automation
│   ├── chatgpt/              # ChatGPT interaction logic
│   ├── cli/                  # CLI interface
│   ├── config/               # Configuration
│   └── file/                 # File utilities
└── manual_fetch.txt          # Manual selectors reference
```

## 🎯 Manual Selectors Used

From `manual_fetch.txt`:

- **Input**: `#prompt-textarea`
- **Submit**: `#composer-submit-button` 
- **Response**: `#thread > div > div.relative...`
- **New Chat**: `#stage-slideover-sidebar > div > div.opacity...`
- **History**: `#history`

## 🔑 Key Go Scraper Techniques

### 1. **Anti-Automation Flags**
```go
chromedp.Flag("enable-automation", false)
chromedp.Flag("disable-blink-features", "AutomationControlled")
```

### 2. **Reload Trick**
```go
chromedp.Navigate(targetURL)
time.Sleep(3 * time.Second)
chromedp.Reload() // Handle challenges
```

### 3. **Cookie Loading**
```go
browser.LoadCookiesAction() // Load real session
```

### 4. **XPath Selectors**
```go
chromedp.SendKeys(`//*[@id="prompt-textarea"]`, message, chromedp.BySearch)
chromedp.Click(`//*[@id="composer-submit-button"]`, chromedp.BySearch)
```

## 💡 Usage Tips

- **Keep browser visible** - CLI needs browser window open
- **Valid cookies** - Ensure `cookies/chatgpt.json` has valid session
- **Manual login** - Login manually first if cookies expired
- **Wait for responses** - CLI waits for ChatGPT to respond

## 🔧 Troubleshooting

### Browser Issues:
```bash
# Make sure Chrome is installed and updated
# Check if cookies are valid
```

### Selector Issues:
```bash
# Selectors are from manual_fetch.txt
# If ChatGPT UI changes, update selectors
```

### Connection Issues:
```bash
# Check internet connection
# Verify ChatGPT is accessible
# Try manual login first
```

## 🆚 Advantages Over Python

| Feature | Python Selenium | Go ChromeDP |
|---------|----------------|-------------|
| **Performance** | Slower | Faster |
| **Memory** | Higher usage | Lower usage |
| **Detection** | Easier to detect | Harder to detect |
| **Stability** | Less stable | More stable |
| **Dependencies** | Many | Minimal |

## 🎯 CLI Workflow

1. **Start CLI** - `go run main.go`
2. **Browser opens** - ChatGPT loads with cookies
3. **Interactive prompt** - Type messages or commands
4. **Real-time chat** - Send/receive messages
5. **History access** - Browse previous chats
6. **New conversations** - Start fresh anytime

## 🇲🇾 Success Factors

- **Manual selectors** - Proven to work
- **Go scraper techniques** - Bypass detection
- **Real cookies** - Maintain session
- **CLI interface** - Easy to use
- **Reliable automation** - Stable performance

---

**🚀 Ready untuk real ChatGPT CLI automation dengan Go power!**