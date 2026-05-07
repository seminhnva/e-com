package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/seminhnva/e-com/internal/store"
)

type userKey string

const userCtx userKey = "user"

type FollowUser struct {
	UserID int64 `json:"user_id" validate:"required"`
}

// getUserHandler godoc
//
//	@Summary		Get user by ID
//	@Description	Fetch a single user by their ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int	true	"User ID"
//	@Success		200		{object}	store.User
//	@Failure		400		{object}	error	"Invalid user ID"
//	@Failure		404		{object}	error	"User not found"
//	@Failure		500		{object}	error	"Internal server error"
//	@Router			/users/{userID} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.handleError(w, r, err)
		return
	}
}

// followUserHandler godoc
//
//	@Summary		Follow a user
//	@Description	Follow a user by their ID
//	@Tags			users
//	@Accept			json
//	@Param			userID	path	int			true	"User ID to follow"
//	@Param			payload	body	FollowUser	true	"Follower user ID"
//	@Success		204		"No Content"
//	@Failure		400		{object}	error	"Invalid request"
//	@Failure		404		{object}	error	"User not found"
//	@Failure		409		{object}	error	"Already following"
//	@Failure		500		{object}	error	"Internal server error"
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser := getUserFromCtx(r)
	// TODO: revert back to authenticated user ID instead of passing in the payload
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badReqestResponse(w, r, err)
		return
	}

	if err := app.store.Follower.Follow(r.Context(), payload.UserID, followedUser.ID); err != nil {
		app.handleError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.handleError(w, r, err)
		return
	}
}

// unfollowUserHandler godoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by their ID
//	@Tags			users
//	@Accept			json
//	@Param			userID	path	int			true	"User ID to unfollow"
//	@Param			payload	body	FollowUser	true	"Follower user ID"
//	@Success		204		"No Content"
//	@Failure		400		{object}	error	"Invalid request"
//	@Failure		404		{object}	error	"User not found"
//	@Failure		500		{object}	error	"Internal server error"
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followedUser := getUserFromCtx(r)
	// TODO: revert back to authenticated user ID instead of passing in the payload

	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badReqestResponse(w, r, err)
		return
	}

	if err := app.store.Follower.Unfollow(r.Context(), payload.UserID, followedUser.ID); err != nil {
		app.handleError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.handleError(w, r, err)
		return
	}
}

func (app *application) userByIdContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
		if err != nil {
			app.badReqestResponse(w, r, err)
			return
		}
		user, err := app.store.Users.GetById(r.Context(), userID)
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
		ctx := context.WithValue(r.Context(), userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
