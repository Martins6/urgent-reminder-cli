package cmd

import (
	"fmt"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"urgent-reminder/internal/display"
	"urgent-reminder/internal/models"
	"urgent-reminder/internal/service"
	"urgent-reminder/internal/storage"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new reminder",
	Long:  `Add a new reminder with interactive prompts for title, date, and recurrence options.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewJSONStore()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}

		reminderService := service.NewReminderService(store)
		displayObj := display.NewDisplay(noColor)

		titlePrompt := promptui.Prompt{
			Label: "Title",
			Validate: func(input string) error {
				if input == "" {
					return fmt.Errorf("title cannot be empty")
				}
				return nil
			},
		}
		title, err := titlePrompt.Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}

		recurrentPrompt := promptui.Select{
			Label: "Is this reminder recurrent?",
			Items: []string{"No", "Yes"},
		}
		_, isRecurrent, err := recurrentPrompt.Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}

		nextID, err := reminderService.GetNextID()
		if err != nil {
			return fmt.Errorf("failed to get next ID: %w", err)
		}

		var reminder *models.Reminder

		if isRecurrent == "No" {
			datePrompt := promptui.Prompt{
				Label: "Date (YYYY-MM-DD)",
				Validate: func(input string) error {
					_, err := time.Parse("2006-01-02", input)
					if err != nil {
						return fmt.Errorf("invalid date format, use YYYY-MM-DD")
					}
					return nil
				},
			}
			dateStr, err := datePrompt.Run()
			if err != nil {
				return fmt.Errorf("prompt failed: %w", err)
			}

			dueDate, _ := time.Parse("2006-01-02", dateStr)

			timePrompt := promptui.Prompt{
				Label:     "Time (HH:MM, optional, press Enter to skip)",
				IsConfirm: false,
			}
			timeStr, err := timePrompt.Run()
			if err != nil {
				return fmt.Errorf("prompt failed: %w", err)
			}

			reminder = models.NewReminder(nextID, title, dueDate)
			if timeStr != "" {
				_, err := time.Parse("15:04", timeStr)
				if err != nil {
					return fmt.Errorf("invalid time format, use HH:MM")
				}
				reminder.Time = timeStr
			}
		} else {
			recurrentTypePrompt := promptui.Select{
				Label: "Recurrence type",
				Items: []string{"Weekly", "Bi-weekly", "Monthly"},
			}
			_, recurrentTypeStr, err := recurrentTypePrompt.Run()
			if err != nil {
				return fmt.Errorf("prompt failed: %w", err)
			}

			var recurrentType models.RecurrentType
			switch recurrentTypeStr {
			case "Weekly":
				recurrentType = models.RecurrentWeekly
			case "Bi-weekly":
				recurrentType = models.RecurrentBiWeekly
			case "Monthly":
				recurrentType = models.RecurrentMonthly
			}

			datePrompt := promptui.Prompt{
				Label: "Start date (YYYY-MM-DD)",
				Validate: func(input string) error {
					_, err := time.Parse("2006-01-02", input)
					if err != nil {
						return fmt.Errorf("invalid date format, use YYYY-MM-DD")
					}
					return nil
				},
			}
			dateStr, err := datePrompt.Run()
			if err != nil {
				return fmt.Errorf("prompt failed: %w", err)
			}

			dueDate, _ := time.Parse("2006-01-02", dateStr)

			reminder = models.NewRecurrentReminder(nextID, title, dueDate, recurrentType)

			if recurrentType == models.RecurrentWeekly || recurrentType == models.RecurrentBiWeekly {
				dayPrompt := promptui.Select{
					Label: "Select days (multi-select)",
					Items: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
				}
				_, day, err := dayPrompt.Run()
				if err != nil {
					return fmt.Errorf("prompt failed: %w", err)
				}
				reminder.RecurrentDays = append(reminder.RecurrentDays, day)

				continuePrompt := promptui.Select{
					Label: "Add more days?",
					Items: []string{"No", "Yes"},
				}
				_, cont, err := continuePrompt.Run()
				for cont == "Yes" {
					_, day, err := dayPrompt.Run()
					if err != nil {
						return fmt.Errorf("prompt failed: %w", err)
					}
					reminder.RecurrentDays = append(reminder.RecurrentDays, day)
					_, cont, err = continuePrompt.Run()
				}
			} else if recurrentType == models.RecurrentMonthly {
				dayOfMonthPrompt := promptui.Prompt{
					Label: "Day of month (1-31)",
					Validate: func(input string) error {
						var day int
						_, err := fmt.Sscanf(input, "%d", &day)
						if err != nil || day < 1 || day > 31 {
							return fmt.Errorf("enter a number between 1 and 31")
						}
						return nil
					},
				}
				dayStr, err := dayOfMonthPrompt.Run()
				if err != nil {
					return fmt.Errorf("prompt failed: %w", err)
				}
				var dayOfMonth int
				fmt.Sscanf(dayStr, "%d", &dayOfMonth)
				reminder.RecurrentDayOfMonth = dayOfMonth
			}

			timePrompt := promptui.Prompt{
				Label:     "Time (HH:MM, optional, press Enter to skip)",
				IsConfirm: false,
			}
			timeStr, err := timePrompt.Run()
			if err != nil {
				return fmt.Errorf("prompt failed: %w", err)
			}
			if timeStr != "" {
				_, err := time.Parse("15:04", timeStr)
				if err != nil {
					return fmt.Errorf("invalid time format, use HH:MM")
				}
				reminder.Time = timeStr
			}
		}

		if err := reminderService.AddReminder(reminder); err != nil {
			return fmt.Errorf("failed to add reminder: %w", err)
		}

		displayObj.PrintSuccess("✓ Reminder added successfully!")
		displayObj.PrintEmpty()
		displayObj.PrintInfo(fmt.Sprintf("ID: %d", reminder.ID))
		displayObj.PrintInfo(fmt.Sprintf("Title: %s", reminder.Title))
		displayObj.PrintInfo(fmt.Sprintf("Date: %s", reminder.FormatDueDate()))
		if reminder.Time != "" {
			displayObj.PrintInfo(fmt.Sprintf("Time: %s", reminder.Time))
		}
		if reminder.IsRecurrent {
			displayObj.PrintInfo(fmt.Sprintf("Recurrent: %s", reminder.RecurrentType))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
