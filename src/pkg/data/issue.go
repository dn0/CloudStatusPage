package data

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"

	"cspage/pkg/db"
	"cspage/pkg/pb"
)

//goland:noinspection GoUnnecessarilyExportedIdentifiers
const (
	IssueTypeAll      IssueType = 0
	IssueTypeAlert    IssueType = 1
	IssueTypeIncident IssueType = 2

	IssuesURLSuffix    = "/issues"
	AlertsURLSuffix    = "/alerts"
	IncidentsURLSuffix = "/incidents"
	alertURLPrefix     = "/alert"
	incidentURLPrefix  = "/incident"
	openURLSuffix      = "/open"
	closedURLSuffix    = "/closed"
	OpenIssuesURL      = IssuesURLSuffix + openURLSuffix

	agentLinkName = "Monitoring agent details"
	probeLinkName = "Monitoring probe details"

	chartTimeBufferMin         = 120 * time.Minute //nolint:revive // This is minimum not minutes.
	chartTimeBufferPingMin     = 30 * time.Minute  //nolint:revive // This is minimum not minutes.
	chartTimeBufferDurationDiv = 12 * time.Hour
)

//nolint:gochecknoglobals // These are constants.
var (
	IssueStatusURLSuffix = map[pb.IncidentStatus]string{
		pb.IncidentStatus_INCIDENT_ANY:    "",
		pb.IncidentStatus_INCIDENT_OPEN:   openURLSuffix,
		pb.IncidentStatus_INCIDENT_CLOSED: closedURLSuffix,
	}
	IssueStatusURLTitle = map[pb.IncidentStatus]string{
		pb.IncidentStatus_INCIDENT_ANY:    "All",
		pb.IncidentStatus_INCIDENT_OPEN:   "Open",
		pb.IncidentStatus_INCIDENT_CLOSED: "Closed",
	}
)

type IssueType uint8

type Issue struct {
	Incident
	Type    IssueType
	CloudId string

	AlertType        pb.AlertType
	AlertProbeName   string
	AlertProbeAction uint32
	AlertData        *AlertData
	AlertIncidentId  string
}

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type AffectedService struct {
	CloudRegion     string `json:"cloud_region"`
	Name            string `json:"name"`
	Group           string `json:"group"`
	ProbeName       string `json:"probe_name"`
	ProbeActionName string `json:"probe_action_name"`
}

func (s *AffectedService) Hash() string {
	return s.CloudRegion + "|" + s.ProbeActionName
}

func (i *Issue) Cloud() *Cloud {
	return CloudMap[i.CloudId]
}

func (i *Issue) Regions() []*CloudRegion {
	regions := make([]*CloudRegion, len(i.CloudRegions))
	for idx, region := range i.CloudRegions {
		regions[idx] = &CloudRegion{Name: region}
	}
	return regions
}

func (i *Issue) Duration(now time.Time) time.Duration {
	if i.TimeEnd == nil {
		return now.Sub(i.TimeBegin)
	}
	return i.TimeEnd.Sub(i.TimeBegin)
}

func (i *Issue) Severity() pb.IncidentSeverity {
	if i.Type == IssueTypeAlert {
		return AlertTypeToIncidentSeverity(i.AlertType)
	}
	return i.Incident.Severity
}

func (i *Issue) Status() pb.IncidentStatus {
	if i.Type == IssueTypeAlert {
		switch pb.AlertStatus(i.Incident.Status) {
		case pb.AlertStatus_ALERT_OPEN:
			return pb.IncidentStatus_INCIDENT_OPEN
		case pb.AlertStatus_ALERT_CLOSED_AUTO, pb.AlertStatus_ALERT_CLOSED_MANUAL:
			return pb.IncidentStatus_INCIDENT_CLOSED
		case pb.AlertStatus_ALERT_UNKNOWN:
			return pb.IncidentStatus_INCIDENT_ANY
		}
	}
	return i.Incident.Status
}

func (i *Issue) AffectedServices() []AffectedService {
	if i.Type == IssueTypeAlert {
		return []AffectedService{{
			CloudRegion:     i.CloudRegions[0],
			Name:            i.AlertData.ServiceName,
			Group:           i.AlertData.ServiceGroup,
			ProbeName:       i.AlertProbeName,
			ProbeActionName: i.AlertData.ProbeActionName,
		}}
	}
	return i.Incident.Data.Services
}

