package probe

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/pb"
)

const (
	sqlSelectAgentDetails = `SELECT
  mon_agent.id    		    AS "id",
  mon_agent.status          AS "status",
  NOW() - mon_agent.started AS "uptime",
  mon_agent.version         AS "version",
  mon_agent.cloud_region    AS "cloud_region",
  mon_agent.cloud_zone      AS "cloud_zone"
FROM {schema}.mon_agent AS mon_agent
WHERE mon_agent.id = $1`
	sqlSelectAgentErrors = `SELECT COALESCE(SUM(mon_job.errors), 0) FROM {schema}.mon_job AS mon_job
WHERE mon_job.agent_id = $1 AND mon_job.time > $2`
	sqlSelectProbeErrors = `SELECT COALESCE(SUM(mon_job.errors), 0) FROM {schema}.mon_job AS mon_job
LEFT JOIN {schema}.mon_agent AS mon_agent ON mon_agent.id = mon_job.agent_id
WHERE mon_job.name = $1 AND mon_job.time > $2 AND mon_agent.cloud_region = $3`
)

type agentDetails struct {
	Id          string
	Status      pb.AgentAction
	Uptime      time.Duration
	Version     string
	CloudRegion string
	CloudZone   string

	errors int
}

//nolint:wrapcheck // Error is properly logged by the caller.
func (v *DetailsView) getAgentDetails(ctx context.Context, cloud string, job *data.Job) (*agentDetails, error) {
	since := time.Now().Add(-defaultTimeSpan)
	if !job.IsPing() {
		obj := &agentDetails{}
		query := db.WithSchema(sqlSelectProbeErrors, cloud)
		err := v.dbc.Read.QueryRow(ctx, query, job.Name, since, job.CloudRegion).Scan(&obj.errors)
		return obj, err
	}
	rows, _ := v.dbc.Read.Query(ctx, db.WithSchema(sqlSelectAgentDetails, cloud), job.AgentId)
	obj, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[agentDetails])
	if err == nil {
		err = v.dbc.Read.QueryRow(ctx, db.WithSchema(sqlSelectAgentErrors, cloud), job.AgentId, since).Scan(&obj.errors)
	}
	return obj, err
}
