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

## 💡 Usage Tips
- **Valid cookies** - Ensure `cookies/chatgpt.json` has valid session
- **Manual login** - Login manually first if cookies expired
- **Wait for responses** - CLI waits for ChatGPT to respond

## 🔧 Troubleshooting

### Browser Issues:
```bash
# Make sure Chrome is installed and updated
# Check if cookies are valid
```

### Connection Issues:
```bash
# Check internet connection
# Verify ChatGPT is accessible
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

## 🤝 Contributing

We welcome contributions to improve this ChatGPT CLI tool! Whether you're fixing bugs, adding new features, or improving documentation, your help is appreciated.

### How to Contribute:

1. **Fork the repository** on GitHub
2. **Create a feature branch**: `git checkout -b feature/your-feature-name`
3. **Make your changes** and test them thoroughly
4. **Commit your changes**: `git commit -m "Add your descriptive commit message"`
5. **Push to your branch**: `git push origin feature/your-feature-name`
6. **Open a Pull Request** with a clear description of your changes

### What We're Looking For:

- 🐛 **Bug fixes** - Help us squash those pesky bugs
- ✨ **New features** - Enhance the CLI experience
- 📚 **Documentation** - Improve README, code comments, or add examples
- 🔧 **Performance improvements** - Make it faster and more efficient
- 🧪 **Tests** - Add test coverage for better reliability

### Guidelines:

- Follow Go best practices and conventions
- Test your changes before submitting
- Keep commits focused and descriptive
- Update documentation if needed

### Questions or Ideas?

Feel free to open an issue to discuss new features, report bugs, or ask questions. We're happy to help and collaborate!

---
