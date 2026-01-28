package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/goldie/ellie-cli/internal/api"
	"github.com/goldie/ellie-cli/internal/models"
	"github.com/spf13/cobra"
)

var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "User operations",
}

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Show current user",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		user, err := client.GetCurrentUser()
		if err != nil {
			return err
		}

		return outputUser(user)
	},
}

var usageCmd = &cobra.Command{
	Use:   "usage",
	Short: "Show API usage statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		usage, err := client.GetAPIUsage()
		if err != nil {
			return err
		}

		return outputUsage(usage)
	},
}

func init() {
	usersCmd.AddCommand(meCmd)
	usersCmd.AddCommand(usageCmd)
}

func outputUser(user *models.User) error {
	if IsJSONOutput() {
		data, err := json.MarshalIndent(user, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	fmt.Printf("Name:   %s\n", user.Name)
	fmt.Printf("Email:  %s\n", user.Email)
	fmt.Printf("ID:     %s\n", user.ID)
	return nil
}

func outputUsage(usage *models.APIUsage) error {
	if IsJSONOutput() {
		data, err := json.MarshalIndent(usage, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	fmt.Printf("Date:      %s\n", usage.Today.Date)
	fmt.Printf("Used:      %d / %d requests\n", usage.Today.Used, usage.Today.Limit)
	fmt.Printf("Remaining: %d\n", usage.Today.Remaining)
	fmt.Printf("Resets:    %s\n", usage.ResetAt)
	fmt.Printf("Rate:      %d req/min\n", usage.RateLimit.RequestsPerMinute)
	return nil
}
