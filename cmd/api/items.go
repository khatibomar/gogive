package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/khatibomar/gogive/internal/data"
)

func (app *application) createItemHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new item")
}

func (app *application) showItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	item := data.Item{
		ID:         id,
		CreatedAt:  time.Now(),
		Title:      "chohata",
		Categories: []string{"albisa", "weapon", "mom tools"},
		Version:    1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"item": item}, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
