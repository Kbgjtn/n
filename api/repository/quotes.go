package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/Kbgjtn/notethingness-api.git/api/model"
	"github.com/Kbgjtn/notethingness-api.git/types"
)

// QuotesRepository is a struct that defines the repository for quotes
type QuotesRepository struct {
	db *sql.DB
}

// New returns a new instance of QuotesRepository
func New(db *sql.DB) *QuotesRepository {
	return &QuotesRepository{db}
}

// Get returns a quote by id
func (r QuotesRepository) Get(
	ctx context.Context,
	args *model.QuoteURLParams,
) (model.Quote, error) {
	query := ` SELECT * FROM "quotes" WHERE "id" = $1 LIMIT 1 `
	row := r.db.QueryRowContext(ctx, query, args.ID)
	var quote model.Quote

	if err := row.Scan(
		&quote.ID,
		&quote.Content,
		&quote.AuthorID,
		&quote.CategoryID,
		&quote.CreatedAt,
		&quote.UpdatedAt,
	); err != nil {
		return quote, err
	}
	return quote, nil
}

// List returns a list of quotes with pagination
func (r QuotesRepository) List(
	ctx context.Context,
	pag *types.Pageable,
) (model.Quotes, error) {
	query := ` SELECT *, COUNT(*) OVER() AS total  
		FROM "quotes" 
		ORDER by "id"
		LIMIT $1 OFFSET $2 `
	rows, err := r.db.QueryContext(ctx, query, pag.Limit, pag.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var quotes []model.Quote

	for rows.Next() {
		var quote model.Quote
		if err := rows.Scan(
			&quote.ID,
			&quote.Content,
			&quote.AuthorID,
			&quote.CategoryID,
			&quote.CreatedAt,
			&quote.UpdatedAt,
			&pag.Total,
		); err != nil {
			return nil, err
		}
		quotes = append(quotes, quote)
	}
	return quotes, nil
}

// Create creates a new quote
func (r QuotesRepository) Create(
	ctx context.Context,
	payload *model.QuoteRequestPayload,
) (model.Quote, error) {
	now := time.Now().UTC()
	query := ` INSERT INTO "quotes" 
		("content", "author_id", "created_at", "updated_at", "category_id")
		VALUES ($1, $2, $3, $4, $5)
		RETURNING * `

	row := r.db.QueryRowContext(
		ctx, query,
		payload.Content, payload.AuthorID, now, now, payload.CategoryID,
	)

	var quote model.Quote

	err := row.Scan(
		&quote.ID,
		&quote.Content,
		&quote.AuthorID,
		&quote.CategoryID,
		&quote.CreatedAt,
		&quote.UpdatedAt,
	)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Constraint == "quotes_category_id_fkey" {
			fmt.Println(pqErr.Message)
			return quote, fmt.Errorf(
				"error: field \"category_id\": %d does not exist or is invalid",
				payload.CategoryID,
			)
		}

		if ok && pqErr.Constraint == "quotes_author_id_fkey" {
			fmt.Println(pqErr.Message)
			return quote, fmt.Errorf(
				"error: field \"author_id\": %d does not exist or is invalid",
				payload.AuthorID,
			)
		}

		return quote, err
	}

	return quote, nil
}

// Delete deletes a quote by id
func (r QuotesRepository) Delete(ctx context.Context, args *model.QuoteURLParams) error {
	query := ` DELETE FROM "quotes" WHERE "id" = $1 `
	_, err := r.db.ExecContext(ctx, query, args.ID)
	if err != nil {
		return err
	}
	return nil
}

// Update updates a quote by id
func (r QuotesRepository) Update(
	ctx context.Context,
	params *model.QuoteURLParams,
	payload *model.QuoteRequestPayload,
) (model.Quote, error) {
	query := ` UPDATE "quotes" 
		SET "content" = $1, "author_id" = $2, "updated_at" = $3, "category_id" = $4
		WHERE "id" = $5
		RETURNING * `

	row := r.db.QueryRowContext(
		ctx, query,
		payload.Content, payload.AuthorID, time.Now().UTC(), payload.CategoryID, params.ID,
	)
	var quote model.Quote
	err := row.Scan(
		&quote.ID,
		&quote.Content,
		&quote.AuthorID,
		&quote.CategoryID,
		&quote.CreatedAt,
		&quote.UpdatedAt,
	)
	if err != nil {
		return quote, fmt.Errorf("error: quote with \"id\" %d not found", params.ID)
	}
	return quote, nil
}
