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

type MatrixView struct {
	dbc *db.Clients
}

func NewMatrixView(dbc *db.Clients) *MatrixView {
	return &MatrixView{
		dbc: dbc,
	}
}

func (v *MatrixView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	views.CacheControl(w, r, views.CacheMaxAgeDefault, nil)
	c, err := v.handler(r)
	views.Render(w, r, c, err)
}

//nolint:wrapcheck // Error is properly logged by the caller.
func (v *MatrixView) handler(r *http.Request) (templ.Component, error) {
	cloud, err := data.GetCloud(chi.URLParam(r, "cloud"))
	if err != nil {
		return nil, err
	}

	regions, err := data.GetCloudRegions(r.Context(), v.dbc.Read, cloud.Id)
	if err != nil {
		return nil, err
	}

	probes, err := data.GetProbeDefinitions(r.Context(), v.dbc.Read, cloud.Id)
	if err != nil {
		return nil, err
	}
	probes = append(probes, data.NewPingProbeDefinition())

	issueQuery := data.IssueQuery{
		Clouds: []string{cloud.Id},
		Status: pb.IncidentStatus_INCIDENT_OPEN,
	}
	issues, err := issueQuery.GetAll(r.Context(), v.dbc.Read)
	if err != nil {
		return nil, err
	}

	//nolint:contextcheck // Confused linter.
	return templates.Base(templates.NavName(cloud.Id), matrixTempl(
		cloud,
		regions,
		probes,
		newIssueMap(issues),
	)), nil
}
