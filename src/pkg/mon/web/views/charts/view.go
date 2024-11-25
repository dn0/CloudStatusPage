package charts

import (
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/mon/web/views"
)

type baseForm struct {
	untilNow bool

	Since    time.Time `form:"since,omitempty"`
	Until    time.Time `form:"until,omitempty"`
	Alert    string    `form:"alert,omitempty"`
	Incident string    `form:"incident,omitempty"`
}

type baseView struct {
	dbc     *db.Clients
	encoder *views.FormEncoder
	decoder *views.FormDecoder
}

type View struct {
	agent *agentView
	probe *probeView
}

func NewView(dbc *db.Clients) *View {
	return &View{
		agent: newAgentView(dbc),
		probe: newProbeView(dbc),
	}
}

func (v *View) NumCharts(r *http.Request, probe *data.ProbeDefinition) int {
	if probe.IsPingDefinition() {
		return v.agent.numCharts(r)
	}
	return v.probe.numCharts(r, probe)
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := v.handler(w, r)
	views.Render(w, r, c, err)
}

//nolint:wrapcheck // Error is properly logged by the caller.
func (v *View) handler(w http.ResponseWriter, r *http.Request) (templ.Component, error) {
	cloud, err := data.GetCloud(chi.URLParam(r, "cloud"))
	if err != nil {
		return nil, err
	}

	region, err := data.GetCloudRegion(r.Context(), v.probe.dbc.Read, cloud.Id, chi.URLParam(r, "region"))
	if err != nil {
		return nil, err
	}

	probe, err := data.GetProbeDefinition(r.Context(), v.probe.dbc.Read, cloud.Id, chi.URLParam(r, "probe"))
	if err != nil {
		return nil, err
	}

	if probe.IsPingDefinition() {
		views.CacheControl(w, r, views.CacheMaxAgeDefault, agentCacheWhitelist)
		return v.agent.templ(r, cloud, region, probe)
	}
	views.CacheControl(w, r, views.CacheMaxAgeDefault, probeCacheWhitelist)
	return v.probe.templ(r, cloud, region, probe)
}

//nolint:wrapcheck // Error is properly logged by the caller.
func (v *baseView) getIssue(r *http.Request, cloud *data.Cloud, form *baseForm) (*data.Issue, error) {
	switch {
	case form.Alert != "":
		return data.GetAlertIssue(r.Context(), v.dbc.Read, cloud.Id, form.Alert)
	case form.Incident != "":
		return data.GetIncidentIssue(r.Context(), v.dbc.Read, cloud.Id, form.Incident)
	default:
		//nolint:nilnil // I'm in a hurry, don't have time to implement sentinels.
		return nil, nil
	}
}

func (v *baseView) validateBaseForm(
	r *http.Request,
	cloud *data.Cloud,
	form *baseForm,
	defaultTimeSpan, maxTimeSpan time.Duration,
) (*data.Issue, error) {
	if form.Until.IsZero() {
		form.Until = time.Now()
		form.untilNow = true
	}

	if form.Since.IsZero() {
		form.Since = form.Until.Add(-defaultTimeSpan)
	}

	if form.Until.Sub(form.Since) > maxTimeSpan {
		form.Since = form.Until.Add(-maxTimeSpan)
	}

	return v.getIssue(r, cloud, form)
}
