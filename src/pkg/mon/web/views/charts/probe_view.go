package charts

import (
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/go-echarts/go-echarts/v2/opts"

	"cspage/pkg/config"
	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/mon/web/views"
)

const (
	probeDefaultTimeSpan = 25 * time.Hour * config.DefaultTimeSpanMultiplier
	probeMaxTimeSpan     = 336 * time.Hour * config.DefaultTimeSpanMultiplier
)

//nolint:gochecknoglobals // This is a constant.
var probeCacheWhitelist = views.CacheAllowedQueryParams{
	"since":    views.CacheAllowedAnyValues,
	"until":    views.CacheAllowedAnyValues,
	"action":   views.CacheAllowedAnyValues,
	"alert":    views.CacheAllowedAnyValues,
	"incident": views.CacheAllowedAnyValues,
}

type probeAction struct {
	action *data.ProbeAction
	series []opts.LineData
	errors []opts.ScatterData
}

type probeViewForm struct {
	baseForm
	Actions []uint32 `form:"action,omitempty"`
}

type probeView struct {
	baseView
}

func newProbeView(dbc *db.Clients) *probeView {
	return &probeView{baseView{
		dbc:     dbc,
		encoder: views.NewFormEncoder(),
		decoder: views.NewFormDecoder(),
	}}
}

func (v *probeView) templ(
	r *http.Request,
	cloud *data.Cloud,
	region *data.CloudRegion,
	probe *data.ProbeDefinition,
) (templ.Component, error) {
	form, probeResults, err := v.validateForm(r, probe)
	if err != nil {
		return nil, err
	}

	issue, err := v.validateBaseForm(r, cloud, &form.baseForm, probeDefaultTimeSpan, probeMaxTimeSpan)
	if err != nil {
		return nil, err
	}

	if err := v.getProbeResults(r.Context(), cloud.Id, region, probe, form, probeResults); err != nil {
		return nil, err
	}

	//nolint:prealloc // Can't preallocate in this form.
	var drawCharts []*lineChart
	for _, actionGroup := range data.GroupActionIDs(form.Actions) {
		chart, err := v.makeChart(cloud, region, probe, *form, issue, actionGroup, probeResults)
		if err != nil {
			return nil, err
		}
		drawCharts = append(drawCharts, chart)
	}

	return chartsTempl(drawCharts), nil
}

//nolint:wrapcheck // Error is properly logged by the caller.
func (v *probeView) validateForm(
	r *http.Request,
	probe *data.ProbeDefinition,
) (*probeViewForm, map[uint32]*probeAction, error) {
	var form probeViewForm
	if err := v.decoder.Decode(&form, r.URL.Query()); err != nil {
		return nil, nil, err
	}

	actions := make([]uint32, 0)
	probeResults := map[uint32]*probeAction{}

	if len(form.Actions) > 0 {
		actionMap := probe.Config.ActionMap()
		for _, formAction := range form.Actions {
			action, ok := actionMap[formAction]
			if !ok {
				return nil, nil, &data.InvalidInputError{Field: "action"}
			}
			actions = append(actions, action.Id)
			probeResults[action.Id] = &probeAction{action: action}
		}
	} else {
		for _, action := range probe.Config.Actions {
			if action.Id > data.ProbeMaxDisplayActionId {
				continue
			}
			actions = append(actions, action.Id)
			probeResults[action.Id] = &probeAction{action: &action}
		}
	}

	form.Actions = actions

	return &form, probeResults, nil
}

//nolint:wrapcheck,varnamelen // Error is properly logged by the caller.
func (v *probeView) makeChart(
	cloud *data.Cloud,
	region *data.CloudRegion,
	probe *data.ProbeDefinition,
	form probeViewForm,
	issue *data.Issue,
	actions []uint32,
	probeResults map[uint32]*probeAction,
) (*lineChart, error) {
	first := probeResults[actions[0]]
	title := first.action.FullGroupName(probe) + " in " + region.Name

	form.Actions = actions
	qs, err := v.encoder.Encode(&form)
	if err != nil {
		return nil, err
	}
	link := probe.DetailsURL(cloud, region) + "?" + qs.Encode()

	probeInterval := probe.Config.Interval()
	chart := newLineChart(title, link, "Duration (ms)", "ms")
	chartErrors := newScatterChart()
	for i, aid := range actions {
		pr := probeResults[aid]
		if form.untilNow {
			pr.series = appendDummyLineItem(pr.series, form.Until)
		}
		chart.addSeries(pr.action.Name, pr.series, pr.action.ChartStack, probeInterval, issue)
		if pr.errors != nil {
			chartErrors.addSeries("error:"+pr.action.Name, chartErrorColors[i%len(chartErrorColors)], pr.errors)
		}
	}
	chart.Overlap(chartErrors)

	return chart, nil
}

func (v *probeView) numCharts(r *http.Request, probe *data.ProbeDefinition) int {
	if r.URL.Query().Has("actions") {
		return 1
	}
	return len(probe.Config.ActionGroupIDs())
}
