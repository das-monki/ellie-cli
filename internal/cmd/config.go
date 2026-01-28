package cmd

import (
	"fmt"
	"strings"

	"github.com/goldie/ellie-cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage CLI configuration",
}

var setAPIKeyCmd = &cobra.Command{
	Use:   "set-api-key <key>",
	Short: "Set the API key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := strings.TrimSpace(args[0])
		if apiKey == "" {
			return fmt.Errorf("API key cannot be empty")
		}

		if err := config.SetAPIKey(apiKey); err != nil {
			return fmt.Errorf("failed to save API key: %w", err)
		}

		configDir, _ := config.GetConfigDir()
		fmt.Printf("API key saved to %s/config.yaml\n", configDir)
		return nil
	},
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.GetConfig()

		fmt.Println("Configuration:")
		fmt.Printf("  Base URL: %s\n", cfg.BaseURL)

		if cfg.APIKey != "" {
			// Mask the API key, showing only first 8 and last 4 characters
			masked := maskAPIKey(cfg.APIKey)
			fmt.Printf("  API Key:  %s\n", masked)
		} else {
			fmt.Println("  API Key:  (not set)")
		}

		configDir, _ := config.GetConfigDir()
		fmt.Printf("\nConfig file: %s/config.yaml\n", configDir)
		fmt.Println("\nAPI key priority:")
		fmt.Println("  1. ELLIE_API_KEY environment variable")
		fmt.Println("  2. ELLIE_API_KEY_FILE environment variable (path to file)")
		fmt.Println("  3. Config file")

		return nil
	},
}

var setBaseURLCmd = &cobra.Command{
	Use:   "set-base-url <url>",
	Short: "Set the API base URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		baseURL := strings.TrimSpace(args[0])
		if baseURL == "" {
			return fmt.Errorf("base URL cannot be empty")
		}

		if err := config.SetBaseURL(baseURL); err != nil {
			return fmt.Errorf("failed to save base URL: %w", err)
		}

		fmt.Printf("Base URL set to: %s\n", baseURL)
		return nil
	},
}

func init() {
	configCmd.AddCommand(setAPIKeyCmd)
	configCmd.AddCommand(showConfigCmd)
	configCmd.AddCommand(setBaseURLCmd)
}

func maskAPIKey(key string) string {
	if len(key) <= 12 {
		return strings.Repeat("*", len(key))
	}
	return key[:8] + strings.Repeat("*", len(key)-12) + key[len(key)-4:]
}
