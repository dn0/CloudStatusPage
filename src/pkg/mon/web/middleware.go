package web

import (
	"net/http"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"

	"cspage/pkg/mon/web/middleware"
)

const (
	static contentType = "STATIC"
	api    contentType = "API"
	web    contentType = "WEB"
)

type contentType string

type middlewareFun = func(http.Handler) http.Handler

func defaultMiddleware(t contentType, debug bool) []middlewareFun {
	ware := []middlewareFun{
		middleware.RealIP,
	}
	switch t {
	case api:
		ware = append(
			ware,
			chiMiddleware.GetHead,
			middleware.ContentTypeApplicationJSON,
		)
	case static:
		ware = append(
			ware,
			chiMiddleware.GetHead,
			// http.FileServer should remove the 'Cache-Control' header when the file does not exist
			middleware.CacheControl(cacheStaticMaxAge, cacheStaticWhitelist),
		)
	case web:
		ware = append(
			ware,
			chiMiddleware.RedirectSlashes,
			chiMiddleware.GetHead,
			middleware.ContentTypeTextHTML,
			middleware.CSP,
		)
	}
	if !debug {
		ware = append(
			ware,
			middleware.StrictTransportSecurity,
		)
	}
	return ware
}

func basicAuth(cfg *Config) middlewareFun {
	return chiMiddleware.BasicAuth(cfg.HTTPBasicAuthRealm, map[string]string{
		cfg.HTTPBasicAuthUsername: cfg.HTTPBasicAuthPassword,
	})
}
