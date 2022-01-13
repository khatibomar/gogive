package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/khatibomar/gogive/internal/validator"
	"github.com/lib/pq"
)

type ItemModel struct {
	DB *sql.DB
}

func (i ItemModel) Insert(item *Item) error {
	query := `
	INSERT INTO items (name, categories)
	VALUES ($1,$2)
	RETURNING id,created_at,version`

	args := []interface{}{item.Name, pq.Array(item.Categories)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return i.DB.QueryRowContext(ctx, query, args...).Scan(&item.ID, &item.CreatedAt, &item.Version)
}

func (i ItemModel) Get(id int64) (*Item, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT id , created_at , name , categories,version
		FROM items
		WHERE id=$1`

	var item Item

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := i.DB.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.CreatedAt,
		&item.Name,
		pq.Array(&item.Categories),
		&item.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &item, nil
}

func (i ItemModel) Update(item *Item) error {
	query := `
        UPDATE items
        SET name = $1, categories = $2, version = version + 1
        WHERE id = $3 AND version = $4
        RETURNING version`

	args := []interface{}{
		item.Name,
		pq.Array(item.Categories),
		item.ID,
		item.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := i.DB.QueryRowContext(ctx, query, args...).Scan(&item.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (i ItemModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM items
	where id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := i.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (i ItemModel) GetAll(name string, categories []string, filters Filters) ([]*Item, error) {
	// http://rachbelaid.com/postgres-full-text-search-is-good-enough
	query := fmt.Sprintf(`
	SELECT id, created_at, name, categories, version
        FROM items
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
        AND (categories @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := i.DB.QueryContext(ctx, query, name, pq.Array(categories))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*Item{}

	for rows.Next() {
		var item Item

		err := rows.Scan(
			&item.ID,
			&item.CreatedAt,
			&item.Name,
			pq.Array(&item.Categories),
			&item.Version,
		)

		if err != nil {
			return nil, err
		}

		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

type Item struct {
	ID         int64     `json:"id"`
	CreatedAt  time.Time `json:"-"`
	Name       string    `json:"name"`
	Categories []string  `json:"categories,omitempty"`
	Version    int32     `json:"version"`
}

func ValidateItem(v *validator.Validator, item *Item) {
	v.Check(item.Name != "", "name", "must be provided")
	v.Check(len(item.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(item.Categories != nil, "categories", "must be provided")
	v.Check(len(item.Categories) >= 1, "categories", "must contain at least 1 categorie")
	v.Check(len(item.Categories) <= 5, "categories", "must not contain more than 5 categories")
	v.Check(validator.Unique(item.Categories), "categories", "must not contain duplicate values")
}
