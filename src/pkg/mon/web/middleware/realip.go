package middleware

import (
	"net/http"
	"strings"
)

//nolint:gochecknoglobals // These are constants.
var (
	xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	xRealIP       = http.CanonicalHeaderKey("X-Real-IP")
)

func RealIP(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rip := realIP(r); rip != "" {
			r.RemoteAddr = rip
		}
		h.ServeHTTP(w, r)
	})
}

func realIP(r *http.Request) string {
	if ip := r.Header.Get(xRealIP); ip != "" {
		return ip
	} else if xff := r.Header.Get(xForwardedFor); xff != "" {
		// X-Forwarded-For: <supplied-value>,<client-ip>,<load-balancer-ip><GFE-IP><backend-IP>
		ips := strings.Split(xff, ",")
		if len(ips) > 1 {
			return ips[len(ips)-2]
		}
	}
	return ""
}
