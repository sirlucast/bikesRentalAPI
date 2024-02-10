package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type (
	pageKey string
)

const (
	PageIDKey pageKey = "page_id"
)

// Pagination middleware is used to extract the next page id from the url query
func Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		PageID := r.URL.Query().Get(string(PageIDKey))
		var intPageID int64 = 0
		var err error
		if PageID != "" {
			intPageID, err = strconv.ParseInt(PageID, 10, 64)
			if err != nil {
				http.Error(w, fmt.Sprintf("couldn't read %s: %v", PageIDKey, err), http.StatusBadRequest)
				return
			}
		}
		ctx := context.WithValue(r.Context(), PageIDKey, intPageID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
