package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/timestamppb"

	"cspage/pkg/db"
	"cspage/pkg/msg"
	"cspage/pkg/pb"
)

const (
	sqlInsertAlert = `INSERT INTO {schema}.mon_alert (
  id, job_id, created, updated, time_begin, time_end, time_check, type, status, cloud_region, probe_name, probe_action, data
 ) VALUES (
  $1, $2, $3, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	sqlUpdateAlert = `UPDATE {schema}.mon_alert SET updated=$2, time_end=$3, time_check=$4, status=$5 WHERE id = $1`
	sqlSelectAlert = `SELECT
  mon_alert.id AS "id",
  mon_alert.job_id AS "job_id",
  COALESCE(mon_alert.incident_id::TEXT, '') AS "incident_id",
  mon_alert.created AS "created",
  mon_alert.updated AS "updated",
  mon_alert.time_begin AS "time_begin",
  mon_alert.time_end AS "time_end",
  mon_alert.time_check AS "time_check",
  mon_alert.type AS "type",
  mon_alert.status AS "status",
  mon_alert.cloud_region AS "cloud_region",
  mon_alert.probe_name AS "probe_name",
  mon_alert.probe_action AS "probe_action",
  mon_alert.data as "data"
FROM {schema}.mon_alert AS mon_alert`
)

type AlertData struct {
	Note               string        `json:"note,omitempty"`
	Trigger            string        `json:"trigger"`
	ServiceName        string        `json:"service_name"`
	ServiceGroup       string        `json:"service_group"`
	ProbeDescription   string        `json:"probe_description"`
	ProbeActionName    string        `json:"probe_action_name"`
	ProbeActionTitle   string        `json:"probe_action_title"`
	ProbeError         string        `json:"probe_error,omitempty"`
	ProbeLatencyValue  time.Duration `json:"probe_latency_value,omitempty"`
	ProbeLatencyAvg    time.Duration `json:"probe_latency_avg,omitempty"`
	ProbeLatencySD     time.Duration `json:"probe_latency_stddev,omitempty"`
	ProbeLatencyZscore float32       `json:"probe_latency_zscore,omitempty"`
}

// Alert ID will be Id here for consistency with ProtoBuf.
type Alert struct {
	job *Job

	Id          string         `json:"id"`
	JobId       string         `json:"job_id"`
	IncidentId  string         `json:"incident_id"`
	Created     time.Time      `json:"created"`
	Updated     time.Time      `json:"updated"`
	TimeBegin   time.Time      `json:"time_begin"`
	TimeEnd     *time.Time     `json:"time_end"`
	TimeCheck   time.Time      `json:"time_check"`
	Type        pb.AlertType   `json:"type"`
	Status      pb.AlertStatus `json:"status"`
	CloudRegion string         `json:"cloud_region"`
	ProbeName   string         `json:"probe_name"`
	ProbeAction uint32         `json:"probe_action"`
	Data        *AlertData     `json:"data"`
}

func (ad *AlertData) ProbeLatencyThreshold(zscore float64) time.Duration {
	return ad.ProbeLatencyAvg + time.Duration(float64(ad.ProbeLatencySD)*zscore)
}

func (a *Alert) affectedService() *AffectedService {
	return &AffectedService{
		CloudRegion:     a.CloudRegion,
		Name:            a.Data.ServiceName,
		Group:           a.Data.ServiceGroup,
		ProbeName:       a.ProbeName,
		ProbeActionName: a.Data.ProbeActionName,
	}
}

func (a *Alert) ToMessage(ctx context.Context, dbc db.Conn, cloud string) (*msg.Message, error) {
	pbMessage, err := a.toProto(ctx, dbc, cloud)
	if err != nil {
		return nil, err
	}
	attrs := msg.NewAttrs(msg.TypeAlert, cloud, a.CloudRegion)
	return msg.NewMessage(pbMessage, attrs), nil
}

func (a *Alert) toProto(ctx context.Context, dbc db.Conn, cloud string) (*pb.Alert, error) {
	if a.job == nil {
		var err error
		if a.job, err = getJob(ctx, dbc, cloud, a.JobId); err != nil {
			return nil, fmt.Errorf("%s: failed to fetch alert's job id=%q: %w", a.Id, a.JobId, err)
		}
	}

	var timeEnd *timestamppb.Timestamp
	if a.TimeEnd != nil {
		timeEnd = timestamppb.New(*a.TimeEnd)
	}

	dataJSON, err := json.Marshal(a.Data)
	if err != nil {
		return nil, fmt.Errorf("%s: failed marshal alert's data: %w", a.Id, err)
	}

	return &pb.Alert{
		Id:          a.Id,
		Job:         a.job.toProto(),
		IncidentId:  a.IncidentId,
		Created:     timestamppb.New(a.Created),
		Updated:     timestamppb.New(a.Updated),
		TimeBegin:   timestamppb.New(a.TimeBegin),
		TimeEnd:     timeEnd,
		TimeCheck:   timestamppb.New(a.TimeCheck),
		Type:        a.Type,
		Status:      a.Status,
		CloudRegion: a.CloudRegion,
		ProbeName:   a.ProbeName,
		ProbeAction: a.ProbeAction,
		Data:        dataJSON,
	}, nil
}

//nolint:varnamelen // Using a (alert) is OK in this context.
func CreateAlert(ctx context.Context, dbc db.Conn, cloud string, a *Alert) error {
	if a.JobId == "" {
		return &db.FieldIsRequiredError{Field: "JobId"}
	}

	a.Id = uuid.Must(uuid.NewV7()).String()
	_, err := dbc.Exec(ctx, db.WithSchema(sqlInsertAlert, cloud),
		a.Id,
		a.JobId,
		a.Created,
		a.TimeBegin,
		a.TimeEnd,
		a.TimeCheck,
		a.Type,
		a.Status,
		a.CloudRegion,
		a.ProbeName,
		a.ProbeAction,
		a.Data,
	)
	//nolint:wrapcheck // Error is properly logged by the caller.
	return err
}

//nolint:varnamelen // Using a (alert) is OK in this context.
func UpdateAlert(ctx context.Context, dbc db.Conn, cloud string, a *Alert) error {
	_, err := dbc.Exec(ctx, db.WithSchema(sqlUpdateAlert, cloud),
		a.Id,
		a.Updated,
		a.TimeEnd,
		a.TimeCheck,
		a.Status,
	)
	//nolint:wrapcheck // Error is properly logged by the caller.
	return err
}

func GetAlerts(
	ctx context.Context,
	dbc db.Conn,
	cloud string,
	s []pb.AlertStatus,
	filters map[string]any,
) ([]*Alert, error) {
	query := db.WithSchema(sqlSelectAlert, cloud) + "\nWHERE mon_alert.status = ANY($1)"
	args := []any{s}
	for k, v := range filters {
		args = append(args, v)
		query += " AND " + k
	}
	rows, _ := dbc.Query(ctx, query, args...)
	//nolint:wrapcheck // Error is properly logged by the caller.
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Alert])
}

func AlertTypeToIncidentSeverity(t pb.AlertType) pb.IncidentSeverity {
	switch t {
	case pb.AlertType_PING_MISSING:
		return pb.IncidentSeverity_INCIDENT_MEDIUM
	case pb.AlertType_PROBE_FAILURE, pb.AlertType_PROBE_TIMEOUT:
		return pb.IncidentSeverity_INCIDENT_MEDIUM
	case pb.AlertType_PROBE_SLOW:
		return pb.IncidentSeverity_INCIDENT_LOW
	default:
		return pb.IncidentSeverity_INCIDENT_NONE
	}
}
