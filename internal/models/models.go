package models

import "encoding/json"

// Task represents a task in the ELLIE planner
type Task struct {
	ID            string          `json:"id"`
	Description   string          `json:"description"`
	Date          json.RawMessage `json:"date,omitempty"`
	Start         json.RawMessage `json:"start,omitempty"`
	DueDate       json.RawMessage `json:"due_date,omitempty"`
	EstimatedTime *int            `json:"estimated_time,omitempty"`
	ActualTime    *int            `json:"actual_time,omitempty"`
	Complete      bool            `json:"complete"`
	CompletedAt   json.RawMessage `json:"completed_at,omitempty"`
	ListID        *string         `json:"listId,omitempty"`
	Label         *string         `json:"label,omitempty"`
	Priority      *int            `json:"priority,omitempty"`
	RecurringID   *string         `json:"recurring_id,omitempty"`
	Recurring     bool            `json:"recurring"`
	CreatedAt     json.RawMessage `json:"created_at,omitempty"`
}

// GetDateString attempts to extract a date string from the Date field
func (t *Task) GetDateString() string {
	if t.Date == nil || string(t.Date) == "null" {
		return ""
	}
	var dateStr string
	if json.Unmarshal(t.Date, &dateStr) == nil {
		return dateStr
	}
	return ""
}

// GetStartString attempts to extract a start time string
func (t *Task) GetStartString() string {
	if t.Start == nil || string(t.Start) == "null" {
		return ""
	}
	var startStr string
	if json.Unmarshal(t.Start, &startStr) == nil {
		return startStr
	}
	return ""
}

// Subtask represents a subtask within a task
type Subtask struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	SortOrder   *int   `json:"sortOrder,omitempty"`
}

// Label represents a label for categorizing tasks
type Label struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// List represents a task list
type List struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Icon        string  `json:"icon,omitempty"`
	AutoLabelID *string `json:"auto_label_id,omitempty"`
}

// User represents the current user
type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// APIUsage represents API usage statistics
type APIUsage struct {
	Today     APIUsageToday     `json:"today"`
	RateLimit APIUsageRateLimit `json:"rateLimit"`
	ResetAt   string            `json:"resetAt"`
}

// APIUsageToday represents today's API usage
type APIUsageToday struct {
	Date      string `json:"date"`
	Used      int    `json:"used"`
	Remaining int    `json:"remaining"`
	Limit     int    `json:"limit"`
}

// APIUsageRateLimit represents rate limit info
type APIUsageRateLimit struct {
	RequestsPerMinute int `json:"requestsPerMinute"`
	WindowMs          int `json:"windowMs"`
}

// CreateTaskRequest represents the request body for creating a task
type CreateTaskRequest struct {
	Description   string  `json:"description"`
	Date          *string `json:"date,omitempty"`
	Start         *string `json:"start,omitempty"`
	EstimatedTime *int    `json:"estimated_time,omitempty"`
	ListID        *string `json:"listId,omitempty"`
	Label         *string `json:"label,omitempty"`
	Priority      *int    `json:"priority,omitempty"`
}

// UpdateTaskRequest represents the request body for updating a task
type UpdateTaskRequest struct {
	Description   *string `json:"description,omitempty"`
	Date          *string `json:"date,omitempty"`
	Start         *string `json:"start,omitempty"`
	EstimatedTime *int    `json:"estimated_time,omitempty"`
	Complete      *bool   `json:"complete,omitempty"`
	ListID        *string `json:"listId,omitempty"`
	Label         *string `json:"label,omitempty"`
	Priority      *int    `json:"priority,omitempty"`
}

// CreateLabelRequest represents the request body for creating a label
type CreateLabelRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// SearchRequest represents the request body for searching tasks
type SearchRequest struct {
	Query string `json:"query"`
}

// DeleteTaskRequest represents the request body for deleting a task
type DeleteTaskRequest struct {
	TaskID string `json:"taskId"`
}

// APIError represents an error response from the API
type APIError struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
