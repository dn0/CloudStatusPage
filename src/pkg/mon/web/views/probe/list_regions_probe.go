package probe

import (
	"context"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5"

	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/pb"
)

const (
	sqlSelectProbeStats = `SELECT
  t1.cloud_region    AS "cloud_region",
  t1.action          AS "action",
  COALESCE(NULLIF(public.average(latency_stats), 'NaN')::bigint, 0) AS "latency_avg",
  COALESCE(NULLIF(public.stddev(latency_stats), 'NaN')::bigint, 0) AS "latency_sd",
  COALESCE(t2.errors, 0) AS "errors"
FROM (
  SELECT
    mon_agent.cloud_region              AS "cloud_region",
    mon_probe.action                    AS "action",
	public.stats_agg(mon_probe.latency) AS "latency_stats"
  FROM {schema}.mon_probe_{table}        AS mon_probe
    LEFT JOIN {schema}.mon_job           AS mon_job    ON mon_job.id = mon_probe.job_id
    LEFT JOIN {schema}.mon_agent         AS mon_agent  ON mon_agent.id = mon_job.agent_id
    LEFT JOIN {schema}.mon_config_region AS mon_region ON mon_region.name = mon_agent.cloud_region
  WHERE mon_probe.time > $1
	AND mon_job.time > $1
    AND mon_region.enabled = TRUE
    AND mon_probe.status = ANY($2)
    AND mon_probe.action <= $4
  GROUP BY mon_agent.cloud_region, mon_probe.action
) AS t1
LEFT JOIN (
  SELECT
    mon_agent.cloud_region AS "cloud_region",
    mon_probe.action       AS "action",
    COUNT(*)               AS "errors"
  FROM {schema}.mon_probe_{table}        AS mon_probe
    LEFT JOIN {schema}.mon_job           AS mon_job    ON mon_job.id = mon_probe.job_id
    LEFT JOIN {schema}.mon_agent         AS mon_agent  ON mon_agent.id = mon_job.agent_id
    LEFT JOIN {schema}.mon_config_region AS mon_region ON mon_region.name = mon_agent.cloud_region
  WHERE mon_probe.time > $1
	AND mon_job.time > $1
    AND mon_region.enabled = TRUE
    AND mon_probe.status = ANY($3)
    AND mon_probe.action <= $4
  GROUP BY mon_agent.cloud_region, mon_probe.action
) AS t2
ON  t1.cloud_region = t2.cloud_region
AND t1.action = t2.action`
)

type probeStatsRow struct {
	CloudRegion string
	Action      uint32
	LatencyAvg  time.Duration
	LatencySD   time.Duration
	Errors      int
}

type probeStats struct {
	m map[string]*probeStatsRow
}

func probeStatsKey(cloudRegion string, probeAction uint32) string {
	return cloudRegion + "|" + strconv.Itoa(int(probeAction))
}

func (p *probeStats) get(cloudRegion string, probeAction uint32) *probeStatsRow {
	if res, ok := p.m[probeStatsKey(cloudRegion, probeAction)]; ok {
		return res
	}
	return &probeStatsRow{}
}

func (p *probeStats) isFastest(cloudRegion string, probeAction uint32) bool {
	if res, ok := p.m[probeStatsKey("fastest", probeAction)]; ok {
		return res.CloudRegion == cloudRegion
	}
	return false
}

func (p *probeStats) isSlowest(cloudRegion string, probeAction uint32) bool {
	if res, ok := p.m[probeStatsKey("slowest", probeAction)]; ok {
		return res.CloudRegion == cloudRegion
	}
	return false
}

func (v *ListRegionsView) getProbeStats(
	ctx context.Context,
	cloud string,
	probe *data.ProbeDefinition,
) (*probeStats, error) {
	query := db.WithSchemaAndTable(sqlSelectProbeStats, cloud, probe.Name)
	since := time.Now().Add(-defaultTimeSpan)
	raws, _ := v.dbc.Read.Query(ctx, query, since, pb.ResultStatus_OK, pb.ResultStatus_ERROR, data.ProbeMaxDisplayActionId)
	rows, err := pgx.CollectRows(raws, pgx.RowToAddrOfStructByName[probeStatsRow])
	if err != nil {
		//nolint:wrapcheck // Error is properly logged by the caller.
		return nil, err
	}

	stats := &probeStats{m: make(map[string]*probeStatsRow)}
	for _, row := range rows {
		stats.m[probeStatsKey(row.CloudRegion, row.Action)] = row
		fastestKey := probeStatsKey("fastest", row.Action)
		if fastest, ok := stats.m[fastestKey]; ok {
			if row.LatencyAvg < fastest.LatencyAvg {
				stats.m[fastestKey] = row
			}
		} else {
			stats.m[fastestKey] = row
		}
		slowestKey := probeStatsKey("slowest", row.Action)
		if slowest, ok := stats.m[slowestKey]; ok {
			if row.LatencyAvg > slowest.LatencyAvg {
				stats.m[slowestKey] = row
			}
		} else {
			stats.m[slowestKey] = row
		}
	}

	return stats, nil
}

func (v *ListRegionsView) probeTempl(
	ctx context.Context,
	cloud *data.Cloud,
	probe *data.ProbeDefinition,
	issues *issueMap,
) (templ.Component, error) {
	regions, err := data.GetCloudRegions(ctx, v.dbc.Read, cloud.Id)
	if err != nil {
		//nolint:wrapcheck // Error is properly logged by the caller.
		return nil, err
	}

	stats, err := v.getProbeStats(ctx, cloud.Id, probe)
	if err != nil {
		return nil, err
	}
	//nolint:contextcheck // Confused linter.
	return listProbeRegionsTempl(
		cloud,
		regions,
		probe,
		issues,
		stats,
	), nil
}
