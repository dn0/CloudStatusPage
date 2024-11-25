package issue

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/mon/web/templates"
	"cspage/pkg/mon/web/views"
)

type DetailsView struct {
	dbc *db.Clients
}

func NewDetailsView(dbc *db.Clients) *DetailsView {
	return &DetailsView{
		dbc: dbc,
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

	title, issue, relatedIssues, err := v.getIssue(r, cloud.Id, chi.URLParam(r, "type"), chi.URLParam(r, "id"))
	if err != nil {
		return nil, err
	}

	return templates.Base(templates.NavIssues, detailsTempl(
		title,
		cloud,
		issue,
		relatedIssues,
	)), nil
}

func (v *DetailsView) getIssue(r *http.Request, cloud, typ, id string) (string, *data.Issue, []*data.Issue, error) {
	var title string
	var issue *data.Issue
	var related []*data.Issue
	var err error

	switch typ {
	case "alert":
		title = "Alerts"
		if issue, err = data.GetAlertIssue(r.Context(), v.dbc.Read, cloud, id); err == nil && issue.AlertIncidentId != "" {
			query := data.IssueQuery{
				Clouds:     []string{cloud},
				IncidentId: issue.AlertIncidentId,
			}
			related, err = query.GetIncidents(r.Context(), v.dbc.Read)
		}
	case "incident":
		title = "Incidents"
		if issue, err = data.GetIncidentIssue(r.Context(), v.dbc.Read, cloud, id); err == nil {
			query := data.IssueQuery{
				Clouds:     []string{cloud},
				IncidentId: issue.Id,
			}
			related, err = query.GetAlerts(r.Context(), v.dbc.Read)
		}
	default:
		err = &data.InvalidInputError{Field: "type"}
	}

	return title, issue, related, err
}
