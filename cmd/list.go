package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"urgent-reminder/internal/display"
	"urgent-reminder/internal/service"
	"urgent-reminder/internal/storage"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List due reminders",
	Long:  `List all reminders that are due or overdue.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewJSONStore()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}

		reminderService := service.NewReminderService(store)
		displayObj := display.NewDisplay(noColor)

		reminders, err := reminderService.GetDueReminders()
		if err != nil {
			return fmt.Errorf("failed to list reminders: %w", err)
		}

		if len(reminders) == 0 {
			displayObj.PrintInfo("No due reminders found.")
			return nil
		}

		sort.Slice(reminders, func(i, j int) bool {
			return reminders[i].ID < reminders[j].ID
		})

		displayObj.PrintBanner()
		displayObj.PrintEmpty()

		for _, reminder := range reminders {
			displayObj.PrintSimpleReminder(
				reminder.ID,
				reminder.Title,
				reminder.FormatDueDate(),
				reminder.FormatTime(),
			)
		}

		displayObj.PrintEmpty()
		displayObj.PrintInfo(fmt.Sprintf("Total: %d URGENT REMINDER(S)", len(reminders)))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
