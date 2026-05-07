package main

import (
	"net/http"

	"github.com/seminhnva/e-com/internal/store"
)

// getFeedHandler godoc
//
//	@Summary		Get user feed
//	@Description	Returns a paginated feed of posts from followed users
//	@Tags			feed
//	@Produce		json
//	@Param			limit	query		int		false	"Number of posts per page (1-100)"	default(10)
//	@Param			page	query		int		false	"Page number"						default(1)
//	@Param			sort	query		string	false	"Sort order (asc or desc)"			default(desc)
//	@Param			search	query		string	false	"Search keyword"
//	@Success		200		{object}	map[string]any
//	@Failure		400		{object}	error	"Invalid query params"
//	@Failure		500		{object}	error	"Internal server error"
//	@Router			/feed [get]
func (app *application) getFeedHandler(w http.ResponseWriter, r *http.Request) {
	fq := store.PaginatedFeedQuery{
		Limit:  10,
		Page:   1,
		Sort:   "desc",
		Search: "",
	}
	fq, err := fq.Parse(r)
	if err != nil {
		app.handleError(w, r, err)
		return
	}
	if err := Validator.Struct(fq); err != nil {
		app.handleError(w, r, err)
		return
	}

	posts, total, err := app.store.Posts.GetUserFeed(r.Context(), int64(102), fq)
	if err != nil {
		app.handleError(w, r, err)
		return
	}

	resp := map[string]any{
		"data":        posts,
		"total":       total,
		"page":        fq.Page,
		"limit":       fq.Limit,
		"total_pages": store.TotalPages(total, fq.Limit),
	}
	if err := app.jsonResponse(w, http.StatusOK, resp); err != nil {
		app.handleError(w, r, err)
		return
	}
}
