# Urgent Reminder Terminal

A modern CLI tool to manage urgent reminders with date-based alerts. Displays reminders with ASCII art in your terminal and integrates with zsh for automatic alerts.

## Features

 - 📅 Add reminders with due dates
 - 📋 List all reminders sorted by due date
 - 🔔 Check for active reminders (within 7 days)
 - ✅ Complete reminders to stop alerts
 - 🔄 Reset reminders for next cycle
 - 🎨 ASCII art display with colored output
 - 💾 XDG-compliant data storage (cross-platform)
 - 🐚 Auto-setup for zsh and bash integration

## Installation

### Prerequisites

- Go 1.16 or later
- Task (go-task) - optional, for build automation

### Install Go

```bash
# macOS
brew install go

# Linux
sudo apt-get install golang  # Ubuntu/Debian
sudo yum install golang      # CentOS/RHEL
```

### Install Task

```bash
# macOS
brew install go-task/tap/go-task

# Linux
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d
```

### Install Urgent Reminder

```bash
# Clone the repository
git clone https://github.com/yourusername/urgent-reminder-terminal.git
cd urgent-reminder-terminal

# Install and auto-setup shell integration
task install-local

# Or install without shell setup
task install

# Or build manually
go build -o urgent-reminder .
sudo install -m 0755 urgent-reminder /usr/local/bin/urgent-reminder
```

### Verify Installation

```bash
urgent-reminder --help
```

## Usage

### Add a Reminder

```bash
urgent-reminder add "Submit project report" 2026-01-15
urgent-reminder add "Pay credit card bill" 2026-01-20
urgent-reminder add "Renew insurance" 2026-02-01
```

### List All Reminders

```bash
urgent-reminder list
```

Output:

