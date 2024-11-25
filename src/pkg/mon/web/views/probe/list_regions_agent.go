package probe

import (
	"context"
	"time"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5"

	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/pb"
)

const (
	sqlSelectAgentInfo = `SELECT
  mon_agent.id    			AS "agent_id",
  mon_region.name           AS "cloud_region",
  mon_region.location       AS "cloud_location",
  mon_agent.version         AS "version",
  NOW() - mon_agent.started AS "uptime",
  COALESCE(SUM(mon_job.errors), 0) AS "errors"
FROM {schema}.mon_agent                AS mon_agent
  LEFT JOIN {schema}.mon_config_region AS mon_region ON mon_region.name = mon_agent.cloud_region
  LEFT JOIN {schema}.mon_job           AS mon_job ON mon_job.agent_id = mon_agent.id
WHERE mon_agent.status = $1
  AND mon_job.time > $2
  AND mon_region.enabled = TRUE
GROUP BY mon_agent.id, mon_agent.version, mon_region.id, mon_region.name, mon_region.location
ORDER BY mon_region.id`
)

type agentInfo struct {
	AgentId       string
	CloudRegion   string
	CloudLocation string
	Version       string
	Uptime        time.Duration
	Errors        int
}

func (a *agentInfo) region() *data.CloudRegion {
	return &data.CloudRegion{
		Name:     a.CloudRegion,
		Location: a.CloudLocation,
	}
}

func (v *ListRegionsView) getAgentInfo(ctx context.Context, cloud string) ([]*agentInfo, error) {
	query := db.WithSchema(sqlSelectAgentInfo, cloud)
	errorsSince := time.Now().Add(-defaultTimeSpan)
	raws, _ := v.dbc.Read.Query(ctx, query, pb.AgentAction_AGENT_START, errorsSince)
	rows, err := pgx.CollectRows(raws, pgx.RowToAddrOfStructByName[agentInfo])
	if err != nil {
		//nolint:wrapcheck // Error is properly logged by the caller.
		return nil, err
	}
	return rows, nil
}

func (v *ListRegionsView) agentTempl(
	ctx context.Context,
	cloud *data.Cloud,
	probe *data.ProbeDefinition,
	issues *issueMap,
) (templ.Component, error) {
	agents, err := v.getAgentInfo(ctx, cloud.Id)
	if err != nil {
		return nil, err
	}
	//nolint:contextcheck // Confused linter.
	return listAgentRegionsTempl(
		cloud,
		probe,
		agents,
		issues,
	), nil
}
