package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/seminhnva/e-com/internal/store"
)

type CreatePostPayload struct {
	Content string   `json:"content" validate:"required,max=1000"`
	Title   string   `json:"title" validate:"required,max=200"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := Validator.Struct(payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	post := &store.Post{
		Content: payload.Content,
		Title:   payload.Title,
		// TODO: get change after auth
		UserID: 1,
		Tags:   []string{"tag1", "tag2"},
	}
	err := app.store.Posts.Create(r.Context(), post)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		app.badReqestResponse(w, r, err)
		return
	}
	post, err := app.store.Posts.GetByID(r.Context(), id)
	switch {
	case errors.Is(err, store.ErrNotFound):
		app.notFoundResponse(w, r, err)
		return
	case err != nil:
		app.internalServerError(w, r, err)
		return
	}

	comments, err := app.store.Comments.GetByPostID(r.Context(), id)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	post.Comments = comments
	if err := writeJSON(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
