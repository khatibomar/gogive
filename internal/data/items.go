package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/khatibomar/gogive/internal/validator"
)

type ItemModel struct {
	DB *sql.DB
}

func (i ItemModel) Insert(item *Item) error {
	query := `
	INSERT INTO items (name,quantity, pcode, user_id, category_id, image_url)
	VALUES ($1,$2,$3,$4,(SELECT id FROM categories where category_name=LOWER($5)),$6)
	RETURNING id,created_at,version`

	args := []interface{}{item.Name, item.Quantity, item.Pcode, item.CreatedBy, item.Category, item.ImageURL}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return i.DB.QueryRowContext(ctx, query, args...).Scan(&item.ID, &item.CreatedAt, &item.Version)
}

func (i ItemModel) Get(id int64) (*Item, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
		SELECT items.id , created_at , name , quantity, category_name, user_id,version
		FROM items LEFT JOIN categories on items.category_id=categories.id
		WHERE items.id=$1`

	var item Item

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := i.DB.QueryRowContext(ctx, query, id).Scan(
		&item.ID,
		&item.CreatedAt,
		&item.Name,
		&item.Quantity,
		&item.Category,
		&item.CreatedBy,
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
		SET name=$1, pcode=$2, category_id=(SELECT id FROM categories WHERE category_name=$3), image_url=$4 , quantity=$5 , version=version+1
        WHERE id = $6 AND version = $7
        RETURNING version`

	args := []interface{}{
		item.Name,
		item.Pcode,
		item.Category,
		item.ImageURL,
		item.Quantity,
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

func (i ItemModel) GetAll(name string, category_name string, filters Filters, cursor *Cursor) ([]*Item, Metadata, error) {
	// http://rachbelaid.com/postgres-full-text-search-is-good-enough
	// NOTE(khatibomar): I strongly feel that my implementation is not safe at all and buggy
	// first I am using interface{} also there may be many edge cases when api
	// grows more and more , so maybe drop performance and use slow offset pagination instead
	// performance is not critical here I know because it's just for fun and no large data
	// will be used :)
	query := `
		SELECT count(*) OVER(),items.id, created_at, name, quantity,category_name as category,pcode, version
		FROM items LEFT JOIN categories on items.category_id=categories.id
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (category_name=$2 OR $2 = '')`
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
	rows, err := stmt.QueryContext(ctx, name, category_name, filters.PageSize)
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
			&item.Quantity,
			&item.Category,
			&item.Pcode,
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
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	CreatedBy int64     `json:"created_by"`
	Quantity  int       `json:"quantity"`
	Pcode     string    `json:"pcode"`
	ImageURL  string    `json:"image_url,omitempty"`
	Version   int32     `json:"version"`
}

func ValidateItem(v *validator.Validator, item *Item) {
	v.Check(item.Name != "", "name", "must be provided")
	v.Check(len(item.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(item.Pcode != "", "location", "must be provided")

	v.Check(item.Quantity > 0, "quantity", "must be provided")

	categories := []string{"vehicules", "mobile phones and accessories", "electronics", "fashion", "pets", "kids and babies", "services", "other"}
	v.Check(item.Category != "", "category", "must be provided")
	v.Check(validator.In(item.Category, categories...), "category", "not in allowed set of categories")
}
