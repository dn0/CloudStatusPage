package charts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	echarts "github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/event"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

type pointStyle = opts.ItemStyle

// Adding some missing geo features.
type geoData struct {
	opts.GeoData
	//nolint:tagliatelle // This is required by the library.
	ItemStyle *pointStyle `json:"itemStyle,omitempty"`
}

type geoConfig struct {
	Roam types.Bool `json:"roam,omitempty"`
	Zoom float32    `json:"zoom,omitempty"`
}

type geoComponent struct {
	opts.GeoComponent
	geoConfig
}

//nolint:govet // There are some issues in echarts which we don't control.
type lineChart struct {
	echarts.Line
}

//nolint:govet // There are some issues in echarts which we don't control.
type scatterChart struct {
	echarts.Scatter
}

//nolint:govet // There are some issues in echarts which we don't control.
type mapChart struct {
	echarts.Geo
	geoConfig
}

type chartType interface {
	*lineChart | *mapChart

	ID() string
	Config() (string, error)
	EventFunctions() []event.Listener
}

func (c *lineChart) ID() string {
	return c.Initialization.ChartID
}

func (c *lineChart) Config() (string, error) {
	c.Validate()
	return jsonNotEscaped(c.JSON())
}

func (c *lineChart) EventFunctions() []event.Listener {
	return c.EventListeners
}

func (c *mapChart) ID() string {
	return c.Initialization.ChartID
}

func (c *mapChart) Config() (string, error) {
	c.Validate()
	obj := c.JSON()
	obj["geo"] = geoComponent{
		GeoComponent: c.GeoComponent,
		geoConfig:    c.geoConfig,
	}
	return jsonNotEscaped(obj)
}

func (c *mapChart) EventFunctions() []event.Listener {
	return c.EventListeners
}

// Function jsonNotEscaped works like JSON(), but it returns a marshaled object whose
// characters will not be escaped in the template. Copied from charts/base.go.
func jsonNotEscaped(obj map[string]any) (string, error) {
	buff := bytes.NewBufferString("")
	enc := json.NewEncoder(buff)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(obj); err != nil {
		return "", fmt.Errorf("failed to create chart config: %w", err)
	}
	return strings.NewReplacer(`__f__"`, "", `"__f__`, "", `__f__)`, "").Replace(buff.String()), nil
}
