package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"urgent-reminder/internal/display"
)

const (
	zshIntegration = `# Urgent Reminder Integration
urgent_reminder_list() {
    if command -v urgent-reminder &>/dev/null; then
        urgent-reminder list 2>/dev/null | grep -q "Total:" && urgent-reminder list
    fi
}
urgent_reminder_list
`

	bashIntegration = `# Urgent Reminder Integration
urgent_reminder_list() {
    if command -v urgent-reminder &>/dev/null; then
        urgent-reminder list 2>/dev/null | grep -q "Total:" && urgent-reminder list
    fi
}
urgent_reminder_list
`
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup shell integration for automatic reminder display",
	Long:  `Setup shell integration that automatically displays reminders when you open a new terminal. Supports both zsh and bash shells.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		displayObj := display.NewDisplay(noColor)

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		shell := os.Getenv("SHELL")
		var configFile string
		var shellName string

		if strings.Contains(shell, "zsh") || fileExists(homeDir+"/.zshrc") {
			configFile = homeDir + "/.zshrc"
			shellName = "zsh"
		} else {
			configFile = homeDir + "/.bashrc"
			shellName = "bash"
		}

		displayObj.PrintInfo(fmt.Sprintf("Detected shell: %s", shellName))
		displayObj.PrintInfo(fmt.Sprintf("Config file: %s", configFile))
		displayObj.PrintEmpty()

		content, err := os.ReadFile(configFile)
		if err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("failed to read config file: %w", err)
			}
			content = []byte{}
		}

		contentStr := string(content)

		if strings.Contains(contentStr, "urgent_reminder_list") {
			displayObj.PrintWarning("Shell integration already configured!")
			displayObj.PrintEmpty()
			displayObj.PrintInfo("The urgent_reminder_list function is already in your shell config.")
			displayObj.PrintEmpty()
			displayObj.PrintInfo("To apply changes (if you just set this up manually):")
			displayObj.PrintInfo(fmt.Sprintf("  source ~/.%src", shellName))
			displayObj.PrintInfo("  # or restart your terminal")
			return nil
		}

		var integration string
		if shellName == "zsh" {
			integration = zshIntegration
		} else {
			integration = bashIntegration
		}

		if contentStr != "" && !strings.HasSuffix(contentStr, "\n") {
			integration = "\n" + integration
		}

		file, err := os.OpenFile(configFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return fmt.Errorf("failed to open config file: %w", err)
		}
		defer file.Close()

		if _, err := file.WriteString(integration); err != nil {
			return fmt.Errorf("failed to write to config file: %w", err)
		}

		displayObj.PrintSuccess("✓ Shell integration configured successfully!")
		displayObj.PrintEmpty()
		displayObj.PrintInfo("The urgent_reminder_list function has been added to your shell config.")
		displayObj.PrintEmpty()
		displayObj.PrintInfo("To apply changes:")
		displayObj.PrintInfo(fmt.Sprintf("  source ~/.%src", shellName))
		displayObj.PrintInfo("  # or restart your terminal")
		displayObj.PrintEmpty()
		displayObj.PrintInfo("This will automatically run 'urgent-reminder list' in new terminals.")

		return nil
	},
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
