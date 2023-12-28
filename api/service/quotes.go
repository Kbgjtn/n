package service

import (
	"context"

	"github.com/Kbgjtn/notethingness-api.git/api/model"
	"github.com/Kbgjtn/notethingness-api.git/api/repository"
	"github.com/Kbgjtn/notethingness-api.git/types"
)

// QuotesService is a service for quotes
type QuotesService struct {
	repo *repository.QuotesRepository
}

// New creates a new QuotesService
func New(repo *repository.QuotesRepository) *QuotesService {
	return &QuotesService{repo}
}

// Get returns a quote
func (s QuotesService) Get(
	ctx context.Context,
	reqDto *model.QuoteURLParams,
) (model.GetQuoteResponse, error) {
	data, err := s.repo.Get(ctx, reqDto)
	return data.CreateResponseDto(), err
}

// List returns a list of quotes
func (s QuotesService) List(
	ctx context.Context,
	pag types.Pageable,
) (model.ListQuoteResponse, error) {
	data, err := s.repo.List(ctx, &pag)
	return data.CreateResponseDto(&pag), err
}

// Create creates a quote
func (s QuotesService) Create(
	ctx context.Context,
	quote *model.QuoteRequestPayload,
) (model.GetQuoteResponse, error) {
	data, err := s.repo.Create(ctx, quote)
	return data.CreateResponseDto(), err
}

// Delete deletes a quote
func (s QuotesService) Delete(ctx context.Context, reqDto *model.QuoteURLParams) error {
	return s.repo.Delete(ctx, reqDto)
}

// Update updates a quote
func (s QuotesService) Update(
	ctx context.Context,
	params *model.QuoteURLParams,
	payload *model.QuoteRequestPayload,
) (model.GetQuoteResponse, error) {
	data, err := s.repo.Update(ctx, params, payload)
	return data.CreateResponseDto(), err
}
