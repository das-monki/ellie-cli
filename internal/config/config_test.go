package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetAPIKey_EnvVar(t *testing.T) {
	// Setup
	originalKey := os.Getenv("ELLIE_API_KEY")
	originalKeyFile := os.Getenv("ELLIE_API_KEY_FILE")
	defer func() {
		os.Setenv("ELLIE_API_KEY", originalKey)
		os.Setenv("ELLIE_API_KEY_FILE", originalKeyFile)
	}()

	os.Setenv("ELLIE_API_KEY", "test-api-key-from-env")
	os.Unsetenv("ELLIE_API_KEY_FILE")

	// Test
	key, err := GetAPIKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != "test-api-key-from-env" {
		t.Errorf("expected 'test-api-key-from-env', got '%s'", key)
	}
}

func TestGetAPIKey_FileEnvVar(t *testing.T) {
	// Setup
	originalKey := os.Getenv("ELLIE_API_KEY")
	originalKeyFile := os.Getenv("ELLIE_API_KEY_FILE")
	defer func() {
		os.Setenv("ELLIE_API_KEY", originalKey)
		os.Setenv("ELLIE_API_KEY_FILE", originalKeyFile)
	}()

	// Create temp file with API key
	tmpDir := t.TempDir()
	keyFile := filepath.Join(tmpDir, "api-key")
	if err := os.WriteFile(keyFile, []byte("test-api-key-from-file\n"), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	os.Unsetenv("ELLIE_API_KEY")
	os.Setenv("ELLIE_API_KEY_FILE", keyFile)

	// Test
	key, err := GetAPIKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != "test-api-key-from-file" {
		t.Errorf("expected 'test-api-key-from-file', got '%s'", key)
	}
}

func TestGetAPIKey_Priority(t *testing.T) {
	// Setup - both env vars set, ELLIE_API_KEY should win
	originalKey := os.Getenv("ELLIE_API_KEY")
	originalKeyFile := os.Getenv("ELLIE_API_KEY_FILE")
	defer func() {
		os.Setenv("ELLIE_API_KEY", originalKey)
		os.Setenv("ELLIE_API_KEY_FILE", originalKeyFile)
	}()

	// Create temp file with API key
	tmpDir := t.TempDir()
	keyFile := filepath.Join(tmpDir, "api-key")
	if err := os.WriteFile(keyFile, []byte("file-key"), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	os.Setenv("ELLIE_API_KEY", "env-key")
	os.Setenv("ELLIE_API_KEY_FILE", keyFile)

	// Test
	key, err := GetAPIKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != "env-key" {
		t.Errorf("expected 'env-key' (priority 1), got '%s'", key)
	}
}

func TestGetAPIKey_NotSet(t *testing.T) {
	// Setup - no key configured
	originalKey := os.Getenv("ELLIE_API_KEY")
	originalKeyFile := os.Getenv("ELLIE_API_KEY_FILE")
	originalHome := os.Getenv("HOME")
	originalXDG := os.Getenv("XDG_CONFIG_HOME")

	defer func() {
		os.Setenv("ELLIE_API_KEY", originalKey)
		os.Setenv("ELLIE_API_KEY_FILE", originalKeyFile)
		os.Setenv("HOME", originalHome)
		if originalXDG != "" {
			os.Setenv("XDG_CONFIG_HOME", originalXDG)
		} else {
			os.Unsetenv("XDG_CONFIG_HOME")
		}
	}()

	// Use temp dir for config to work in sandboxed environments
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	os.Setenv("XDG_CONFIG_HOME", tmpDir)

	os.Unsetenv("ELLIE_API_KEY")
	os.Unsetenv("ELLIE_API_KEY_FILE")

	// Initialize viper with empty config
	if err := Init(); err != nil {
		t.Fatalf("init error: %v", err)
	}

	// Test
	_, err := GetAPIKey()
	if err == nil {
		t.Error("expected error when no API key is set")
	}
}
