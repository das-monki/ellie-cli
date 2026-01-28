package api

import (
	"encoding/json"
	"fmt"

	"github.com/goldie/ellie-cli/internal/models"
)

// GetLabels retrieves all labels
func (c *Client) GetLabels() ([]models.Label, error) {
	resp, err := c.Get("/v1/labels/getLabels")
	if err != nil {
		return nil, err
	}

	var labels []models.Label
	if err := json.Unmarshal(resp, &labels); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return labels, nil
}

// CreateLabel creates a new label
func (c *Client) CreateLabel(req *models.CreateLabelRequest) (*models.Label, error) {
	resp, err := c.Post("/v1/labels/createLabel", req)
	if err != nil {
		return nil, err
	}

	var label models.Label
	if err := json.Unmarshal(resp, &label); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &label, nil
}
