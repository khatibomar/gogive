package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/items", app.listItemsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/items", app.createItemHandler)
	router.HandlerFunc(http.MethodGet, "/v1/items/:id", app.showItemHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/items/:id", app.updateItemHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/items/:id", app.deleteItemHandler)

	return router
}
