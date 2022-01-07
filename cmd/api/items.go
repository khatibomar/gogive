package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/khatibomar/gogive/internal/data"
	"github.com/khatibomar/gogive/internal/validator"
)

func (app *application) createItemHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name       string   `json:"name"`
		Categories []string `json:"categories"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	v.Check(input.Name != "", "name", "must be provided")

	v.Check(input.Categories != nil, "categories", "must be provided")
	v.Check(len(input.Categories) >= 1, "categories", "must contain at least 1 genre")
	v.Check(len(input.Categories) <= 5, "categories", "must not contain more than 5 genres")
	v.Check(validator.Unique(input.Categories), "categories", "must not contain duplicate values")

	if !v.Valid() {
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
