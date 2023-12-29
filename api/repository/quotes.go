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

type QuotesRepository struct {
	store *sql.DB
}

func NewQuoteRepo(s *sql.DB) *QuotesRepository {
	return &QuotesRepository{s}
}

func (db *QuotesRepository) Get(c context.Context, arg interface{}) (interface{}, error) {
	query := ` SELECT * FROM "quotes" WHERE "id" = $1 LIMIT 1 `
	row := db.store.QueryRowContext(c, query, arg.(model.QuoteURLParams).ID)
	var quote model.Quote

	if err := row.Scan(
		&quote.ID,
		&quote.Content,
		&quote.AuthorID,
		&quote.CategoryID,
		&quote.CreatedAt,
		&quote.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return quote, nil
}

func (db *QuotesRepository) List(
	c context.Context,
	arg types.Pageable,
) (interface{}, interface{}, error) {
	query := ` SELECT *, COUNT(*) OVER() AS total  FROM "quotes" ORDER by "id" LIMIT $1 OFFSET $2 `
	rows, err := db.store.QueryContext(c, query, arg.Limit, arg.Offset)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var quotes model.Quotes

	for rows.Next() {
		var quote model.Quote
		if err := rows.Scan(
			&quote.ID,
			&quote.Content,
			&quote.AuthorID,
			&quote.CategoryID,
			&quote.CreatedAt,
			&quote.UpdatedAt,
			&arg.Total,
		); err != nil {
			return nil, nil, err
		}
		quotes = append(quotes, quote)
	}

	return quotes, arg, nil
}

func (db *QuotesRepository) Create(c context.Context, arg interface{}) (interface{}, error) {
	now := time.Now().UTC()
	query := ` INSERT INTO "quotes" 
		("content", "author_id", "created_at", "updated_at", "category_id")
		VALUES ($1, $2, $3, $4, $5)
		RETURNING * `

	payload := arg.(model.QuoteRequestPayload)

	row := db.store.QueryRowContext(
		c, query,
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

func (db *QuotesRepository) Update(
	c context.Context,
	arg interface{},
	pay interface{},
) (interface{}, error) {
	param := arg.(model.QuoteURLParams)
	payload := pay.(model.QuoteRequestPayload)
	query := ` UPDATE "quotes" 
		SET "content" = $1, "author_id" = $2, "updated_at" = $3, "category_id" = $4
		WHERE "id" = $5
		RETURNING * `

	row := db.store.QueryRowContext(
		c, query,
		payload.Content, payload.AuthorID, time.Now().UTC(), payload.CategoryID, param.ID,
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
		return quote, fmt.Errorf("error: quote with \"id\" %d not found", param.ID)
	}
	return quote, nil
}

func (db *QuotesRepository) Delete(c context.Context, arg interface{}) error {
	query := ` DELETE FROM "quotes" WHERE "id" = $1 `
	_, err := db.store.ExecContext(c, query, arg.(model.QuoteURLParams).ID)
	if err != nil {
		return err
	}

	return nil
}
