package store

import (
	"math"
	"net/http"
	"strconv"
	"strings"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=100"`
	Page   int    `json:"page" validate:"gte=1"`
	Sort   string `json:"sort" validate:"oneof=asc desc"`
	Search string `json:"search" validate:"max=100"`
}

func (fq *PaginatedFeedQuery) Offset() int {
	return (fq.Page - 1) * fq.Limit
}

func TotalPages(total, limit int) int {
	return int(math.Ceil(float64(total) / float64(limit)))
}

func (fq *PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	if v := qs.Get("limit"); v != "" {
		l, err := strconv.Atoi(v)
		if err != nil {
			return *fq, err
		}
		fq.Limit = l
	}
	if v := qs.Get("page"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return *fq, err
		}
		fq.Page = p
	}
	if v := qs.Get("sort"); v != "" {
		fq.Sort = v
	}
	if v := strings.TrimSpace(qs.Get("search")); v != "" {
		fq.Search = v
	}

	return *fq, nil
}
