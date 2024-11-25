package charts

import (
	"strings"
	"time"

	echarts "github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"

	"cspage/pkg/data"
)

const (
	dataZoomStart             = 90
	dataZoomMinCount          = 800
	markAreaIssueTimeMod      = 2
	lineChartTooltipFormatter = "echartLineTooltipFormatterFactory('{unit}')"
	defaultStack              = "stack"
	noStack                   = "no"
)

func newLineChart(title, titleLink, yAxisName, unit string) *lineChart {
	chart := &lineChart{*echarts.NewLine()}
	chart.Initialization.Theme = "auto" // To avoid automatic injections of "light" color scheme
	chart.Title.Title = title
	chart.Title.Link = titleLink
	chart.XYAxis.XAxisList[0].Type = "time"
	chart.XYAxis.YAxisList[0].Name = yAxisName
	chart.Tooltip.Trigger = "axis"
	chart.Tooltip.Formatter = opts.FuncOpts(strings.Replace(lineChartTooltipFormatter, "{unit}", unit, 1))
	chart.Toolbox.Feature = &opts.ToolBoxFeature{
		Restore: &opts.ToolBoxFeatureRestore{},
	}
	chart.DataZoomList = []opts.DataZoom{{
		Type: "slider",
	}}
	return chart
}

// Only for overlay (not for rendering directly).
func newScatterChart() *scatterChart {
	return &scatterChart{*echarts.NewScatter()}
}

func (c *lineChart) setDefaultDataZoom(dataCount int) {
	if dataCount < dataZoomMinCount {
		return
	}
	c.DataZoomList[0].Start = dataZoomStart
}

//nolint:mnd // Magic numbers are for style.
func (c *lineChart) addSeries(name string, lineData []opts.LineData, stack string, freq time.Duration, issue *data.Issue) {
	// Set zoom according to first series length
	if len(c.MultiSeries) == 0 && issue == nil {
		c.setDefaultDataZoom(len(lineData))
	}

	series := echarts.SingleSeries{
		Name:       name,
		Type:       types.ChartLine,
		Data:       lineData,
		Smooth:     opts.Bool(false),
		Symbol:     "circle",
		SymbolSize: 2,
		LineStyle: &opts.LineStyle{
			Width: 1.0,
		},
	}

	if stack == "" {
		stack = defaultStack
	}
	if stack != noStack {
		series.Stack = stack
		series.LineStyle.Width = 0.8
		series.AreaStyle = &opts.AreaStyle{
			Opacity: 0.4,
		}
	}

	if issue == nil {
		series.MarkLines = newMarkLineAverage()
	} else {
		series.MarkLines = newMarkLineFromCloudIssue(issue)
		series.MarkAreas = newMarkAreaFromCloudIssue(freq, issue)
	}

	c.MultiSeries = append(c.MultiSeries, series)
}

//nolint:mnd // Magic numbers are for style.
func (c *scatterChart) addSeries(name, color string, scatterData []opts.ScatterData) {
	series := echarts.SingleSeries{
		Name:       name,
		Type:       types.ChartScatter,
		Color:      color,
		Data:       scatterData,
		Symbol:     "rect",
		SymbolSize: 6,
	}
	c.MultiSeries = append(c.MultiSeries, series)
}

//nolint:mnd // Magic numbers are for style.
func newMarkLineAverage() *opts.MarkLines {
	return &opts.MarkLines{
		Data: []any{
			opts.MarkLineNameTypeItem{
				Name: "Average",
				Type: "average",
			},
		},
		MarkLineStyle: opts.MarkLineStyle{
			Symbol:     []string{"circle", "pin"},
			SymbolSize: 8,
			Label: &opts.Label{
				Formatter:     "avg:\n{c}",
				Align:         "left",
				VerticalAlign: "middle",
			},
		},
	}
}

//nolint:mnd // Magic numbers are for style.
func newMarkLineFromCloudIssue(issue *data.Issue) *opts.MarkLines {
	if issue.AlertData == nil || issue.AlertData.ProbeLatencyAvg == 0 {
		return nil
	}
	return &opts.MarkLines{
		Data: []any{
			opts.MarkLineNameYAxisItem{
				Name:  "Threshold",
				YAxis: issue.AlertData.ProbeLatencyAvg.Round(time.Millisecond).Milliseconds(),
			},
		},
		MarkLineStyle: opts.MarkLineStyle{
			Symbol:     []string{"pin"},
			SymbolSize: 8,
			LineStyle: &opts.LineStyle{
				Color:   "#f87171",
				Opacity: 0.7,
			},
			Label: &opts.Label{
				Formatter:     "x̄:\n{c}",
				Align:         "left",
				VerticalAlign: "middle",
			},
		},
	}
}

func newMarkAreaFromCloudIssue(freq time.Duration, issue *data.Issue) *opts.MarkAreas {
	if freq <= time.Minute {
		freq = 0 // Display real begin and end of the issue window
	}

	var timeEnd time.Time
	if issue.TimeEnd == nil {
		timeEnd = time.Now()
	} else {
		timeEnd = issue.TimeEnd.Add(-freq / markAreaIssueTimeMod)
	}
	timeBegin := issue.TimeBegin.Add(-freq / markAreaIssueTimeMod)

	return &opts.MarkAreas{
		Data: []any{
			[]opts.MarkAreaData{
				{
					XAxis: timeBegin,
				},
				{
					XAxis: timeEnd,
				},
			},
		},
	}
}

func appendDummyLineItem(lineData []opts.LineData, t time.Time) []opts.LineData {
	return append(lineData, opts.LineData{Value: []any{t, nil}})
}
