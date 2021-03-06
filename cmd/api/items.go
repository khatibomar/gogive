package main

import (
	"errors"
	"fmt"
	"net/http"

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
		Name     string `json:"name"`
		Quantity int    `json:"quantity"`
		Category string `json:"category"`
		Pcode    string `json:"pcode"`
		ImageURL string `json:"image_url,omitempty"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	item := &data.Item{
		Name:      input.Name,
		Category:  input.Category,
		Quantity:  input.Quantity,
		Pcode:     input.Pcode,
		ImageURL:  input.ImageURL,
		CreatedBy: app.contextGetUser(r).ID,
	}

	v := validator.New()

	// Here we are passing validator as dependency instead of decleare it
	// in the function itself because we may re-use this validator later
	if data.ValidateItem(v, item); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Items.Insert(item)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	// When sending a HTTP response, we want to include a Location header to let the
	// client know which URL they can find the newly-created resource at. We make an
	// empty http.Header map and then use the Set() method to add a new Location header,
	// interpolating the system-generated ID for our new movie in the URL.
	headers.Set("Location", fmt.Sprintf("/v1/items/%d", item.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"item": item}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	item, err := app.models.Items.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"item": item}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	item, err := app.models.Items.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// here using a pointer to a string because in case of partially
	// we need to differentiate between the user send zero values
	// or provided nothing at all
	var input struct {
		Name     *string `json:"name"`
		Quantity *int    `json:"quantity"`
		Category *string `json:"category"`
		Pcode    *string `json:"pcode"`
		ImageURL *string `json:"image_url"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		item.Name = *input.Name
	}

	if input.Category != nil {
		item.Category = *input.Category
	}

	if input.Pcode != nil {
		item.Pcode = *input.Pcode
	}

	if input.ImageURL != nil {
		item.ImageURL = *input.ImageURL
	}

	if input.Quantity != nil {
		item.Quantity = *input.Quantity
	}

	v := validator.New()

	if data.ValidateItem(v, item); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
	}

	err = app.models.Items.Update(item)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"item": item}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Items.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "item successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listItemsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string
		Category string
		Filters  data.Filters
		Cursor   *data.Cursor
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Category = app.readString(qs, "category", "")

	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "-id", "-name"}

	if qs.Has("last_sort_val") || qs.Has("last_id") {
		input.Cursor = &data.Cursor{}
		input.Cursor.LastSortVal = app.readString(qs, "last_sort_val", "")
		input.Cursor.LastID = app.readInt64(qs, "last_id", 0, v)
	}

	data.ValidateFilters(v, input.Filters)
	if input.Cursor != nil {
		data.ValidateCursors(v, *input.Cursor)
	}
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	items, metadata, err := app.models.Items.GetAll(input.Name, input.Category, input.Filters, input.Cursor)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"items": items, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
