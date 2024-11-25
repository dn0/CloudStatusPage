package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"

	"cspage/pkg/db"
	"cspage/pkg/msg"
	"cspage/pkg/pb"
)

const (
	incidentOutageThreshold = 3

	sqlInsertIncident = `INSERT INTO {schema}.mon_incident (
  id, created, updated, time_begin, time_end, severity, status, cloud_regions, data
 ) VALUES (
  $1, $2, $2, $3, $4, $5, $6, $7, $8)`
	sqlUpdateIncident = `UPDATE {schema}.mon_incident
SET updated=$2, time_begin=$3, time_end=$4, severity=$5, status=$6, cloud_regions=$7, data=$8
WHERE id=$1`
	sqlUpdateIncidentData  = `UPDATE {schema}.mon_incident SET data=$2 WHERE id=$1`
	sqlLinkAlertToIncident = `UPDATE {schema}.mon_alert SET updated=$2, incident_id=$3 WHERE id = $1`
	sqlSelectIncident      = `SELECT
  mon_incident.id AS "id",
  mon_incident.created AS "created",
  mon_incident.updated AS "updated",
  mon_incident.time_begin AS "time_begin",
  mon_incident.time_end AS "time_end",
  mon_incident.severity AS "severity",
  mon_incident.status AS "status",
  mon_incident.cloud_regions AS "cloud_regions",
  mon_incident.data as "data"
FROM {schema}.mon_incident AS mon_incident`
)

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type IncidentLink struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type IncidentData struct {
	Note      string            `json:"note,omitempty"`
	Services  []AffectedService `json:"services,omitempty"`
	Summary   string            `json:"summary"`
	Links     []IncidentLink    `json:"links,omitempty"`
	TwitterId string            `json:"twitter_id,omitempty"`
}

//goland:noinspection GoUnnecessarilyExportedIdentifiers
type Incident struct {
	Id           string              `json:"id"`
	Created      time.Time           `json:"created"`
	Updated      time.Time           `json:"updated"`
	TimeBegin    time.Time           `json:"time_begin"`
	TimeEnd      *time.Time          `json:"time_end"`
	Severity     pb.IncidentSeverity `json:"severity"`
	Status       pb.IncidentStatus   `json:"status"`
	CloudRegions []string            `json:"cloud_regions"`
	Data         *IncidentData       `json:"data"`

	alerts []*Alert
}

func (inc *Incident) DetailsURL(cloud string) string {
	return "/cloud/" + cloud + incidentURLPrefix + "/" + inc.Id
}

func (inc *Incident) ToMessage(ctx context.Context, dbc db.Conn, cloud string) (*msg.Message, error) {
	pbMessage, err := inc.toProto(ctx, dbc, cloud)
	if err != nil {
		return nil, err
	}
	attrs := msg.NewAttrs(msg.TypeIncident, cloud, strings.Join(inc.CloudRegions, ","))
	return msg.NewMessage(pbMessage, attrs), nil
}

func (inc *Incident) toProto(ctx context.Context, dbc db.Conn, cloud string) (*pb.Incident, error) {
	alerts := make([]*pb.Alert, len(inc.alerts))
	for j, alert := range inc.alerts {
		var err error
		if alerts[j], err = alert.toProto(ctx, dbc, cloud); err != nil {
			return nil, fmt.Errorf("%s: failed to fetch incident's alert id=%s: %w", inc.Id, alert.Id, err)
		}
	}

	regions := make([]*pb.CloudRegion, len(inc.CloudRegions))
	for j, region := range inc.CloudRegions {
		regions[j] = &pb.CloudRegion{Region: region}
	}

	var timeEnd *timestamppb.Timestamp
	if inc.TimeEnd != nil {
		timeEnd = timestamppb.New(*inc.TimeEnd)
	}

	dataJSON, err := json.Marshal(inc.Data)
	if err != nil {
		return nil, fmt.Errorf("%s: failed marshal incidents's data: %w", inc.Id, err)
	}

	return &pb.Incident{
		Id:           inc.Id,
		Created:      timestamppb.New(inc.Created),
		Updated:      timestamppb.New(inc.Updated),
		TimeBegin:    timestamppb.New(inc.TimeBegin),
		TimeEnd:      timeEnd,
		Severity:     inc.Severity,
		Status:       inc.Status,
		CloudRegions: regions,
		Alerts:       alerts,
		Data:         dataJSON,
	}, nil
}

//nolint:cyclop // This has to be complex.
func (inc *Incident) makeSummary(cloud string, slowOnly bool) {
	region := inc.CloudRegions[0]
	if len(inc.CloudRegions) > 1 {
		region = "multiple regions"
	}

	problem := "problem"
	if slowOnly {
		problem = "performance degradation"
	}

	if len(inc.Data.Services) == 1 {
		inc.Data.Summary = fmt.Sprintf("%s %s in %s", inc.Data.Services[0].Name, problem, region)
		return
	}

	// Multiple services, but maybe they share the same service group (e.g. compute)
	serviceName := inc.Data.Services[0].Name
	serviceGroup := inc.Data.Services[0].Group
	for _, svc := range inc.Data.Services {
		if serviceName != svc.Name {
			serviceName = ""
		}
		if serviceGroup != svc.Group {
			serviceGroup = ""
			break
		}
	}

	switch {
	case serviceName != "":
		inc.Data.Summary = fmt.Sprintf("%s %s in %s", serviceName, problem, region)
	case serviceGroup != "":
		inc.Data.Summary = fmt.Sprintf("%s %s in %s", serviceGroup, problem, region)
	default:
		switch {
		case slowOnly:
			problem = "performance degradations"
		case len(inc.Data.Services) >= incidentOutageThreshold:
			problem = "outage"
		default:
			problem = "problems"
		}
		inc.Data.Summary = fmt.Sprintf("%s %s in %s", CloudMap[cloud].Name, problem, region)
	}
}

