package model

import (
	"errors"
	"strconv"

	"github.com/Kbgjtn/notethingness-api.git/types"
	"github.com/Kbgjtn/notethingness-api.git/util"
)

// Category represents a category
type Category struct {
	ID    int    `json:"id"    example:"1"`
	Label string `json:"label" example:"My Category"`
}

type CategoryRequestPayload struct {
	Label string `json:"label" example:"My Category"`
}

func (c CategoryRequestPayload) Validate() error {
	if c.Label == "" {
		return errors.New("error: label is required")
	}

	if len(c.Label) > 255 {
		return errors.New("error: label must be less than 255 characters")
	}

	// label should be unique and only allow alphabetical characters
	if _, err := strconv.Atoi(c.Label); err == nil {
		return errors.New("error: label must be a string")
	}

	// check if label is contained symbols
	if !util.ContainsOnlyAlphabet(c.Label) {
		return errors.New("error: label must only contain alphabetical characters")
	}

	return nil
}

func (c *Category) ToJSON(code int, message string) types.JSONResult {
	return types.JSONResult{
		Data:    c,
		Code:    code,
		Message: message,
	}
}

func (c *Category) Validate() error {
	if c.Label == "" {
		return errors.New("error: label is required")
	}
	return nil
}

type RequestURLParam struct {
	ID int `json:"id" example:"1"`
}

func ParseParams(v string) (RequestURLParam, error) {
	var r RequestURLParam

	if v == "" {
		return r, errors.New("error: id is required")
	}

	parsed, err := strconv.Atoi(v)
	if err != nil {
		return r, errors.New("error: id must be a number greater than 0")
	}

	if parsed <= 0 {
		return r, errors.New("error: id is required and must be a number greater than 0")
	}
	r.ID = parsed
	return r, nil
}

func (c *Category) toJSON() types.JSONResult {
	return types.JSONResult{
		Data:    c,
		Code:    200,
		Message: "success",
	}
}

type Categories []Category

func (c Categories) ToJSON(pag types.Pageable) types.JSONResultWithPaginate {
	pag.Calc()
	return types.JSONResultWithPaginate{
		Data:     c,
		Code:     200,
		Message:  "success",
		Paginate: &pag,
		Length:   len(c),
	}
}
