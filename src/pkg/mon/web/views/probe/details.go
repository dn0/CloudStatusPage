package probe

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/mon/web/templates"
	"cspage/pkg/mon/web/views"
	"cspage/pkg/mon/web/views/charts"
	"cspage/pkg/pb"
)

type DetailsView struct {
	dbc    *db.Clients
	charts *charts.View
}

func NewDetailsView(dbc *db.Clients, chartsView *charts.View) *DetailsView {
	return &DetailsView{
		dbc:    dbc,
		charts: chartsView,
	}
}

func (v *DetailsView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	views.CacheControl(w, r, views.CacheMaxAgeDefault, nil)
	c, err := v.handler(r)
	views.Render(w, r, c, err)
}

//nolint:wrapcheck // Error is properly logged by the caller.
func (v *DetailsView) handler(r *http.Request) (templ.Component, error) {
	cloud, err := data.GetCloud(chi.URLParam(r, "cloud"))
	if err != nil {
		return nil, err
	}

	region, err := data.GetCloudRegion(r.Context(), v.dbc.Read, cloud.Id, chi.URLParam(r, "region"))
	if err != nil {
		return nil, err
	}

	probe, err := data.GetProbeDefinition(r.Context(), v.dbc.Read, cloud.Id, chi.URLParam(r, "probe"))
	if err != nil {
		return nil, err
	}

	job, err := data.GetLastJob(r.Context(), v.dbc.Read, cloud.Id, probe.Name, region.Name)
	if err != nil {
		return nil, err
	}

	details, err := v.getAgentDetails(r.Context(), cloud.Id, job)
	if err != nil {
		return nil, err
	}

	issueQuery := newIssueQuery(cloud, region, probe)
	issueCount, err := issueQuery.CountAll(r.Context(), v.dbc.Read)
	if err != nil {
		return nil, err
	}

	//nolint:contextcheck // Confused linter.
	return templates.Base(templates.NavName(cloud.Id), detailsTempl(
		cloud,
		region,
		probe,
		job,
		details,
		issueCount.Total(),
		v.charts.NumCharts(r, probe),
		r.URL.Query(),
	)), nil
}

func newIssueQuery(cloud *data.Cloud, region *data.CloudRegion, probe *data.ProbeDefinition) *data.IssueQuery {
	return &data.IssueQuery{
		Clouds:      []string{cloud.Id},
		Status:      pb.IncidentStatus_INCIDENT_ANY,
		CloudRegion: region.Name,
		ProbeName:   probe.Name,
	}
}
