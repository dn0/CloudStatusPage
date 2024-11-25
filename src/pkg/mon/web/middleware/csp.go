package middleware

import (
	"net/http"
)

//nolint:lll // Long CSP header.
const (
	cspHeaderName  = "Content-Security-Policy"
	cspHeaderValue = "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com;"
)

func CSP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(cspHeaderName, cspHeaderValue)
		next.ServeHTTP(w, r)
	})
}