```
  _   _          _   _          __        __                 _       _ 
 | | | | ___  | | | |   ___   \ \      / /__  _ __  | |   __| |
 | |_| |/ _ \ | | | |  / _ \   \ \ /\ / / _ \| '__| | |  / _` |
 |  _  |  __/ | | | | | (_) |   \ V  V / (_) | |    | | | (_| |
 |_| |_|\___| |_| |_|  \___/     \_/\_/ \___/|_|    |_|  \__,_|

============================================================

[ACTIVE] Submit project report -- Due in 5 days
[ACTIVE] Pay credit card bill -- Due in 10 days
[ACTIVE] Renew insurance -- Due in 22 days

============================================================
Total: 3 reminder(s)
Data: /home/user/.local/share/urgent-reminder/reminders.json
```

### Check for Active Reminders

```bash
urgent-reminder check
```

This command:

- Displays reminders due within 7 days
- Exits with code 1 if reminders found (for zsh integration)
- Exits with code 0 if no reminders (silent)

### Complete a Reminder

```bash
# Complete by ID
urgent-reminder complete 20250102120000

# Complete all reminders
urgent-reminder complete all
```

### Reset a Reminder

```bash
# Reset by ID (moves to next year and re-enables alerts)
urgent-reminder reset 20250102120000

# Reset all reminders
urgent-reminder reset all
```

### View Configuration

```bash
urgent-reminder config
```

Shows:

- Platform information
- Data file location
- XDG environment variables
- Customization options

### Setup Shell Integration

```bash
urgent-reminder setup
```

This command automatically:

- Detects your shell (zsh or bash)
- Adds integration to your config file (`~/.zshrc` or `~/.bashrc`)
- Prevents duplicate entries
- Displays configuration details

Or use `task install-local` to install the binary and set up shell integration in one step.

The integration adds this function to your shell config:

```bash
# Urgent Reminder Integration
urgent_reminder_list() {
    if command -v urgent-reminder &>/dev/null; then
        urgent-reminder list
    fi
}
urgent_reminder_list
```

Apply changes:

```bash
# For zsh
source ~/.zshrc

# For bash
source ~/.bashrc

# Or restart your terminal
```

## Data Storage

### XDG-Compliant Storage

This tool follows the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html) for cross-platform compatibility.

Default locations:

- **Linux/macOS**: `~/.local/share/urgent-reminder/reminders.json`
- **Windows**: Not yet supported

### Customization

Change data location:

```bash
export XDG_DATA_HOME=/custom/path
```

## Environment Variables

### Colors

Disable colored output:

```bash
# Method 1: Environment variable
export URGENT_REMINDER_NO_COLOR=1

# Method 2: Command flag
urgent-reminder list --no-color
```

### XDG Paths

```bash
export XDG_DATA_HOME=/custom/data/path
export XDG_CONFIG_HOME=/custom/config/path
export XDG_CACHE_HOME=/custom/cache/path
```

## Task Commands

If you have Task installed:

```bash
# Build the binary
task build

# Install to system
task install

# Install and auto-setup shell integration
task install-local

# Setup shell integration manually
urgent-reminder setup

# Clean build artifacts
task clean

# Run in development mode
task dev

# Run linters
task lint

# Run tests
task test
```

## Cross-Platform Compatibility

This tool is designed to work identically on macOS and Linux by following the XDG Base Directory Specification.

### macOS

On modern macOS systems, XDG paths work seamlessly:

- Data: `~/.local/share/urgent-reminder/`
- Config: `~/.config/`
- Cache: `~/.cache/`

### Linux

Standard XDG paths:

- Data: `~/.local/share/urgent-reminder/` (or `$XDG_DATA_HOME`)
- Config: `~/.config/` (or `$XDG_CONFIG_HOME`)
- Cache: `~/.cache/` (or `$XDG_CACHE_HOME`)

## Troubleshooting

### "command not found: urgent-reminder"

**Solution**: Ensure the binary is in your PATH:

```bash
which urgent-reminder
# If not found, check:
ls -la /usr/local/bin/urgent-reminder
ls -la $GOPATH/bin/urgent-reminder
```

### Shell integration not working

**Solution**: Verify the function is in your shell config:

```bash
# For zsh
grep "urgent_reminder_list" ~/.zshrc

# For bash
grep "urgent_reminder_list" ~/.bashrc
```

If not present, run:

```bash
urgent-reminder setup
```

Then reload:

```bash
# For zsh
source ~/.zshrc

# For bash
source ~/.bashrc
```

### Reminders not displaying in terminal

**Solution**: Check that the check command works:

```bash
urgent-reminder check
```

If it works manually but not in new terminals:

- Verify zsh integration is set up correctly
- Check that errors are suppressed (using `2>/dev/null`)

### Date format errors

**Solution**: Always use YYYY-MM-DD format:

```bash
# Correct
urgent-reminder add "Task" 2026-01-15

# Incorrect
urgent-reminder add "Task" 01/15/2026
urgent-reminder add "Task" Jan 15, 2026
```

### Permission denied when writing data

**Solution**: Check XDG directory permissions:

```bash
ls -la ~/.local/share/
```

If it doesn't exist, create it:

```bash
mkdir -p ~/.local/share/urgent-reminder
```

## Development

### Project Structure

```
.
├── cmd/
│   ├── root.go          # Root command
│   ├── add.go           # Add reminder command
│   ├── list.go          # List reminders command
│   ├── check.go         # Check reminders command
│   ├── complete.go      # Complete/reset commands
│   ├── setup.go         # Zsh integration setup
│   └── config.go        # Configuration display
├── internal/
│   ├── models/
│   │   └── reminder.go      # Reminder data model
│   ├── storage/
│   │   └── json_store.go    # XDG-compliant JSON storage
│   ├── service/
│   │   └── reminder_service.go  # Business logic
│   └── display/
│       └── banner.go            # ASCII art display
├── main.go                   # Entry point
├── Taskfile.yml              # Build automation
├── go.mod                    # Go module
└── README.md                 # This file
```

### Building

```bash
# Development
task dev

# Production build
task build
task install
```

### Dependencies

- [github.com/spf13/cobra](https://github.com/spf13/cobra) - CLI framework
- [github.com/common-nighthawk/go-figure](https://github.com/common-nighthawk/go-figure) - ASCII art
- [github.com/fatih/color](https://github.com/fatih/color) - Colored output

## License

MIT License - feel free to use this tool for personal or commercial projects.

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## Credits

Built with:

- Go programming language
- Cobra CLI framework
- XDG Base Directory Specification
- Go-task for build automation
