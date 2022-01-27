package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Items  ItemModel
	Users  UserModel
	Tokens TokenModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Items:  ItemModel{DB: db},
		Tokens: TokenModel{DB: db},
		Users:  UserModel{DB: db},
	}
}
