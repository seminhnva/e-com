package main

import (
	"net/http"

	"github.com/seminhnva/e-com/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=1000"`
	PostID  int64  `json:"post_id" validate:"omitempty"`
	UserID  int64  `json:"user_id" validate:"required"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)
	var payload CreateCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badReqestResponse(w, r, err)
		return
	}

	if err := Validator.Struct(payload); err != nil {
		app.badReqestResponse(w, r, err)
		return
	}
	comment := &store.Comment{
		Content: payload.Content,
		PostID:  post.ID,
		UserID:  payload.UserID,
	}
	if err := app.store.Comments.Create(r.Context(), comment); err != nil {
		app.handleError(w, r, err)
		return
	}

	if error := app.jsonResponse(w, http.StatusCreated, comment); error != nil {
		app.internalServerError(w, r, error)
		return
	}
}
