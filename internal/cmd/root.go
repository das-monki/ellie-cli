package cmd

import (
	"os"

	"github.com/goldie/ellie-cli/internal/config"
	"github.com/spf13/cobra"
)

var jsonOutput bool

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:          "ellie",
	Short:        "CLI for the Ellie Daily Planner",
	Long:         `A command-line interface for interacting with the Ellie Daily Planner API.`,
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return config.Init()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")

	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(tasksCmd)
	rootCmd.AddCommand(labelsCmd)
	rootCmd.AddCommand(listsCmd)
	rootCmd.AddCommand(usersCmd)
}

// IsJSONOutput returns whether JSON output is enabled
func IsJSONOutput() bool {
	return jsonOutput
}
