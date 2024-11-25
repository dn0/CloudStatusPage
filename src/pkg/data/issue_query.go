package data

import (
	"context"
	"maps"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"cspage/pkg/db"
	"cspage/pkg/pb"
)

type IssueCounter struct {
	Alerts    int
	Incidents int
}

type IssueQuery struct {
	*Paginator
	Counter IssueCounter

	Since       time.Time
	Until       time.Time
	Clouds      []string
	Status      pb.IncidentStatus
	CloudRegion string
	ProbeName   string
	IncidentId  string // GetAlerts will return all child alerts and GetIncidents of course only one incident
}

func (c *IssueCounter) Total() int {
	return c.Alerts + c.Incidents
}

func (c *IssueCounter) Get(typ IssueType) int {
	switch typ {
	case IssueTypeAll:
		return c.Alerts + c.Incidents
	case IssueTypeAlert:
		return c.Alerts
	case IssueTypeIncident:
		return c.Incidents
	}
	return -1
}

func (c *IssueCounter) Inc(typ IssueType) {
	//nolint:exhaustive // Only actual issue types make sense here.
	switch typ {
	case IssueTypeAlert:
		c.Alerts++
	case IssueTypeIncident:
		c.Incidents++
	}
}

func (q *IssueQuery) GetAlerts(ctx context.Context, dbc db.Client) ([]*Issue, error) {
	if err := q.activatePagination(ctx, dbc, IssueTypeAlert); err != nil {
		return nil, err
	}
	queries, args := q.queryAlerts(sqlSelectAlertIssue)
	return q.getIssues(ctx, dbc, queries, args)
}

func (q *IssueQuery) GetIncidents(ctx context.Context, dbc db.Client) ([]*Issue, error) {
	if err := q.activatePagination(ctx, dbc, IssueTypeIncident); err != nil {
		return nil, err
	}
	queries, args := q.queryIncidents(sqlSelectIncidentIssue)
	return q.getIssues(ctx, dbc, queries, args)
}

func (q *IssueQuery) GetAll(ctx context.Context, dbc db.Client) ([]*Issue, error) {
	if err := q.activatePagination(ctx, dbc, IssueTypeAll); err != nil {
		return nil, err
	}
	queries, args := q.queryAlerts(sqlSelectAlertIssue)
	queries2, args2 := q.queryIncidents(sqlSelectIncidentIssue)
	maps.Copy(args, args2)
	queries = append(queries, queries2...)
	return q.getIssues(ctx, dbc, queries, args)
}

func (q *IssueQuery) CountAll(ctx context.Context, dbc db.Client) (IssueCounter, error) {
	queries, args := q.queryAlerts(sqlCountAlertIssues)
	queries2, args2 := q.queryIncidents(sqlCountIncidentIssues)
	maps.Copy(args, args2)
	queries = append(queries, queries2...)
	return q.countIssues(ctx, dbc, queries, args)
}

func (q *IssueQuery) getIssues(ctx context.Context, dbc db.Client, queries []string, args pgx.NamedArgs) ([]*Issue, error) {
	query := "WITH issues AS (\n" + strings.Join(queries, "\nUNION ALL\n") +
		")\nSELECT * FROM issues ORDER BY created DESC"
	if q.Paginator != nil {
		limit, offset := q.getLimitOffset()
		query += " LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)
	}
	rows, _ := dbc.Query(ctx, query, args)
	//nolint:wrapcheck // Error is properly logged by the caller.
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Issue])
}

func (q *IssueQuery) countIssues(
	ctx context.Context,
	dbc db.Client,
	queries []string,
	args pgx.NamedArgs,
) (IssueCounter, error) {
	query := "WITH issues AS (\n" + strings.Join(queries, "\nUNION\n") +
		")\nSELECT SUM(num_alerts) AS alerts, SUM(num_incidents) AS incidents FROM issues"
	rows, _ := dbc.Query(ctx, query, args)
	//nolint:wrapcheck // Error is properly logged by the caller.
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[IssueCounter])
}

