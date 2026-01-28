package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/goldie/ellie-cli/internal/api"
	"github.com/goldie/ellie-cli/internal/models"
	"github.com/spf13/cobra"
)

var listsCmd = &cobra.Command{
	Use:   "lists",
	Short: "Manage lists",
}

var listListsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all lists",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		lists, err := client.GetLists()
		if err != nil {
			return err
		}

		return outputLists(lists)
	},
}

func init() {
	listsCmd.AddCommand(listListsCmd)
}

func outputLists(lists []models.List) error {
	if IsJSONOutput() {
		data, err := json.MarshalIndent(lists, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	if len(lists) == 0 {
		fmt.Println("No lists found")
		return nil
	}

	for _, list := range lists {
		printList(&list)
	}
	return nil
}

func printList(list *models.List) {
	if list.Icon != "" {
		fmt.Printf("%s %s\n", list.Icon, list.Title)
	} else {
		fmt.Printf("• %s\n", list.Title)
	}
	fmt.Printf("  ID: %s\n", list.ID)
}
