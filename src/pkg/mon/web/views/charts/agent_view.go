package charts

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/go-echarts/go-echarts/v2/opts"

	"cspage/pkg/config"
	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/mon/web/views"
)

const (
	chartTypeDiff        = "diff"
	agentDefaultTimeSpan = 12 * time.Hour * config.DefaultTimeSpanMultiplier
	agentMaxTimeSpan     = 168 * time.Hour * config.DefaultTimeSpanMultiplier
)

//nolint:gochecknoglobals // This is a constant.
var jobMetrics = []*chartMetric{
	{
		name: "job_duration",
		fields: []dbField{
			{column: "ROUND(mon_job.took / 1000000::numeric, 2)", alias: "{probe_name}"},
		},
		chartName: "Monitoring Agent Job",
		chartAxis: "Duration",
		unit:      "ms",
	},
}

//nolint:gochecknoglobals // This is a constant.
var agentMetrics = []*chartMetric{
	{
		name: "proc_cpu",
		fields: []dbField{
			{column: "ROUND(proc_cpu_percent::numeric, 3)", alias: "proc_cpu_percent"},
		},
		chartName: "Monitoring Agent CPU Usage",
		chartAxis: "Process CPU",
		unit:      "%",
	},
	{
		name: "proc_mem",
		fields: []dbField{
			{column: "ROUND(proc_mem_rss / 1048576.0, 2)", alias: "proc_mem_rss"},
		},
		chartName: "Monitoring Agent Memory Usage",
		chartAxis: "Process memory",
		unit:      "MB",
	},
	{
		name: "os_cpu",
		fields: []dbField{
			// {column: "os_cpu_idle", alias: "idle"},
			{column: "os_cpu_user", alias: "user"},
			{column: "os_cpu_system", alias: "system"},
			{column: "os_cpu_iowait", alias: "iowait"},
			{column: "os_cpu_nice", alias: "nice"},
			{column: "os_cpu_irq", alias: "irq"},
			{column: "os_cpu_softirq", alias: "softirq"},
			{column: "os_cpu_steal", alias: "steal"},
		},
		chartType: chartTypeDiff,
		chartName: "Monitoring Agent's OS CPU Usage",
		chartAxis: "OS CPU time",
		unit:      "sec",
	},
	{
		name: "os_mem",
		fields: []dbField{
			{column: "ROUND(os_mem_free / 1048576.0, 2)", alias: "free"},
			{column: "ROUND(os_mem_cached / 1048576.0, 2)", alias: "cached"},
			{column: "ROUND(os_mem_buffers / 1048576.0, 2)", alias: "buffers"},
			{column: "ROUND(os_mem_used / 1048576.0, 2)", alias: "used"},
		},
		chartName: "Monitoring Agent's OS Memory Usage",
		chartAxis: "OS memory",
		unit:      "MB",
	},
}

//nolint:gochecknoglobals // This is a constant.
var agentCacheWhitelist = views.CacheAllowedQueryParams{
	"since":    views.CacheAllowedAnyValues,
	"until":    views.CacheAllowedAnyValues,
	"metric":   views.CacheAllowedAnyValues,
	"alert":    views.CacheAllowedAnyValues,
	"incident": views.CacheAllowedAnyValues,
}

type dbField struct {
	column string
	alias  string
}

type chartMetric struct {
	name      string
	fields    []dbField
	chartName string
	chartType string
	chartAxis string
	unit      string
}

type agentViewForm struct {
	baseForm
	Metric string `form:"metric,omitempty"`
}

type agentView struct {
	baseView
}

func newAgentView(dbc *db.Clients) *agentView {
	return &agentView{baseView{
		dbc:     dbc,
		encoder: views.NewFormEncoder(),
		decoder: views.NewFormDecoder(),
	}}
}

func (v *agentView) templ(
	r *http.Request,
	cloud *data.Cloud,
	region *data.CloudRegion,
	probe *data.ProbeDefinition,
) (templ.Component, error) {
	form, baseQuery, metrics, err := v.validateForm(r)
	if err != nil {
		return nil, err
	}

	issue, err := v.validateBaseForm(r, cloud, &form.baseForm, agentDefaultTimeSpan, agentMaxTimeSpan)
	if err != nil {
		return nil, err
	}

	chartSeries, err := v.getData(r.Context(), cloud.Id, region, probe, form, baseQuery, metrics)
	if err != nil {
		return nil, err
	}

	drawCharts := make([]*lineChart, len(metrics))
	for i, metric := range metrics {
		chart, err := v.makeChart(cloud, region, probe, *form, issue, metric, chartSeries)
		if err != nil {
			return nil, err
		}
		drawCharts[i] = chart
	}

	return chartsTempl(drawCharts), nil
}

//nolint:wrapcheck // Error is properly logged by the caller.
func (v *agentView) validateForm(r *http.Request) (*agentViewForm, string, []*chartMetric, error) {
	var form agentViewForm
	if err := v.decoder.Decode(&form, r.URL.Query()); err != nil {
		return nil, "", nil, err
	}

	sql, metrics := validateMetric(form.Metric)

	if sql == "" || len(metrics) == 0 {
		return nil, "", nil, &data.InvalidInputError{Field: "metric"}
	}

	return &form, sql, metrics, nil
}

//nolint:wrapcheck // Error is properly logged by the caller.
func (v *agentView) makeChart(
	cloud *data.Cloud,
	region *data.CloudRegion,
	probe *data.ProbeDefinition,
	form agentViewForm,
	issue *data.Issue,
	metric *chartMetric,
	series map[string][]opts.LineData,
) (*lineChart, error) {
	title := fmt.Sprintf("%s %s in %s", cloud.Name, metric.chartName, region.Name)
	axisName := fmt.Sprintf("%s (%s)", metric.chartAxis, metric.unit)

	form.Metric = metric.name
	qs, err := v.encoder.Encode(&form)
	if err != nil {
		return nil, err
	}
	link := probe.DetailsURL(cloud, region) + "?" + qs.Encode()
	// field aliases can contain placeholders
	replacer := strings.NewReplacer(
		"{probe_name}", probe.Name,
	)

	probeInterval := probe.Config.Interval()
	chart := newLineChart(title, link, axisName, metric.unit)
	for _, field := range metric.fields {
		lineData := series[field.alias]
		if metric.chartType == chartTypeDiff && len(lineData) > 0 {
			lineData = diffSeries[float32](lineData)
		}
		if form.untilNow {
			lineData = appendDummyLineItem(lineData, form.Until)
		}
		chart.addSeries(replacer.Replace(field.alias), lineData, defaultStack, probeInterval, issue)
	}

	return chart, nil
}

func (v *agentView) numCharts(r *http.Request) int {
	if r.URL.Query().Has("metric") {
		return 1
	}
	return len(agentMetrics)
}

func validateMetric(name string) (string, []*chartMetric) {
	if name == "" {
		return sqlSelectAgentStats, agentMetrics
	}

	for _, metric := range agentMetrics {
		if metric.name == name {
			return sqlSelectAgentStats, []*chartMetric{metric}
		}
	}
	for _, metric := range jobMetrics {
		if metric.name == name {
			return sqlSelectJobStats, []*chartMetric{metric}
		}
	}

	return "", nil
}
