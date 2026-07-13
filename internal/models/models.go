package models

import (
	"encoding/json"
	"time"
)

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

// ParseTimestamp decodes the two shapes the API uses for an instant: an ISO 8601
// string, and a Firestore timestamp object ({"_seconds":…,"_nanoseconds":…}).
// Which one comes back depends on the endpoint -- createTask and updateTask echo
// strings, while getTask and byDate return the Firestore objects -- so both have
// to be handled or the timestamp silently reads as absent.
func ParseTimestamp(raw json.RawMessage) (time.Time, bool) {
	if len(raw) == 0 || string(raw) == "null" {
		return time.Time{}, false
	}

	var isoStr string
	if json.Unmarshal(raw, &isoStr) == nil {
		if isoStr == "" {
			return time.Time{}, false
		}
		if parsed, err := time.Parse(time.RFC3339, isoStr); err == nil {
			return parsed, true
		}
		if parsed, err := time.Parse("2006-01-02", isoStr); err == nil {
			return parsed, true
		}
		return time.Time{}, false
	}

	var ts struct {
		Seconds     *int64 `json:"_seconds"`
		Nanoseconds int64  `json:"_nanoseconds"`
	}
	if json.Unmarshal(raw, &ts) == nil && ts.Seconds != nil {
		return time.Unix(*ts.Seconds, ts.Nanoseconds).UTC(), true
	}

	return time.Time{}, false
}

// GetDate returns the task's date, if it has one.
func (t *Task) GetDate() (time.Time, bool) {
	return ParseTimestamp(t.Date)
}

// GetStart returns the task's start instant, if it has one.
func (t *Task) GetStart() (time.Time, bool) {
	return ParseTimestamp(t.Start)
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
