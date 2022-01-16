package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	middlewares := alice.New(app.recoverPanic)

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/items", app.listItemsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/items", app.createItemHandler)
	router.HandlerFunc(http.MethodGet, "/v1/items/:id", app.showItemHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/items/:id", app.updateItemHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/items/:id", app.deleteItemHandler)

	if app.config.limiter.enabled {
		rl, cleanup := NewRateLimiter()
		go cleanup(rl)
		return middlewares.Then(rl.RateLimit(app, router))
	}

	return middlewares.Then(router)
}
