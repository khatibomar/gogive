package data

import (
	"database/sql"
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

	return i.DB.QueryRow(query, args...).Scan(&item.ID, &item.CreatedAt, &item.Version)
}

func (i ItemModel) Get(id int64) (*Item, error) {
	return nil, nil
}

func (i ItemModel) Update(item *Item) error {
	return nil
}

func (i ItemModel) Delete(id int64) error {
	return nil
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
