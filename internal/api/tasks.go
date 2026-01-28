package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/goldie/ellie-cli/internal/models"
)

// GetTask retrieves a task by ID
func (c *Client) GetTask(taskID string) (*models.Task, error) {
	path := fmt.Sprintf("/v1/tasks/getTask?taskId=%s", url.QueryEscape(taskID))
	resp, err := c.Get(path)
	if err != nil {
		return nil, err
	}

	var task models.Task
	if err := json.Unmarshal(resp, &task); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &task, nil
}

// GetTasksByDate retrieves tasks for a specific date
func (c *Client) GetTasksByDate(date, timeZone string) ([]models.Task, error) {
	path := fmt.Sprintf("/v1/tasks/byDate?date=%s", url.QueryEscape(date))
	if timeZone != "" {
		path += fmt.Sprintf("&timeZone=%s", url.QueryEscape(timeZone))
	}

	resp, err := c.Get(path)
	if err != nil {
		return nil, err
	}

	var tasks []models.Task
	if err := json.Unmarshal(resp, &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return tasks, nil
}

// GetTasksByList retrieves tasks for a specific list
func (c *Client) GetTasksByList(listID string) ([]models.Task, error) {
	path := fmt.Sprintf("/v1/tasks/byList?listId=%s", url.QueryEscape(listID))

	resp, err := c.Get(path)
	if err != nil {
		return nil, err
	}

	var tasks []models.Task
	if err := json.Unmarshal(resp, &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return tasks, nil
}

// GetBraindump retrieves unscheduled tasks
func (c *Client) GetBraindump() ([]models.Task, error) {
	resp, err := c.Get("/v1/tasks/getBraindump")
	if err != nil {
		return nil, err
	}

	var tasks []models.Task
	if err := json.Unmarshal(resp, &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return tasks, nil
}

// CreateTask creates a new task
func (c *Client) CreateTask(req *models.CreateTaskRequest) (*models.Task, error) {
	resp, err := c.Post("/v1/tasks/createTask", req)
	if err != nil {
		return nil, err
	}

	var task models.Task
	if err := json.Unmarshal(resp, &task); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &task, nil
}

// UpdateTask updates an existing task
func (c *Client) UpdateTask(taskID string, req *models.UpdateTaskRequest) (*models.Task, error) {
	path := fmt.Sprintf("/v1/tasks/updateTask/%s", url.PathEscape(taskID))
	resp, err := c.Post(path, req)
	if err != nil {
		return nil, err
	}

	var task models.Task
	if err := json.Unmarshal(resp, &task); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &task, nil
}

// MarkTaskComplete marks a task as complete
func (c *Client) MarkTaskComplete(taskID string) (*models.Task, error) {
	path := fmt.Sprintf("/v1/tasks/markTaskAsComplete?taskId=%s", url.QueryEscape(taskID))
	resp, err := c.Post(path, nil)
	if err != nil {
		return nil, err
	}

	var task models.Task
	if err := json.Unmarshal(resp, &task); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &task, nil
}

// DeleteTask deletes a task
func (c *Client) DeleteTask(taskID string) error {
	req := &models.DeleteTaskRequest{TaskID: taskID}
	_, err := c.Post("/v1/tasks/deleteTask", req)
	return err
}

// SearchTasks searches for tasks
func (c *Client) SearchTasks(query string) ([]models.Task, error) {
	req := &models.SearchRequest{Query: query}
	resp, err := c.Post("/v1/tasks/search", req)
	if err != nil {
		return nil, err
	}

	var tasks []models.Task
	if err := json.Unmarshal(resp, &tasks); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return tasks, nil
}
