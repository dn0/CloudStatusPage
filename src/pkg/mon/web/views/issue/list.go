package issue

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

const (
	perPageDefault  = 25
	perPageEmbedded = 5
)

//nolint:gochecknoglobals // This is a constant.
var cacheWhitelist = views.CacheAllowedQueryParams{
	"page": views.CacheAllowedAnyValues,
}

type listForm struct {
	Page int `form:"page,omitempty"`
}

type ListView struct {
	dbc      *db.Clients
	decoder  *views.FormDecoder
	perPage  int
	embedded bool
}

func NewListView(dbc *db.Clients) *ListView {
	return &ListView{
		dbc:      dbc,
		decoder:  views.NewFormDecoder(),
		perPage:  perPageDefault,
		embedded: false,
	}
}

func NewEmbeddedListView(dbc *db.Clients) *ListView {
	return &ListView{
		dbc:      dbc,
		decoder:  views.NewFormDecoder(),
		perPage:  perPageEmbedded,
		embedded: true,
	}
}

func (v *ListView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	views.CacheControl(w, r, views.CacheMaxAgeDefault, cacheWhitelist)
	c, err := v.handler(r)
	views.Render(w, r, c, err)
}

//nolint:wrapcheck // Error is properly logged by the caller.
func (v *ListView) handler(r *http.Request) (templ.Component, error) {
	query := &data.IssueQuery{}

	cloud, err := v.getClouds(r, query)
	if err != nil {
		return nil, err
	}

	region, err := v.getRegion(r, query, cloud)
	if err != nil {
		return nil, err
	}

	probe, err := v.getProbe(r, query, cloud)
	if err != nil {
		return nil, err
	}

	status, err := v.getStatus(r, query)
	if err != nil {
		return nil, err
	}

	var form listForm
	if errr := v.decoder.Decode(&form, r.URL.Query()); errr != nil {
		return nil, errr
	}
	query.Paginator = data.NewPaginator(form.Page, v.perPage)

	typ := chi.URLParam(r, "type")
	title, urlSuffix, issues, err := v.getIssues(r, typ, query)
	if err != nil {
		return nil, err
	}

	main := listTempl(
		title,
		urlSuffix,
		cloud,
		region,
		probe,
		status,
		issues,
		query,
		v.embedded,
	)
	if v.embedded {
		return main, nil
	}
	return templates.Base(templates.NavIssues, main), nil
}

//nolint:wrapcheck,nilnil // Error is properly logged by the caller.
func (v *ListView) getClouds(r *http.Request, query *data.IssueQuery) (*data.Cloud, error) {
	cloudParam := chi.URLParam(r, "cloud")
	if cloudParam == "" {
		query.Clouds = data.CloudIds
		return nil, nil
	}

	cloud, err := data.GetCloud(cloudParam)
	if err != nil {
		return nil, err
	}

	query.Clouds = []string{cloud.Id}
	return cloud, nil
}

//nolint:wrapcheck,nilnil // Error is properly logged by the caller.
func (v *ListView) getRegion(r *http.Request, query *data.IssueQuery, cloud *data.Cloud) (*data.CloudRegion, error) {
	if cloud == nil {
		return nil, nil
	}

	regionParam := chi.URLParam(r, "region")
	if regionParam == "" {
		return nil, nil
	}

	region, err := data.GetCloudRegion(r.Context(), v.dbc.Read, cloud.Id, regionParam)
	if err != nil {
		return nil, err
	}

	query.CloudRegion = region.Name
	return region, nil
}

//nolint:wrapcheck,nilnil // Error is properly logged by the caller.
func (v *ListView) getProbe(r *http.Request, query *data.IssueQuery, cloud *data.Cloud) (*data.ProbeDefinition, error) {
	if cloud == nil {
		return nil, nil
	}

	probeParam := chi.URLParam(r, "probe")
	if probeParam == "" {
		return nil, nil
	}

	probe, err := data.GetProbeDefinition(r.Context(), v.dbc.Read, cloud.Id, probeParam)
	if err != nil {
		return nil, err
	}

	query.ProbeName = probe.Name
	return probe, nil
}

func (v *ListView) getStatus(r *http.Request, query *data.IssueQuery) (pb.IncidentStatus, error) {
	var status pb.IncidentStatus
	switch chi.URLParam(r, "status") {
	case "", "all":
		status = pb.IncidentStatus_INCIDENT_ANY
	case "open":
		status = pb.IncidentStatus_INCIDENT_OPEN
	case "closed":
		status = pb.IncidentStatus_INCIDENT_CLOSED
	default:
		return 0, &data.InvalidInputError{Field: "status"}
	}

	query.Status = status
	return status, nil
}

func (v *ListView) getIssues(r *http.Request, typ string, query *data.IssueQuery) (string, string, []*data.Issue, error) {
	var title string
	var urlSuffix string
	var issues []*data.Issue
	var err error

	switch typ {
	case "issues":
		title = "Issues"
		urlSuffix = data.IssuesURLSuffix
		issues, err = query.GetAll(r.Context(), v.dbc.Read)
	case "alerts":
		title = "Alerts"
		urlSuffix = data.AlertsURLSuffix
		issues, err = query.GetAlerts(r.Context(), v.dbc.Read)
	case "incidents":
		title = "Incidents"
		urlSuffix = data.IncidentsURLSuffix
		issues, err = query.GetIncidents(r.Context(), v.dbc.Read)
	default:
		err = &data.InvalidInputError{Field: "type"}
	}

	return title, urlSuffix, issues, err
}
