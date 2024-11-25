package probe

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/mon/web/templates"
	"cspage/pkg/mon/web/views"
	"cspage/pkg/pb"
)

type ListProbesView struct {
	dbc *db.Clients
}

func NewListProbesView(dbc *db.Clients) *ListProbesView {
	return &ListProbesView{
		dbc: dbc,
	}
}

func (v *ListProbesView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	views.CacheControl(w, r, views.CacheMaxAgeDefault, nil)
	c, err := v.handler(r)
	views.Render(w, r, c, err)
}

//nolint:wrapcheck // Error is properly logged by the caller, and it's a complicated handler indeed.
func (v *ListProbesView) handler(r *http.Request) (templ.Component, error) {
	cloud, err := data.GetCloud(chi.URLParam(r, "cloud"))
	if err != nil {
		return nil, err
	}

	region, err := data.GetCloudRegion(r.Context(), v.dbc.Read, cloud.Id, chi.URLParam(r, "region"))
	if err != nil {
		return nil, err
	}

	probes, err := data.GetProbeDefinitions(r.Context(), v.dbc.Read, cloud.Id)
	if err != nil {
		return nil, err
	}
	probes = append(probes, data.NewPingProbeDefinition())

	issueQuery := data.IssueQuery{
		Clouds:      []string{cloud.Id},
		Status:      pb.IncidentStatus_INCIDENT_OPEN,
		CloudRegion: region.Name,
	}
	issues, err := issueQuery.GetAll(r.Context(), v.dbc.Read)
	if err != nil {
		return nil, err
	}

	//nolint:contextcheck // Confused linter.
	return templates.Base(templates.NavName(cloud.Id), listProbesTempl(
		cloud,
		region,
		probes,
		newIssueMap(issues),
	)), nil
}
