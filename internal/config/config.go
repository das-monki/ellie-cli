package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	configFileName = "config"
	configFileType = "yaml"
	configDirName  = "ellie"
)

// Config holds the application configuration
type Config struct {
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
}

// DefaultBaseURL is the default API base URL
const DefaultBaseURL = "https://api.ellieplanner.com"

// GetConfigDir returns the configuration directory path
func GetConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config directory: %w", err)
	}
	return filepath.Join(configDir, configDirName), nil
}

// Init initializes the configuration
func Init() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.AddConfigPath(configDir)

	// Set defaults
	viper.SetDefault("base_url", DefaultBaseURL)

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	return nil
}

// GetAPIKey returns the API key using priority:
// 1. ELLIE_API_KEY environment variable
// 2. ELLIE_API_KEY_FILE environment variable (reads from file)
// 3. Config file
func GetAPIKey() (string, error) {
	// Priority 1: Direct environment variable
	if apiKey := os.Getenv("ELLIE_API_KEY"); apiKey != "" {
		return strings.TrimSpace(apiKey), nil
	}

	// Priority 2: File-based environment variable (for agenix/secrets management)
	if keyFile := os.Getenv("ELLIE_API_KEY_FILE"); keyFile != "" {
		content, err := os.ReadFile(keyFile)
		if err != nil {
			return "", fmt.Errorf("failed to read API key file: %w", err)
		}
		return strings.TrimSpace(string(content)), nil
	}

	// Priority 3: Config file
	if apiKey := viper.GetString("api_key"); apiKey != "" {
		return apiKey, nil
	}

	return "", fmt.Errorf("API key not configured. Set ELLIE_API_KEY, ELLIE_API_KEY_FILE, or run 'ellie config set-api-key <key>'")
}

// SetAPIKey saves the API key to the config file
func SetAPIKey(apiKey string) error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	viper.Set("api_key", apiKey)

	configPath := filepath.Join(configDir, configFileName+"."+configFileType)
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetBaseURL returns the API base URL
func GetBaseURL() string {
	if baseURL := os.Getenv("ELLIE_BASE_URL"); baseURL != "" {
		return baseURL
	}
	return viper.GetString("base_url")
}

// SetBaseURL saves the base URL to the config file
func SetBaseURL(baseURL string) error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	viper.Set("base_url", baseURL)

	configPath := filepath.Join(configDir, configFileName+"."+configFileType)
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfig returns the current configuration
func GetConfig() *Config {
	apiKey, _ := GetAPIKey()
	return &Config{
		APIKey:  apiKey,
		BaseURL: GetBaseURL(),
	}
}
