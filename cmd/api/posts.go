package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/seminhnva/e-com/internal/store"
)

type postKey string

const postCtx postKey = "post"

type CreatePostPayload struct {
	Content string   `json:"content" validate:"required,max=1000"`
	Title   string   `json:"title" validate:"required,max=200"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Content *string `json:"content" validate:"omitempty,max=1000"`
	Title   *string `json:"title" validate:"omitempty,max=200"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badReqestResponse(w, r, err)
		return
	}
	if err := Validator.Struct(payload); err != nil {
		app.badReqestResponse(w, r, err)
		return
	}

	post := &store.Post{
		Content: payload.Content,
		Title:   payload.Title,
		// TODO: get change after auth
		UserID: 101,
		Tags:   payload.Tags,
	}
	if err := app.store.Posts.Create(r.Context(), post); err != nil {
		app.handleError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)
	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.handleError(w, r, err)
		return
	}
	post.Comments = comments
	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		app.badReqestResponse(w, r, err)
		return
	}
	if err = app.store.Posts.Delete(r.Context(), id); err != nil {
		app.handleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badReqestResponse(w, r, err)
		return
	}
	if err := Validator.Struct(payload); err != nil {
		app.badReqestResponse(w, r, err)
		return
	}

	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		app.handleError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postID := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(postID, 10, 64)
		if err != nil {
			app.badReqestResponse(w, r, err)
			return
		}
		post, err := app.store.Posts.GetByID(r.Context(), id)
		log.Printf("err: %v | ctx.Err: %v", err, r.Context().Err())

		if err != nil {
			app.handleError(w, r, err)
			return
		}
		ctx := context.WithValue(r.Context(), postCtx, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}
