package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/goldie/ellie-cli/internal/api"
	"github.com/goldie/ellie-cli/internal/models"
	"github.com/spf13/cobra"
)

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "Manage tasks",
}

var getTaskCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get a task by ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		task, err := client.GetTask(args[0])
		if err != nil {
			return err
		}

		return outputTask(task)
	},
}

var listTasksCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks for a date",
	RunE: func(cmd *cobra.Command, args []string) error {
		date, _ := cmd.Flags().GetString("date")
		timeZone, _ := cmd.Flags().GetString("timezone")

		if date == "" {
			return fmt.Errorf("--date flag is required")
		}

		client, err := api.NewClient()
		if err != nil {
			return err
		}

		tasks, err := client.GetTasksByDate(date, timeZone)
		if err != nil {
			return err
		}

		return outputTasks(tasks)
	},
}

var byListCmd = &cobra.Command{
	Use:   "by-list",
	Short: "List tasks by list ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		listID, _ := cmd.Flags().GetString("list-id")

		if listID == "" {
			return fmt.Errorf("--list-id flag is required")
		}

		client, err := api.NewClient()
		if err != nil {
			return err
		}

		tasks, err := client.GetTasksByList(listID)
		if err != nil {
			return err
		}

		return outputTasks(tasks)
	},
}

var braindumpCmd = &cobra.Command{
	Use:   "braindump",
	Short: "Get unscheduled tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		tasks, err := client.GetBraindump()
		if err != nil {
			return err
		}

		return outputTasks(tasks)
	},
}

var createTaskCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new task",
	RunE: func(cmd *cobra.Command, args []string) error {
		desc, _ := cmd.Flags().GetString("desc")
		date, _ := cmd.Flags().GetString("date")
		start, _ := cmd.Flags().GetString("start")
		estimatedTime, _ := cmd.Flags().GetInt("estimated-time")
		listID, _ := cmd.Flags().GetString("list-id")
		label, _ := cmd.Flags().GetString("label")
		priority, _ := cmd.Flags().GetInt("priority")

		if desc == "" {
			return fmt.Errorf("--desc flag is required")
		}

		req := &models.CreateTaskRequest{
			Description: desc,
		}

		if date != "" {
			req.Date = &date
		}
		if start != "" {
			req.Start = &start
		}
		if estimatedTime > 0 {
			req.EstimatedTime = &estimatedTime
		}
		if listID != "" {
			req.ListID = &listID
		}
		if label != "" {
			req.Label = &label
		}
		if priority > 0 {
			req.Priority = &priority
		}

		client, err := api.NewClient()
		if err != nil {
			return err
		}

		task, err := client.CreateTask(req)
		if err != nil {
			return err
		}

		return outputTask(task)
	},
}

var updateTaskCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]

		desc, _ := cmd.Flags().GetString("desc")
		date, _ := cmd.Flags().GetString("date")
		start, _ := cmd.Flags().GetString("start")
		estimatedTime, _ := cmd.Flags().GetInt("estimated-time")
		complete, _ := cmd.Flags().GetBool("complete")
		listID, _ := cmd.Flags().GetString("list-id")
		label, _ := cmd.Flags().GetString("label")
		priority, _ := cmd.Flags().GetInt("priority")

		req := &models.UpdateTaskRequest{}

		if cmd.Flags().Changed("desc") {
			req.Description = &desc
		}
		if cmd.Flags().Changed("date") {
			req.Date = &date
		}
		if cmd.Flags().Changed("start") {
			req.Start = &start
		}
		if cmd.Flags().Changed("estimated-time") {
			req.EstimatedTime = &estimatedTime
		}
		if cmd.Flags().Changed("complete") {
			req.Complete = &complete
		}
		if cmd.Flags().Changed("list-id") {
			req.ListID = &listID
		}
		if cmd.Flags().Changed("label") {
			req.Label = &label
		}
		if cmd.Flags().Changed("priority") {
			req.Priority = &priority
		}

		client, err := api.NewClient()
		if err != nil {
			return err
		}

		task, err := client.UpdateTask(taskID, req)
		if err != nil {
			return err
		}

		return outputTask(task)
	},
}

var completeTaskCmd = &cobra.Command{
	Use:   "complete <id>",
	Short: "Mark a task as complete",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		task, err := client.MarkTaskComplete(args[0])
		if err != nil {
			return err
		}

		return outputTask(task)
	},
}

var deleteTaskCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		if err := client.DeleteTask(args[0]); err != nil {
			return err
		}

		if !IsJSONOutput() {
			fmt.Println("Task deleted successfully")
		}
		return nil
	},
}

var searchTasksCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search tasks",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := api.NewClient()
		if err != nil {
			return err
		}

		tasks, err := client.SearchTasks(args[0])
		if err != nil {
			return err
		}

		return outputTasks(tasks)
	},
}

var agendaCmd = &cobra.Command{
	Use:   "agenda",
	Short: "Get daily agenda including recurring tasks",
	Long:  "Fetches all tasks for a date including recurring tasks. Unlike 'list', this shows the full daily agenda.",
	RunE: func(cmd *cobra.Command, args []string) error {
		date, _ := cmd.Flags().GetString("date")

		if date == "" {
			return fmt.Errorf("--date flag is required")
		}

		client, err := api.NewClient()
		if err != nil {
			return err
		}

		tasks, err := client.GetTasksForDate(date)
		if err != nil {
			return err
		}

		return outputTasks(tasks)
	},
}

func init() {
	// list command flags
	listTasksCmd.Flags().String("date", "", "Date in YYYY-MM-DD format (required)")
	listTasksCmd.Flags().String("timezone", "", "Timezone (e.g., America/New_York)")

	// by-list command flags
	byListCmd.Flags().String("list-id", "", "List ID (required)")

	// agenda command flags
	agendaCmd.Flags().String("date", "", "Date in YYYY-MM-DD format (required)")

	// create command flags
	createTaskCmd.Flags().String("desc", "", "Task description (required)")
	createTaskCmd.Flags().String("date", "", "Date in YYYY-MM-DD format")
	createTaskCmd.Flags().String("start", "", "Start time")
	createTaskCmd.Flags().Int("estimated-time", 0, "Estimated time in seconds")
	createTaskCmd.Flags().String("list-id", "", "List ID")
	createTaskCmd.Flags().String("label", "", "Label ID")
	createTaskCmd.Flags().Int("priority", 0, "Priority (1-4)")

	// update command flags
	updateTaskCmd.Flags().String("desc", "", "Task description")
	updateTaskCmd.Flags().String("date", "", "Date in YYYY-MM-DD format")
	updateTaskCmd.Flags().String("start", "", "Start time")
	updateTaskCmd.Flags().Int("estimated-time", 0, "Estimated time in seconds")
	updateTaskCmd.Flags().Bool("complete", false, "Mark as complete")
	updateTaskCmd.Flags().String("list-id", "", "List ID")
	updateTaskCmd.Flags().String("label", "", "Label ID")
	updateTaskCmd.Flags().Int("priority", 0, "Priority (1-4)")

	tasksCmd.AddCommand(getTaskCmd)
	tasksCmd.AddCommand(listTasksCmd)
	tasksCmd.AddCommand(byListCmd)
	tasksCmd.AddCommand(braindumpCmd)
	tasksCmd.AddCommand(createTaskCmd)
	tasksCmd.AddCommand(updateTaskCmd)
	tasksCmd.AddCommand(completeTaskCmd)
	tasksCmd.AddCommand(deleteTaskCmd)
	tasksCmd.AddCommand(searchTasksCmd)
	tasksCmd.AddCommand(agendaCmd)
}

func outputTask(task *models.Task) error {
	if IsJSONOutput() {
		data, err := json.MarshalIndent(task, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	printTask(task)
	return nil
}

func outputTasks(tasks []models.Task) error {
	if IsJSONOutput() {
		data, err := json.MarshalIndent(tasks, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}

	for i, task := range tasks {
		if i > 0 {
			fmt.Println()
		}
		printTask(&task)
	}
	return nil
}

func printTask(task *models.Task) {
	status := "[ ]"
	if task.Complete {
		status = "[x]"
	}

	fmt.Printf("%s %s\n", status, task.Description)
	fmt.Printf("    ID: %s\n", task.ID)

	if dateStr := task.GetDateString(); dateStr != "" {
		fmt.Printf("    Date: %s\n", dateStr)
	}

	if startStr := task.GetStartString(); startStr != "" {
		fmt.Printf("    Start: %s\n", startStr)
	}

	if task.EstimatedTime != nil && *task.EstimatedTime > 0 {
		minutes := *task.EstimatedTime / 60
		if minutes > 0 {
			fmt.Printf("    Estimated: %d min\n", minutes)
		}
	}

	if task.Priority != nil {
		fmt.Printf("    Priority: %s\n", priorityString(*task.Priority))
	}

	if task.Label != nil {
		fmt.Printf("    Label: %s\n", *task.Label)
	}

	if task.ListID != nil {
		fmt.Printf("    List: %s\n", *task.ListID)
	}
}

func priorityString(p int) string {
	switch p {
	case 1:
		return "Low"
	case 2:
		return "Medium"
	case 3:
		return "High"
	case 4:
		return "Urgent"
	default:
		return fmt.Sprintf("%d", p)
	}
}
