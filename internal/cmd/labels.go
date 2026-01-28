package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/goldie/ellie-cli/internal/api"
	"github.com/goldie/ellie-cli/internal/models"
	"github.com/spf13/cobra"
)

var labelsCmd = &cobra.Command{
	Use:   "labels",
	Short: "Manage labels",
}

var listLabelsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all labels",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		labels, err := client.GetLabels()
		if err != nil {
			return err
		}

		return outputLabels(labels)
	},
}

var createLabelCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new label",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		color, _ := cmd.Flags().GetString("color")

		if name == "" {
			return fmt.Errorf("--name flag is required")
		}
		if color == "" {
			return fmt.Errorf("--color flag is required")
		}

		req := &models.CreateLabelRequest{
			Name:  name,
			Color: color,
		}

		client, err := api.NewClient()
		if err != nil {
			return err
		}

		label, err := client.CreateLabel(req)
		if err != nil {
			return err
		}

		return outputLabel(label)
	},
}

func init() {
	createLabelCmd.Flags().String("name", "", "Label name (required)")
	createLabelCmd.Flags().String("color", "", "Label color in hex format, e.g., #FF5733 (required)")

	labelsCmd.AddCommand(listLabelsCmd)
	labelsCmd.AddCommand(createLabelCmd)
}

func outputLabel(label *models.Label) error {
	if IsJSONOutput() {
		data, err := json.MarshalIndent(label, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	printLabel(label)
	return nil
}

func outputLabels(labels []models.Label) error {
	if IsJSONOutput() {
		data, err := json.MarshalIndent(labels, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	if len(labels) == 0 {
		fmt.Println("No labels found")
		return nil
	}

	for _, label := range labels {
		printLabel(&label)
	}
	return nil
}

func printLabel(label *models.Label) {
	fmt.Printf("• %s (%s)\n", label.Name, label.Color)
	fmt.Printf("  ID: %s\n", label.ID)
}
