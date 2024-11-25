package web

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"

	"cspage/pkg/config"
	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/mon/web/middleware"
	"cspage/pkg/mon/web/templates"
	"cspage/pkg/mon/web/views"
	"cspage/pkg/mon/web/views/about"
	"cspage/pkg/mon/web/views/agent"
	"cspage/pkg/mon/web/views/charts"
	"cspage/pkg/mon/web/views/home"
	"cspage/pkg/mon/web/views/issue"
	"cspage/pkg/mon/web/views/probe"
)

const (
	staticDir         = "./srv/mon-web"
	cacheStaticMaxAge = 24 * time.Hour
)

//nolint:gochecknoglobals // These are constants.
var (
	pong                 = []byte("pong")
	cacheStaticWhitelist = views.CacheAllowedQueryParams{
		"version": views.CacheAllowedQueryValues{
			config.Version: struct{}{},
			"0":            struct{}{},
			"1":            struct{}{},
			"2":            struct{}{},
		},
	}
)

//nolint:lll,varnamelen // Long lines and r(outer) for better readability.
func Handlers(cfg *Config, dbc *db.Clients) map[string]http.Handler {
	views.CachingDisabled = cfg.Debug
	issueListView := issue.NewListView(dbc)
	issueDetailsView := issue.NewDetailsView(dbc)
	chartsView := charts.NewView(dbc)
	probeIssuesView := issue.NewEmbeddedListView(dbc)
	probeDetailsView := probe.NewDetailsView(dbc, chartsView)

	r := chi.NewRouter()
	r.Use(defaultMiddleware(web, cfg.Debug)...)
	r.Route("/cloud/{cloud:^("+strings.Join(data.CloudIds, "|")+")$}", func(r chi.Router) {
		// Probe charts
		r.Get("/region/{region:^[A-Za-z0-9-]+$}/probe/{probe:^[A-Za-z0-9_]+$}/charts", chartsView.ServeHTTP)
		// Probe issues, embedded variant
		r.Get("/region/{region:^[A-Za-z0-9-]+$}/probe/{probe:^[A-Za-z0-9_]+$}/{type:issues}/embedded", probeIssuesView.ServeHTTP)
		// Details (details + issues + charts)
		r.Get("/region/{region:^[A-Za-z0-9-]+$}/probe/{probe:^[A-Za-z0-9_]+$}", probeDetailsView.ServeHTTP)
		// Tables
		r.Get("/region/{region:^[A-Za-z0-9-]+$}", probe.NewListProbesView(dbc).ServeHTTP)
		r.Get("/probe/{probe:^[A-Za-z0-9_]+$}", probe.NewListRegionsView(dbc).ServeHTTP)
		r.Get("/", probe.NewMatrixView(dbc).ServeHTTP)
		// Issues
		r.Get("/region/{region:^[A-Za-z0-9-]+$}/{type:^(issues|alerts|incidents)$}/{status:^(open|closed|all)$}", issueListView.ServeHTTP)
		r.Get("/region/{region:^[A-Za-z0-9-]+$}/{type:^(issues|alerts|incidents)$}", issueListView.ServeHTTP)
		r.Get("/{type:^(issues|alerts|incidents)$}/{status:^(open|closed|all)$}", issueListView.ServeHTTP)
		r.Get("/{type:^(issues|alerts|incidents)$}", issueListView.ServeHTTP)
		// Issue details
		r.With(middleware.ValidateUUID).Get("/{type:^(alert|incident)$}/{id}", issueDetailsView.ServeHTTP)
	})
	r.Get("/{type:^(issues|alerts|incidents)$}/{status:^(open|closed|all)$}", issueListView.ServeHTTP)
	r.Get("/{type:^(issues|alerts|incidents)$}", issueListView.ServeHTTP)
	r.Get("/world", charts.NewMapView(dbc).ServeHTTP)
	r.Get("/about", about.NewView().ServeHTTP)
	r.Get("/", home.NewView().ServeHTTP)

	return map[string]http.Handler{
		"/":                 r,
		templates.StaticURL: newStaticServer(staticDir, templates.StaticURL, defaultMiddleware(static, cfg.Debug)...),
		"/api":              apiHandler(cfg, dbc),
	}
}

//nolint:varnamelen // Long lines and r(outer) for better readability.
func apiHandler(cfg *Config, dbc *db.Clients) http.Handler {
	r := chi.NewRouter()
	r.Use(defaultMiddleware(api, cfg.Debug)...)
	r.Get("/ping", apiPing)
	r.Route("/cloud/{cloud:^("+strings.Join(data.CloudIds, "|")+")$}", func(r chi.Router) {
		// Secret Area
		r.Group(func(r chi.Router) {
			r.Use(basicAuth(cfg))
			// List of running agents
			r.Get("/agent", agent.NewListView(dbc).ServeHTTP)
		})
	})
	return r
}

func apiPing(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pong)
}
