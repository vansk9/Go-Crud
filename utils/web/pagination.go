package web

import (
	"net/http"
	"strconv"
)

// PaginationParams contains the pagination parameters from the request.
type PaginationParams struct {
	Page     int
	PageSize int
}

// DefaultPageSize is the default number of items per page.
const DefaultPageSize = 10

// NewPaginationParams reads "page" and "pageSize" from the query string
// and applies defaults/validation.
func NewPaginationParams(r *http.Request) PaginationParams {
	q := r.URL.Query()

	page := 1
	if p := q.Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}

	pageSize := DefaultPageSize
	if ps := q.Get("pageSize"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			pageSize = v
		}
	}

	return PaginationParams{Page: page, PageSize: pageSize}
}

// CalculateOffset returns the SQL offset for LIMIT/OFFSET queries.
func (p PaginationParams) CalculateOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetPaginationResponse builds the PaginationResponse, using the incoming request URL
// to generate "next" and "prev" links.
func (p PaginationParams) GetPaginationResponse(r *http.Request, total int64) PaginationResponse {
	return PaginationResponse{
		Prev:  p.prevURL(r),
		Next:  p.nextURL(r, total),
		Total: total,
	}
}

// prevURL returns the URL for the previous page, or empty if on the first page.
func (p PaginationParams) prevURL(r *http.Request) string {
	if p.Page <= 1 {
		return ""
	}
	u := *r.URL
	q := u.Query()
	q.Set("page", strconv.Itoa(p.Page-1))
	q.Set("pageSize", strconv.Itoa(p.PageSize))
	u.RawQuery = q.Encode()
	return u.String()
}

// nextURL returns the URL for the next page, or empty if there are no more items.
func (p PaginationParams) nextURL(r *http.Request, total int64) string {
	if total <= int64(p.Page*p.PageSize) {
		return ""
	}
	u := *r.URL
	q := u.Query()
	q.Set("page", strconv.Itoa(p.Page+1))
	q.Set("pageSize", strconv.Itoa(p.PageSize))
	u.RawQuery = q.Encode()
	return u.String()
}

// PaginationResponse is returned in the JSON envelope.
type PaginationResponse struct {
	Prev  string `json:"prev"`
	Next  string `json:"next"`
	Total int64  `json:"total"`
}
