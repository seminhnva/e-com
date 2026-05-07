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

// createCommentHandler godoc
//
//	@Summary		Create a comment
//	@Description	Add a comment to a post
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		int						true	"Post ID"
//	@Param			payload	body		CreateCommentPayload	true	"Comment payload"
//	@Success		201		{object}	store.Comment
//	@Failure		400		{object}	error	"Invalid request body"
//	@Failure		404		{object}	error	"Post not found"
//	@Failure		500		{object}	error	"Internal server error"
//	@Router			/posts/{postID}/comments [post]
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
