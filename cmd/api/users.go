package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/seminhnva/e-com/internal/store"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badReqestResponse(w, r, err)
		return
	}
	data, err := app.store.Users.GetById(r.Context(), userID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundResponse(w, r, err)
			return
		default:
			app.handleError(w, r, err)
			return
		}
	}
	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.handleError(w, r, err)
		return
	}
}
