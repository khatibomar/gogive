package data

import (
	"time"

	"github.com/khatibomar/gogive/internal/validator"
)

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
