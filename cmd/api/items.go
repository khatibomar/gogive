package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/khatibomar/gogive/internal/data"
	"github.com/khatibomar/gogive/internal/validator"
)

func (app *application) createItemHandler(w http.ResponseWriter, r *http.Request) {
	// The problem with decoding directly into a Item struct
	// is that a client could provide the keys id and version in their JSON request,
	// and the corresponding values would be decoded without any error
	// into the ID and Version fields of the Item struct
	// even though we don’t want them to be.
	var input struct {
		Name       string   `json:"name"`
		Categories []string `json:"categories"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	item := &data.Item{
		Name:       input.Name,
		Categories: input.Categories,
	}

	v := validator.New()

	// Here we are passing validator as dependency instead of decleare it
	// in the function itself because we may re-use this validator later
	if data.ValidateItem(v, item); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	item := data.Item{
		ID:         id,
		CreatedAt:  time.Now(),
		Name:       "chohata",
		Categories: []string{"albisa", "weapon", "mom tools"},
		Version:    1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"item": item}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
