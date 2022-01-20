package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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

func (i ItemModel) GetAll(name string, categories []string, filters Filters, cursor *Cursor) ([]*Item, Metadata, error) {
	// http://rachbelaid.com/postgres-full-text-search-is-good-enough
	// NOTE(khatibomar): I strongly feel that my implementation is not safe at all and buggy
	// first I am using interface{} also there may be many edge cases when api
	// grows more and more , so maybe drop performance and use slow offset pagination instead
	// performance is not critical here I know because it's just for fun and no large data
	// will be used :)
	query := `
		SELECT count(*) OVER(),id, created_at, name, categories, version
		FROM items
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (categories @> $2 OR $2 = '{}')
	`
	if cursor != nil {
		if strings.HasPrefix(filters.Sort, "-") {
			if filters.sortColumn() == "id" {
				if cursor.LastID != 0 {
					query += fmt.Sprintf(`AND id < %d `, cursor.LastID)
				}
			} else {
				if cursor.LastSortVal != "" {
					query += fmt.Sprintf(`AND %s < '%s' `, filters.sortColumn(), cursor.LastSortVal)
				}
			}
		} else {
			if filters.sortColumn() == "id" {
				query += fmt.Sprintf(`AND id > %d `, cursor.LastID)
			} else {
				query += fmt.Sprintf(`AND row(%s,id) > ('%s',%d) `, filters.sortColumn(), cursor.LastSortVal, cursor.LastID)
			}
		}
	}

	query += fmt.Sprintf(`
		ORDER BY %s %s, id ASC
		fetch first $3 rows only`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// using statement instead to help from sql injection
	stmt, err := i.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, Metadata{}, err
	}
	rows, err := stmt.QueryContext(ctx, name, pq.Array(categories), filters.PageSize)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0

	items := []*Item{}

	for rows.Next() {
		var item Item

		err := rows.Scan(
			&totalRecords,
			&item.ID,
			&item.CreatedAt,
			&item.Name,
			pq.Array(&item.Categories),
			&item.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	if len(items) > 0 {
		if cursor == nil {
			cursor = &Cursor{}
		}
		cursor.LastID = items[len(items)-1].ID
		if filters.sortColumn() == "id" {
			cursor.LastSortVal = items[len(items)-1].ID
		} else {
			cursor.LastSortVal = items[len(items)-1].Name
		}
	}

	metadata := newMetadata(totalRecords, filters.PageSize, cursor)

	return items, metadata, nil
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
