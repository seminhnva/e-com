package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/seminhnva/e-com/internal/store"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("internal server error: %s path: %s error: %v", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}
func (app *application) badReqestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request: %s path: %s error: %v", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found: %s path: %s error: %v", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusNotFound, "resource not found")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("conflict: %s path: %s error: %v", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusConflict, "resource was updated by another request, please retry")
}

func (app *application) timeoutResponse(w http.ResponseWriter, r *http.Request) {
	app.logger.Warnf("request timeout: %s path: %s", r.Method, r.URL.Path)
	writeJSONError(w, http.StatusServiceUnavailable, "request timeout, please try again")
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}
func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) handleError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, context.DeadlineExceeded), r.Context().Err() == context.DeadlineExceeded:
		app.timeoutResponse(w, r)
	case errors.Is(err, context.Canceled), r.Context().Err() == context.Canceled:
		return
	case errors.Is(err, store.ErrNotFound):
		app.notFoundResponse(w, r, err)
	case errors.Is(err, store.ErrEditConflict):
		app.conflictResponse(w, r, err)
	case errors.Is(err, store.ErrConflict):
		writeJSONError(w, http.StatusConflict, err.Error())
	case errors.Is(err, store.ErrDuplicateEmail):
		writeJSONError(w, http.StatusConflict, err.Error())
	case errors.Is(err, store.ErrDuplicateUsername):
		writeJSONError(w, http.StatusConflict, err.Error())
	default:
		app.internalServerError(w, r, err)
	}
}
