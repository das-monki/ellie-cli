package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/goldie/ellie-cli/internal/config"
	"github.com/goldie/ellie-cli/internal/models"
)

func setupTestServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func setupTestClient(t *testing.T, serverURL string) *Client {
	// Set env vars for the test
	originalKey := os.Getenv("ELLIE_API_KEY")
	originalURL := os.Getenv("ELLIE_BASE_URL")
	originalHome := os.Getenv("HOME")
	originalXDG := os.Getenv("XDG_CONFIG_HOME")

	t.Cleanup(func() {
		os.Setenv("ELLIE_API_KEY", originalKey)
		os.Setenv("ELLIE_BASE_URL", originalURL)
		os.Setenv("HOME", originalHome)
		if originalXDG != "" {
			os.Setenv("XDG_CONFIG_HOME", originalXDG)
		} else {
			os.Unsetenv("XDG_CONFIG_HOME")
		}
	})

	// Use temp dir for config to work in sandboxed environments
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	os.Setenv("ELLIE_API_KEY", "test-api-key")
	os.Setenv("ELLIE_BASE_URL", serverURL)

	if err := config.Init(); err != nil {
		t.Fatalf("failed to init config: %v", err)
	}

	client, err := NewClient()
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return client
}

func TestClient_AuthHeader(t *testing.T) {
	var receivedAuth string
	server := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		receivedAuth = r.Header.Get("x-api-key")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	})
	defer server.Close()

	client := setupTestClient(t, server.URL)
	_, _ = client.Get("/test")

	if receivedAuth != "test-api-key" {
		t.Errorf("expected 'test-api-key', got '%s'", receivedAuth)
	}
}

func TestClient_GetUser(t *testing.T) {
	server := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/users/me" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}

		response := models.User{
			ID:    "user-123",
			Name:  "Test User",
			Email: "test@example.com",
		}
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	client := setupTestClient(t, server.URL)
	user, err := client.GetCurrentUser()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.ID != "user-123" {
		t.Errorf("expected ID 'user-123', got '%s'", user.ID)
	}
	if user.Name != "Test User" {
		t.Errorf("expected Name 'Test User', got '%s'", user.Name)
	}
}

func TestClient_GetTask(t *testing.T) {
	server := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/tasks/getTask" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("taskId") != "task-456" {
			t.Errorf("expected taskId 'task-456', got '%s'", r.URL.Query().Get("taskId"))
		}

		response := models.Task{
			ID:          "task-456",
			Description: "Test task",
			Complete:    false,
		}
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	client := setupTestClient(t, server.URL)
	task, err := client.GetTask("task-456")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task.ID != "task-456" {
		t.Errorf("expected ID 'task-456', got '%s'", task.ID)
	}
	if task.Description != "Test task" {
		t.Errorf("expected Description 'Test task', got '%s'", task.Description)
	}
}

func TestClient_CreateTask(t *testing.T) {
	server := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/tasks/createTask" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var req models.CreateTaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request: %v", err)
		}
		if req.Description != "New task" {
			t.Errorf("expected Description 'New task', got '%s'", req.Description)
		}

		response := models.Task{
			ID:          "new-task-789",
			Description: "New task",
			Complete:    false,
		}
		json.NewEncoder(w).Encode(response)
	})
	defer server.Close()

	client := setupTestClient(t, server.URL)
	task, err := client.CreateTask(&models.CreateTaskRequest{
		Description: "New task",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task.ID != "new-task-789" {
		t.Errorf("expected ID 'new-task-789', got '%s'", task.ID)
	}
}

func TestClient_APIError(t *testing.T) {
	server := setupTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Forbidden"})
	})
	defer server.Close()

	client := setupTestClient(t, server.URL)
	_, err := client.GetCurrentUser()
	if err == nil {
		t.Error("expected error for 401 response")
	}
}
