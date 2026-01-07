package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	noColor bool
)

var rootCmd = &cobra.Command{
	Use:   "urgent-reminder",
	Short: "A simple CLI tool to manage reminders",
	Long: `Urgent Reminder is a simple CLI tool to manage reminders with due dates.
Supports both single and recurrent reminders (weekly, bi-weekly, monthly).

Data is stored in XDG-compliant locations:
  - Linux/macOS: ~/.local/share/urgent-reminder/

Commands:
  add         - Add a new reminder (interactive)
  list        - List due reminders
  check [id]  - Mark a reminder as complete
  config-list - List config file locations
  setup       - Setup shell integration`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")

	if os.Getenv("URGENT_REMINDER_NO_COLOR") == "1" {
		noColor = true
	}
}
