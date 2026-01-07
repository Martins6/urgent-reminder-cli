package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"urgent-reminder/internal/display"
	"urgent-reminder/internal/storage"
)

var configListCmd = &cobra.Command{
	Use:   "config-list",
	Short: "List config file locations",
	Long:  `Show the locations of configuration and data files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewJSONStore()
		if err != nil {
			return fmt.Errorf("failed to initialize storage: %w", err)
		}

		displayObj := display.NewDisplay(noColor)

		displayObj.PrintHeader("Configuration Files")
		displayObj.PrintEmpty()
		displayObj.PrintInfo(fmt.Sprintf("Data file: %s", store.GetDataPath()))

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		dataHome := os.Getenv("XDG_DATA_HOME")
		if dataHome == "" {
			dataHome = filepath.Join(homeDir, ".local", "share")
		}
		appDataPath := filepath.Join(dataHome, "urgent-reminder")
		displayObj.PrintInfo(fmt.Sprintf("Data directory: %s", appDataPath))

		displayObj.PrintEmpty()
		displayObj.PrintInfo("To change data location, set XDG_DATA_HOME:")
		displayObj.PrintInfo("  export XDG_DATA_HOME=/custom/path")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configListCmd)
}
