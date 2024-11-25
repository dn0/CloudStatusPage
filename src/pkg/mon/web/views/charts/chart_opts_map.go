package charts

import (
	echarts "github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/event"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

const (
	mapChartClickEvent       = "echartMapClickEventListener"
	mapChartTooltipFormatter = "echartMapTooltipFormatter"
)

//nolint:gochecknoglobals // This is a constant.
var chartErrorColors = []string{
	"red",
	"#ff0266",
	"#b00020",
	"#e23636",
}

//nolint:mnd // Magic numbers are for style.
func newMapChart(typ string) *mapChart {
	chart := &mapChart{Geo: *echarts.NewGeo()}
	chart.Initialization.Theme = "auto" // To avoid automatic injections of "light" color scheme
	chart.GeoComponent.Map = typ
	chart.GeoComponent.Silent = opts.Bool(true)
	// Our extension of GeoComponent
	chart.geoConfig.Roam = opts.Bool(true)
	chart.geoConfig.Zoom = 1.1
	// -
	chart.Legend.Bottom = "70rem"
	chart.Legend.ItemStyle = &opts.ItemStyle{
		Opacity:     0.9,
		BorderWidth: 0.5,
	}
	chart.Tooltip.Enterable = opts.Bool(true)
	chart.Tooltip.Formatter = opts.FuncOpts(mapChartTooltipFormatter)
	chart.Toolbox.Top = "64rem"
	chart.Toolbox.Feature = &opts.ToolBoxFeature{
		Restore: &opts.ToolBoxFeatureRestore{},
	}
	chart.EventListeners = []event.Listener{{
		EventName: "click",
		Handler:   mapChartClickEvent,
	}}
	return chart
}

func newMapPoint(name, desc, link string, lat, lon *float32, style *pointStyle) *geoData {
	return &geoData{
		GeoData: opts.GeoData{
			Name:  name,
			Value: []any{lon, lat, desc, link},
		},
		ItemStyle: style,
	}
}

//nolint:mnd // Magic numbers are for style.
func (c *mapChart) addCloudRegionSeries(name, color, symbol string, dataSeries []*geoData) {
	series := echarts.SingleSeries{
		Name:        name,
		Type:        types.ChartScatter,
		Color:       color,
		ColorBy:     "data",
		Symbol:      symbol,
		SymbolSize:  10,
		Data:        dataSeries,
		CoordSystem: types.ChartGeo,
		// Style of disabled regions
		ItemStyle: &opts.ItemStyle{
			BorderColor: "black",
			BorderWidth: 1.0,
			Opacity:     0.1,
		},
	}
	c.MultiSeries = append(c.MultiSeries, series)
}

//nolint:mnd // Magic numbers are for style.
func (c *mapChart) addCloudIssuesSeries(dataSeries []*geoData) {
	series := echarts.SingleSeries{
		Name:        "Issues",
		Type:        types.ChartEffectScatter,
		Color:       "#ff0266",
		ColorBy:     "data",
		Symbol:      "circle",
		SymbolSize:  10,
		Data:        dataSeries,
		CoordSystem: types.ChartGeo,
		// Style of disabled regions
		ItemStyle: &opts.ItemStyle{
			Opacity: 0.8,
		},
		RippleEffect: &opts.RippleEffect{
			Period:    3.5,
			Scale:     4,
			BrushType: "fill",
		},
	}
	c.MultiSeries = append(c.MultiSeries, series)
}
