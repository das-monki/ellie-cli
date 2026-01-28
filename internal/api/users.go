package api

import (
	"encoding/json"
	"fmt"

	"github.com/goldie/ellie-cli/internal/models"
)

// GetCurrentUser retrieves the current user
func (c *Client) GetCurrentUser() (*models.User, error) {
	resp, err := c.Get("/v1/users/me")
	if err != nil {
		return nil, err
	}

	var user models.User
	if err := json.Unmarshal(resp, &user); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &user, nil
}

// GetAPIUsage retrieves API usage statistics
func (c *Client) GetAPIUsage() (*models.APIUsage, error) {
	resp, err := c.Get("/v1/users/apiUsage")
	if err != nil {
		return nil, err
	}

	var usage models.APIUsage
	if err := json.Unmarshal(resp, &usage); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &usage, nil
}
