package main

import (
	"context"
	"net/http"

	"github.com/khatibomar/gogive/internal/data"
)

type contextKey string

const userContextKey = "user"
const itemContextKey = "item"

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}

func (app *application) contextSetItem(r *http.Request, item *data.Item) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, item)
	return r.WithContext(ctx)
}

func (app *application) contextGetItem(r *http.Request) *data.Item {
	item, ok := r.Context().Value(itemContextKey).(*data.Item)
	if !ok {
		panic("missing item value in request context")
	}

	return item
}