func (q *IssueQuery) activatePagination(ctx context.Context, dbc db.Client, typ IssueType) error {
	if q.Paginator == nil {
		return nil
	}
	if q.Paginator.NumPages != 0 { // Paginator is already active
		return nil
	}
	var err error
	q.Counter, err = q.CountAll(ctx, dbc)
	if err != nil {
		return err
	}
	return q.Paginator.SetCount(q.Counter.Get(typ))
}

func (q *IssueQuery) queryAlerts(baseQuery string) ([]string, pgx.NamedArgs) {
	var cond []string
	args := pgx.NamedArgs{}

	if q.IncidentId != "" {
		cond = append(cond, "mon_alert.incident_id = @alert_incident_id")
		args["alert_incident_id"] = q.IncidentId
	}

	if q.CloudRegion != "" {
		cond = append(cond, "mon_alert.cloud_region = @region")
		args["region"] = q.CloudRegion
	}

	switch q.Status {
	case pb.IncidentStatus_INCIDENT_OPEN:
		cond = append(cond, "mon_alert.status = ANY(@alert_status)")
		args["alert_status"] = pb.AlertStatus_OPEN
	case pb.IncidentStatus_INCIDENT_CLOSED:
		cond = append(cond, "mon_alert.status = ANY(@alert_status)")
		args["alert_status"] = pb.AlertStatus_CLOSED
	case pb.IncidentStatus_INCIDENT_ANY:
		cond = append(cond, "mon_alert.status = ANY(@alert_status)")
		args["alert_status"] = pb.AlertStatus_ALL
	}

	if q.ProbeName != "" {
		cond = append(cond, "mon_alert.probe_name = @probe_name")
		args["probe_name"] = q.ProbeName
	}

	return q.makeIssueQueries(sqlTableAlert, baseQuery, cond, args)
}

func (q *IssueQuery) queryIncidents(baseQuery string) ([]string, pgx.NamedArgs) {
	var cond []string
	args := pgx.NamedArgs{}

	if q.IncidentId != "" {
		cond = append(cond, "mon_incident.id = @incident_id")
		args["incident_id"] = q.IncidentId
	}

	if q.CloudRegion != "" {
		cond = append(cond, "@region = ANY(mon_incident.cloud_regions)")
		args["region"] = q.CloudRegion
	}

	switch q.Status {
	case pb.IncidentStatus_INCIDENT_OPEN:
		cond = append(cond, "mon_incident.status = @incident_status")
		args["incident_status"] = pb.IncidentStatus_INCIDENT_OPEN
	case pb.IncidentStatus_INCIDENT_CLOSED:
		cond = append(cond, "mon_incident.status = @incident_status")
		args["incident_status"] = pb.IncidentStatus_INCIDENT_CLOSED
	case pb.IncidentStatus_INCIDENT_ANY:
		cond = append(cond, "mon_incident.status = ANY(@incident_status)")
		args["incident_status"] = []pb.IncidentStatus{pb.IncidentStatus_INCIDENT_OPEN, pb.IncidentStatus_INCIDENT_CLOSED}
	}

	if q.ProbeName != "" {
		cond = append(cond, "mon_alert.probe_name = @incident_probe_name")
		args["incident_probe_name"] = q.ProbeName
		baseQuery += sqlSelectIncidentJoinAlerts
	}

	return q.makeIssueQueries(sqlTableIncident, baseQuery, cond, args)
}

func (q *IssueQuery) makeIssueQueries(
	table string,
	baseQuery string,
	where []string,
	args pgx.NamedArgs,
) ([]string, pgx.NamedArgs) {
	if q.Since.IsZero() {
		q.Since = time.Now().Add(-sqlSelectIssuesDefaultTimeSpan)
	}
	where = append(where, table+".created > @since")
	args["since"] = q.Since

	if !q.Until.IsZero() {
		where = append(where, table+".created <= @until")
		args["until"] = q.Until
	}

	where = append(where, sqlSelectIssuesBaseCondition)
	baseQuery += " WHERE " + strings.Join(where, " AND ")

	queries := make([]string, len(q.Clouds))
	for i, cloud := range q.Clouds {
		query := "( " + db.WithSchema(baseQuery, cloud) + " )"
		queries[i] = query
	}

	return queries, args
}
