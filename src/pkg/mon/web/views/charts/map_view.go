package charts

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"

	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/mon/web/views"
	"cspage/pkg/pb"
)

const (
	mapName = "world"
)

//nolint:gochecknoglobals // This is a constant.
var (
	incidentSeverityColor = map[pb.IncidentSeverity]string{
		pb.IncidentSeverity_INCIDENT_NONE:   "green",
		pb.IncidentSeverity_INCIDENT_LOW:    "yellow",
		pb.IncidentSeverity_INCIDENT_MEDIUM: "orange",
		pb.IncidentSeverity_INCIDENT_HIGH:   "red",
	}
	regionEnabledSymbol = map[bool]string{
		true:  "&check;",
		false: "&times;",
	}
)

type MapView struct {
	dbc *db.Clients
}

func NewMapView(dbc *db.Clients) *MapView {
	return &MapView{
		dbc: dbc,
	}
}

func (v *MapView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	views.CacheControl(w, r, views.CacheMaxAgeDefault, nil)
	c, err := v.handler(r)
	views.Render(w, r, c, err)
}

//nolint:wrapcheck // Error is properly logged by the caller.
func (v *MapView) handler(r *http.Request) (templ.Component, error) {
	issueQuery := data.IssueQuery{
		Clouds: data.CloudIds,
		Status: pb.IncidentStatus_INCIDENT_OPEN,
	}
	issues, err := issueQuery.GetAll(r.Context(), v.dbc.Read)
	if err != nil {
		return nil, err
	}

	regions, err := data.GetCloudRegionsGeo(r.Context(), v.dbc.Read, data.CloudIds)
	if err != nil {
		return nil, err
	}

	return v.templ(r, issues, regions), nil
}

func (v *MapView) templ(r *http.Request, issues []*data.Issue, regions []*data.CloudRegionGeo) templ.Component {
	issueMap := make(map[string][]*data.Issue)
	for _, issue := range issues {
		for _, region := range issue.CloudRegions {
			key := issue.CloudId + "|" + region
			issueMap[key] = append(issueMap[key], issue)
		}
	}

	cloudRegions := make(map[string][]*geoData)
	cloudIssues := make([]*geoData, 0)
	for _, region := range regions {
		if region.Lat == nil || region.Lon == nil {
			if region.Enabled {
				slog.Error("Cloud region is missing geo coordinates", "cloud", region.CloudId, "region", region)
			}
			continue
		}
		pointRegion, pointIssue := newMapPoints(r.Context(), region, issueMap[region.CloudId+"|"+region.Name])
		cloudRegions[region.CloudId] = append(cloudRegions[region.CloudId], pointRegion)
		if pointIssue != nil {
			cloudIssues = append(cloudIssues, pointIssue)
		}
	}

	chart := newMapChart(mapName)
	for _, cloud := range data.Clouds {
		chart.addCloudRegionSeries(cloud.Name, cloud.Color, cloud.Symbol, cloudRegions[cloud.Id])
	}
	chart.addCloudIssuesSeries(cloudIssues)

	//nolint:contextcheck // Context is not needed here.
	return mapChartTempl(chart)
}

//nolint:mnd // Magic numbers are for style.
func newMapPoints(ctx context.Context, region *data.CloudRegionGeo, issues []*data.Issue) (*geoData, *geoData) {
	name := region.Cloud().Name + " " + region.Name
	desc := fmt.Sprintf("Location: %s<br>Enabled: %s", region.Location, regionEnabledSymbol[region.Enabled])

	if !region.Enabled {
		return newMapPoint(name, desc, "", region.Lat, region.Lon, nil), nil
	}

	link := region.Cloud().URLPrefix() + region.URLPrefix()

	if issues == nil {
		point := newMapPoint(name, desc, link, region.Lat, region.Lon, &pointStyle{
			BorderColor: incidentSeverityColor[pb.IncidentSeverity_INCIDENT_NONE],
			BorderWidth: 1.0,
			Opacity:     0.8,
		})
		return point, nil
	}

	var highestSev pb.IncidentSeverity
	issuesDesc := newIssueDescription()
	for _, i := range issues {
		sev := i.Severity()
		if sev > highestSev {
			highestSev = sev
		}
		_ = issuesDesc.add(ctx, i)
	}

	color := incidentSeverityColor[highestSev]
	point := newMapPoint(name, desc, link, region.Lat, region.Lon, &pointStyle{
		BorderColor: color, BorderWidth: 1.5, Opacity: 0.9,
	})
	desc += issuesDesc.String()
	issuePoint := newMapPoint(name, desc, data.OpenIssuesURL, region.Lat, region.Lon, &pointStyle{
		Color: color,
	})

	return point, issuePoint
}
