package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/Kbgjtn/notethingness-api.git/api/model"
	"github.com/Kbgjtn/notethingness-api.git/types"
)

type CategoryRepository struct {
	store *sql.DB
}

func NewCategoryRepo(store *sql.DB) *CategoryRepository {
	return &CategoryRepository{store}
}

func (r CategoryRepository) List(
	ctx context.Context,
	args *types.Pageable,
) (model.Categories, error) {
	query := `SELECT *, COUNT(*) OVER() AS total FROM "categories" ORDER BY "id" LIMIT $1 OFFSET $2`
	rows, err := r.store.QueryContext(ctx, query, args.Limit, args.Offset)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var categories model.Categories

	for rows.Next() {
		var category model.Category

		err := rows.Scan(
			&category.ID,
			&category.Label,
			&args.Total,
		)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}

		args.Calc()
		categories = append(categories, category)
	}

	return categories, nil
}

func (r CategoryRepository) Get(
	ctx context.Context,
	args model.RequestURLParam,
) (model.Category, error) {
	var category model.Category

	query := `SELECT * FROM "categories" WHERE "id" = $1 LIMIT 1`

	row, err := r.store.QueryContext(ctx, query, args.ID)
	if err != nil {
		return category, err
	}

	if row.Next() {
		if err = row.Scan(&category.ID, &category.Label); err != nil {
			return category, err
		}
	}

	return category, nil
}

func (r CategoryRepository) Create(c context.Context, label string) (model.Category, error) {
	query := `INSERT INTO "categories" ("label") VALUES ($1) RETURNING *`
	row := r.store.QueryRowContext(c, query, label)

	var category model.Category

	if err := row.Scan(&category.ID, &category.Label); err != nil {
		pqErr, ok := err.(*pq.Error)

		if ok && pqErr.Constraint == "categories_label_key" {
			return model.Category{}, fmt.Errorf(
				"error: category with label %s already exists", label,
			)
		}

		return model.Category{}, err
	}

	return category, nil
}

func (r CategoryRepository) Delete(c context.Context, args model.RequestURLParam) error {
	query := `DELETE FROM "categories" WHERE "id" = $1`
	_, err := r.store.ExecContext(c, query, args.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r CategoryRepository) Update(
	c context.Context, args model.RequestURLParam, label string,
) (model.Category, error) {
	query := `UPDATE "categories" SET "label" = $1 WHERE "id" = $2 RETURNING *`

	row, err := r.store.QueryContext(c, query, label, args.ID)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Constraint == "categories_label_key" {
			return model.Category{}, fmt.Errorf(
				"error: category with label %s already exists", label,
			)
		}
		return model.Category{}, err
	}

	if row.Next() {
		var category model.Category
		if err := row.Scan(&category.ID, &category.Label); err != nil {
			return category, err
		}
		category.Label = label
		return category, nil
	}

	return model.Category{}, nil
}
