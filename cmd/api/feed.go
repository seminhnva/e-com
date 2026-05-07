package main

import (
	"net/http"

	"github.com/seminhnva/e-com/internal/store"
)

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
