package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/goldie/ellie-cli/internal/api"
	"github.com/goldie/ellie-cli/internal/models"
	"github.com/spf13/cobra"
)

var (
	isoDateTimeRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T`)
	clockTimeRe   = regexp.MustCompile(`^(\d{1,2}):(\d{2})$`)
)

// resolveLocation picks the timezone that a bare HH:MM start time is written in.
// Empty means the machine's local zone, which is what someone saying "16:00" means.
func resolveLocation(timezone string) (*time.Location, error) {
	if timezone == "" {
		return time.Local, nil
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, fmt.Errorf("unknown timezone '%s': expected an IANA name like Europe/Zurich", timezone)
	}
	return loc, nil
}

// formatStartTime converts a HH:MM start time into the ISO datetime the API wants.
// The clock time is read in loc -- not UTC -- so `--start 16:00` means 16:00 where
// the user is, matching how the API already treats --date (stored as local midnight).
// The result is expressed in UTC, since only the instant travels over the wire.
// A start that is already an ISO datetime carries its own offset and is passed through.
func formatStartTime(start, date string, loc *time.Location) (string, error) {
	if start == "" {
		return "", nil
	}

	if isoDateTimeRe.MatchString(start) {
		return start, nil
	}

	if !clockTimeRe.MatchString(start) {
		return "", fmt.Errorf("invalid start time format '%s': expected HH:MM or ISO datetime", start)
	}

	// The clock time alone does not say which day it falls on.
	if date == "" {
		return "", fmt.Errorf("--date is required when using --start with HH:MM format")
	}

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return "", fmt.Errorf("invalid date format '%s': expected YYYY-MM-DD", date)
	}

	var hour, minute int
	if _, err := fmt.Sscanf(start, "%d:%d", &hour, &minute); err != nil {
		return "", fmt.Errorf("invalid start time format '%s': expected HH:MM or ISO datetime", start)
	}
	if hour > 23 || minute > 59 {
		return "", fmt.Errorf("invalid start time '%s': hour must be 0-23 and minute 0-59", start)
	}

	combined := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(),
		hour, minute, 0, 0, loc)

	return combined.UTC().Format(time.RFC3339), nil
}

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
		timezone, _ := cmd.Flags().GetString("timezone")
		estimatedTime, _ := cmd.Flags().GetInt("estimated-time")
		listID, _ := cmd.Flags().GetString("list-id")
		label, _ := cmd.Flags().GetString("label")
		priority, _ := cmd.Flags().GetInt("priority")

		if desc == "" {
			return fmt.Errorf("--desc flag is required")
		}

		loc, err := resolveLocation(timezone)
		if err != nil {
			return err
		}

		// Convert HH:MM time to ISO datetime if needed
		formattedStart, err := formatStartTime(start, date, loc)
		if err != nil {
			return err
		}

		req := &models.CreateTaskRequest{
			Description: desc,
		}

		if date != "" {
			req.Date = &date
		}
		if formattedStart != "" {
			req.Start = &formattedStart
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
		timezone, _ := cmd.Flags().GetString("timezone")
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
			loc, err := resolveLocation(timezone)
			if err != nil {
				return err
			}

			// Convert HH:MM time to ISO datetime if needed
			formattedStart, err := formatStartTime(start, date, loc)
			if err != nil {
				return err
			}
			req.Start = &formattedStart
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
	createTaskCmd.Flags().String("start", "", "Start time as HH:MM (local) or an ISO datetime")
	createTaskCmd.Flags().String("timezone", "", "Timezone the HH:MM start time is in (default: local)")
	createTaskCmd.Flags().Int("estimated-time", 0, "Estimated time in seconds")
	createTaskCmd.Flags().String("list-id", "", "List ID")
	createTaskCmd.Flags().String("label", "", "Label ID")
	createTaskCmd.Flags().Int("priority", 0, "Priority (1-4)")

	// update command flags
	updateTaskCmd.Flags().String("desc", "", "Task description")
	updateTaskCmd.Flags().String("date", "", "Date in YYYY-MM-DD format")
	updateTaskCmd.Flags().String("start", "", "Start time as HH:MM (local) or an ISO datetime")
	updateTaskCmd.Flags().String("timezone", "", "Timezone the HH:MM start time is in (default: local)")
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

	// Timestamps come back as instants; show them in local time, naming the zone so
	// there is no doubt about which clock the number belongs to.
	if date, ok := task.GetDate(); ok {
		fmt.Printf("    Date: %s\n", date.In(time.Local).Format("2006-01-02"))
	}

	if start, ok := task.GetStart(); ok {
		fmt.Printf("    Start: %s\n", start.In(time.Local).Format("15:04 MST"))
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
