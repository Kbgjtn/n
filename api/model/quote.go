package model

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Kbgjtn/notethingness-api.git/types"
)

// Quote is a quote entity model
type Quote struct {
	ID         int64     `json:"id"          example:"1"`
	Content    string    `json:"content"     example:"I am a quote"`
	AuthorID   int64     `json:"author_id"   example:"1"`
	CategoryID int64     `json:"category_id" example:"1"`
	CreatedAt  time.Time `json:"created_at"  example:"2021-01-01T00:00:00Z"`
	UpdatedAt  time.Time `json:"updated_at"  example:"2021-01-01T00:00:00Z"`
}

// CreateResponseDto creates a response dto for a single quote
func (q Quote) CreateResponseDto() types.JSONResult {
	return types.JSONResult{
		Data:    q,
		Code:    200,
		Message: "success",
	}
}

// Quotes is a collection of Quote
type Quotes []Quote

// Len is the number of elements in the collection.
func (q Quotes) Len() int {
	return len(q)
}

// Less reports whether the element with
func (q Quotes) CreateResponseDto(pag *types.Pageable) types.JSONResultWithPaginate {
	if pag.Total > 0 {
		return types.JSONResultWithPaginate{
			Message:  "success",
			Code:     200,
			Data:     q,
			Length:   len(q),
			Paginate: pag,
		}
	}

	return types.JSONResultWithPaginate{
		Code:    200,
		Message: "success",
		Data:    q,
		Length:  len(q),
	}
}

// CreateQuoteRequest is the request payload for creating a quote
type CreateQuoteRequest struct {
	Content    string `json:"content"     example:"I am a quote"`
	AuthorID   int64  `json:"author_id"   example:"1"`
	CategoryID int64  `json:"category_id" example:"1"`
}

type QuoteRequestPayload struct {
	Content    string `json:"content"     example:"I am a quote"`
	AuthorID   int64  `json:"author_id"   example:"1"`
	CategoryID int64  `json:"category_id" example:"1"`
}

// Validate validates the CreateQuoteRequest payload
func (req *QuoteRequestPayload) Validate() (bool, error) {
	if req == nil {
		return false, errors.New("error: payload is required")
	}

	if req.Content == "" {
		return false, errors.New("error: field \"content\" is required and must be a string")
	}

	if req.AuthorID <= 0 {
		return false, errors.New(
			"error: field \"author_id\" is required and must be a positive number",
		)
	}

	if req.CategoryID <= 0 {
		return false, errors.New(
			"error: field \"category_id\" is required and must be a positive number",
		)
	}

	return true, nil
}

// Parse parses the CreateQuoteRequest payload
func (req *QuoteRequestPayload) Parse() Quote {
	return Quote{
		Content:    req.Content,
		AuthorID:   req.AuthorID,
		CategoryID: req.CategoryID,
	}
}

// DataQuote is the response dto for a single quote
type DataQuote struct {
	Data Quote `json:"data" example:"{...}"`
}

// ListDataQuotes is the response dto for a list of quotes
type ListDataQuotes struct {
	Data     Quotes          `json:"data"`
	Length   int             `json:"length"`
	Paginate *types.Pageable `json:"paginate"`
}

// CreateQuoteResponse is the response dto for creating a quote
type CreateQuoteResponse *DataQuote

// ListQuoteResponse is the response dto for a list of quotes
type ListQuoteResponse ListDataQuotes

// ListQuoteRequest is the request payload for a list of quotes
type ListQuoteRequest types.Pageable

// QuoteURLParams is the url parameters for a single quote
type QuoteURLParams struct {
	ID int64 `json:"id" example:"1" validate:"required"`
}

// Parse parses the id url parameter
func (req *QuoteURLParams) Parse(value string) error {
	p, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return fmt.Errorf("param \"id\" is required and must be a number")
	}

	req.ID = p
	return nil
}

// GetQuoteResponse is the response dto for a single quote
type GetQuoteResponse DataQuote

// New creates a new Quote model
func (q *Quote) New(content string, authorID int64) Quote {
	return Quote{
		Content:  content,
		AuthorID: authorID,
	}
}

// Validate validates the Quote model
func (q *Quote) Validate() bool {
	if q == nil {
		return false
	}

	if q.Content == "" {
		return false
	}

	if q.AuthorID <= 0 {
		return false
	}

	return true
}

// CreateListQuoteReponse creates a response dto for a list of quotes
func CreateListQuoteReponse(quotes Quotes) ListQuoteResponse {
	return ListQuoteResponse{
		Data: quotes,
	}
}
