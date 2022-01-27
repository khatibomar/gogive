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

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)

	// If a client sends the same PUT /v1/users/activated request multiple times,
	// the first will succeed (assuming the token is valid) and then
	// any subsequent requests will result in an error being sent to the client
	// (because the token has been used and deleted from the database).
	// But the important thing is that nothing in our application state
	// (i.e. database) changes after that first request.
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	if app.config.limiter.enabled {
		rl, cleanup := NewRateLimiter()
		go cleanup(rl)
		return middlewares.Then(rl.RateLimit(app, router))
	}

	return middlewares.Then(router)
}