//nolint:wrapcheck,cyclop,gocognit,funlen // Errors are properly logged by the caller. Should be complex.
func (inc *Incident) CreateOrUpdateFromAlerts(
	ctx context.Context,
	dbc db.Conn,
	cloud string,
	checkpoint time.Time,
	alerts []*ExtendedAlert,
) error {
	newInc := false
	sql := sqlUpdateIncident
	inc.Status = pb.IncidentStatus_INCIDENT_CLOSED
	inc.alerts = nil

	if inc.Id == "" {
		inc.Id = uuid.Must(uuid.NewV7()).String()
		newInc = true
		sql = sqlInsertIncident
	}

	if inc.Data == nil {
		inc.Data = &IncidentData{}
	}

	dtx, err := dbc.Begin(ctx)
	if err != nil {
		return err
	}

	cloudRegions := make(map[string]struct{})
	affectedServices := make(map[string]*AffectedService)
	slowOnly := true

	for _, alert := range alerts {
		inc.alerts = append(inc.alerts, &alert.Alert) // For protobuf notifications

		cloudRegions[alert.CloudRegion] = struct{}{}
		service := alert.affectedService()
		affectedServices[service.Hash()] = service

		if alert.Type != pb.AlertType_PROBE_SLOW {
			slowOnly = false
		}

		alertSeverity := AlertTypeToIncidentSeverity(alert.Type)
		if alertSeverity > inc.Severity {
			inc.Severity = alertSeverity
		}

		if inc.Created.IsZero() || alert.Created.Before(inc.Created) {
			inc.Created = alert.Created
		}

		if inc.TimeBegin.IsZero() || alert.TimeBegin.Before(inc.TimeBegin) {
			inc.TimeBegin = alert.TimeBegin
		}

		if alert.TimeEnd != nil && (inc.TimeEnd == nil || alert.TimeEnd.After(*inc.TimeEnd)) {
			inc.TimeEnd = alert.TimeEnd
		}

		if alert.Status == pb.AlertStatus_ALERT_OPEN {
			inc.Status = pb.IncidentStatus_INCIDENT_OPEN
		}

		if newInc || alert.IncidentId != inc.Id {
			if _, err = dtx.Exec(ctx, db.WithSchema(sqlLinkAlertToIncident, cloud),
				alert.Id,
				inc.Created,
				inc.Id,
			); err != nil {
				return err
			}
		}
	}

	inc.CloudRegions = nil
	for region := range cloudRegions {
		inc.CloudRegions = append(inc.CloudRegions, region)
	}

	inc.Data.Services = nil
	for _, service := range affectedServices {
		inc.Data.Services = append(inc.Data.Services, *service)
	}

	if inc.Status == pb.IncidentStatus_INCIDENT_OPEN {
		inc.TimeEnd = nil
	}

	if newInc {
		inc.Updated = inc.Created
	} else {
		inc.Updated = checkpoint
	}

	inc.makeSummary(cloud, slowOnly)

	_, err = dtx.Exec(ctx, db.WithSchema(sql, cloud),
		inc.Id,
		inc.Updated,
		inc.TimeBegin,
		inc.TimeEnd,
		inc.Severity,
		inc.Status,
		inc.CloudRegions,
		inc.Data,
	)
	if err != nil {
		return err
	}

	return dtx.Commit(ctx)
}

func UpdateIncidentData(ctx context.Context, dbc db.Conn, cloud string, inc *Incident) error {
	_, err := dbc.Exec(ctx, db.WithSchema(sqlUpdateIncidentData, cloud),
		inc.Id,
		inc.Data,
	)
	//nolint:wrapcheck // Error is properly logged by the caller.
	return err
}

func GetIncidents(
	ctx context.Context,
	dbc db.Conn,
	cloud string,
	s []pb.IncidentStatus,
	filters map[string]any,
) ([]*Incident, error) {
	query := db.WithSchema(sqlSelectIncident, cloud) + "\nWHERE mon_incident.status = ANY($1)"
	args := []any{s}
	for k, v := range filters {
		args = append(args, v)
		query += " AND " + k
	}
	rows, _ := dbc.Query(ctx, query, args...)
	//nolint:wrapcheck // Error is properly logged by the caller.
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Incident])
}

func GetIncident(
	ctx context.Context,
	dbc db.Conn,
	cloud string,
	id string,
) (*Incident, error) {
	query := db.WithSchema(sqlSelectIncident, cloud) + "\nWHERE mon_incident.id = $1"
	rows, _ := dbc.Query(ctx, query, id)
	obj, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[Incident])
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, &db.ObjectNotFoundError{Object: fmt.Sprintf("incident=%s cloud=%s", id, cloud)}
	}
	return obj, nil
}
