package middleware

import (
	"net/http"
	"time"

	"cspage/pkg/mon/web/views"
)

func CacheControl(
	maxAge time.Duration,
	whitelist views.CacheAllowedQueryParams,
) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			views.CacheControl(w, r, maxAge, whitelist)
			next.ServeHTTP(w, r)
		})
	}
}
