package model

import (
	"errors"
)

// Category represents a category
type Category struct {
	ID    int    `json:"id"    example:"1"`
	Label string `json:"label" example:"Category 1"`
}

type CategoryRequestPayload struct {
	Label string `json:"label" example:"Category 1"`
}

func (c *Category) Validate() error {
	if c.Label == "" {
		return errors.New("error: label is required")
	}
	return nil
}