func (i *Issue) Summary() string {
	if i.Type == IssueTypeAlert {
		var prefix string
		if i.AlertType == pb.AlertType_PING_MISSING {
			// Cloud Monitoring Agent: PING_MISSING
			prefix = i.Cloud().Name + " " + i.AlertData.ProbeDescription
		} else {
			// Operation Cloud Probe Name: PROBE_ISSUE
			prefix = ProbeActionFullName(
				i.AlertData.ProbeDescription,
				i.AlertData.ProbeActionName,
				i.AlertData.ProbeActionTitle,
				false,
			)
		}
		return prefix + ": " + i.AlertType.String()
	}
	return i.Incident.Data.Summary
}

func (i *Issue) ListURL() string {
	var suffix string
	if i.Type == IssueTypeAlert {
		suffix = AlertsURLSuffix
	} else {
		suffix = IncidentsURLSuffix
	}
	return "/cloud/" + i.CloudId + suffix
}

func (i *Issue) DetailsURL() string {
	var prefix string
	if i.Type == IssueTypeAlert {
		prefix = alertURLPrefix
	} else {
		prefix = incidentURLPrefix
	}
	return "/cloud/" + i.CloudId + prefix + "/" + i.Id
}

func (i *Issue) Links() []IncidentLink {
	if i.Type == IssueTypeAlert {
		link := IncidentLink{
			Name: probeLinkName,
			URL: fmt.Sprintf(
				"/cloud/%s/region/%s/probe/%s",
				i.CloudId,
				i.CloudRegions[0],
				i.AlertProbeName,
			),
		}
		if i.AlertType == pb.AlertType_PING_MISSING {
			link.Name = agentLinkName
		}
		return []IncidentLink{link}
	}
	return i.Incident.Data.Links
}

func (i *Issue) ChartURL() string {
	if i.Type != IssueTypeAlert {
		return ""
	}
	qs := make(url.Values)
	qs.Set("alert", i.Id)
	if i.AlertType == pb.AlertType_PING_MISSING {
		qs.Set("metric", "job_duration")
	} else {
		qs.Set("action", strconv.Itoa(int(i.AlertProbeAction)))
	}
	buffer := chartTimeBufferMin
	if i.AlertType == pb.AlertType_PING_MISSING {
		buffer = chartTimeBufferPingMin
	}
	buffer += time.Duration((float64(i.Duration(time.Now())) / float64(chartTimeBufferDurationDiv)) * float64(buffer))
	qs.Set("since", i.TimeBegin.Add(-buffer).In(time.UTC).Format(time.RFC3339))
	if i.TimeEnd != nil && i.Status() == pb.IncidentStatus_INCIDENT_CLOSED {
		qs.Set("until", i.TimeEnd.Add(buffer).In(time.UTC).Format(time.RFC3339))
	}
	return fmt.Sprintf(
		"/cloud/%s/region/%s/probe/%s/charts?%s",
		i.CloudId,
		i.CloudRegions[0],
		i.AlertProbeName,
		qs.Encode(),
	)
}

func GetAlertIssue(ctx context.Context, dbc db.Client, cloud, id string) (*Issue, error) {
	query := db.WithSchema(sqlSelectAlertIssue, cloud) + " WHERE mon_alert.id = $1"
	return getIssue(ctx, dbc, query, cloud, id)
}

func GetIncidentIssue(ctx context.Context, dbc db.Client, cloud, id string) (*Issue, error) {
	query := db.WithSchema(sqlSelectIncidentIssue, cloud) + " WHERE mon_incident.id = $1"
	return getIssue(ctx, dbc, query, cloud, id)
}

func getIssue(ctx context.Context, dbc db.Client, query, cloud, id string) (*Issue, error) {
	rows, _ := dbc.Query(ctx, query, id)
	obj, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[Issue])
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, &db.ObjectNotFoundError{Object: fmt.Sprintf("issue=%s cloud=%s", id, cloud)}
	}
	return obj, nil
}
