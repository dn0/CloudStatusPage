package probe

import (
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"cspage/pkg/config"
	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/mon/web/templates"
	"cspage/pkg/mon/web/views"
	"cspage/pkg/pb"
)

const (
	defaultTimeSpan = 24 * time.Hour * config.DefaultTimeSpanMultiplier
)

type ListRegionsView struct {
	dbc *db.Clients
}

func NewListRegionsView(dbc *db.Clients) *ListRegionsView {
	return &ListRegionsView{
		dbc: dbc,
	}
}

func (v *ListRegionsView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	views.CacheControl(w, r, views.CacheMaxAgeDefault, nil)
	c, err := v.handler(r)
	views.Render(w, r, c, err)
}

//nolint:wrapcheck // Error is properly logged by the caller, and it's a complicated handler indeed.
func (v *ListRegionsView) handler(r *http.Request) (templ.Component, error) {
	cloud, err := data.GetCloud(chi.URLParam(r, "cloud"))
	if err != nil {
		return nil, err
	}

	probe, err := data.GetProbeDefinition(r.Context(), v.dbc.Read, cloud.Id, chi.URLParam(r, "probe"))
	if err != nil {
		return nil, err
	}

	issueQuery := data.IssueQuery{
		Clouds:    []string{cloud.Id},
		Status:    pb.IncidentStatus_INCIDENT_OPEN,
		ProbeName: probe.Name,
	}
	issues, err := issueQuery.GetAll(r.Context(), v.dbc.Read)
	if err != nil {
		return nil, err
	}

	var main templ.Component
	if probe.IsPingDefinition() {
		main, err = v.agentTempl(r.Context(), cloud, probe, newIssueMap(issues))
	} else {
		main, err = v.probeTempl(r.Context(), cloud, probe, newIssueMap(issues))
	}
	if err != nil {
		return nil, err
	}

	//nolint:contextcheck // Confused linter.
	return templates.Base(templates.NavName(cloud.Id), main), nil
}
