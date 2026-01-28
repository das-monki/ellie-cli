package api

import (
	"encoding/json"
	"fmt"

	"github.com/goldie/ellie-cli/internal/models"
)

// GetLists retrieves all lists
func (c *Client) GetLists() ([]models.List, error) {
	resp, err := c.Get("/v1/lists/getLists")
	if err != nil {
		return nil, err
	}

	var lists []models.List
	if err := json.Unmarshal(resp, &lists); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return lists, nil
}
