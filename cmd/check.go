package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"urgent-reminder/internal/display"
	"urgent-reminder/internal/service"
	"urgent-reminder/internal/storage"
)

var checkCmd = &cobra.Command{
	Use:   "check [id]",
	Short: "Mark a reminder as complete",
	Long:  `Mark a reminder as complete. If recurrent, it will advance to the next cycle. If not recurrent, it will be deleted.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		idStr := args[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("invalid ID: %s", idStr)
		}

		store, err := storage.NewJSONStore()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}

		reminderService := service.NewReminderService(store)
		displayObj := display.NewDisplay(noColor)

		reminder, err := reminderService.GetReminder(id)
		if err != nil {
			return fmt.Errorf("failed to get reminder: %w", err)
		}

		if reminder.IsRecurrent {
			if err := reminderService.CheckReminder(id); err != nil {
				return fmt.Errorf("failed to update reminder: %w", err)
			}
			updatedReminder, err := reminderService.GetReminder(id)
			if err != nil {
				return fmt.Errorf("failed to get updated reminder: %w", err)
			}
			displayObj.PrintSuccess("✓ Recurrent reminder advanced to next cycle")
			displayObj.PrintEmpty()
			displayObj.PrintInfo(fmt.Sprintf("Next due date: %s", updatedReminder.DueDate.Format("2006-01-02")))
		} else {
			if err := reminderService.CheckReminder(id); err != nil {
				return fmt.Errorf("failed to delete reminder: %w", err)
			}
			displayObj.PrintSuccess("✓ Reminder completed and deleted")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
