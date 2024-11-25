package analyst

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/pb"
)

const (
	sqlSelectProbeResults = `WITH
 agg AS (
  SELECT
    mon_probe.time         AS "time",
    mon_probe.job_id       AS "job_id",
    mon_job.time           AS "job_time",
    mon_agent.cloud_region AS "cloud_region",
    mon_probe.action       AS "action",
    mon_probe.status       AS "status",
    mon_probe.error        AS "error",
    mon_probe.latency      AS "latency",
    public.stats_agg(mon_probe.latency) OVER w AS "latency_stats"
  FROM {schema}.mon_probe_{table} AS mon_probe
    LEFT JOIN {schema}.mon_job AS mon_job ON mon_job.id = mon_probe.job_id
    LEFT JOIN {schema}.mon_agent AS mon_agent ON mon_agent.id = mon_job.agent_id
  WHERE mon_probe.time > @since_window AND mon_probe.time <= @until
    AND mon_job.time > @since_window AND mon_job.time <= @until
    AND mon_probe.action <= @max_action_id
  WINDOW w AS (PARTITION BY mon_agent.cloud_region, mon_probe.action, mon_probe.status
               ORDER BY mon_probe.time ASC RANGE @window PRECEDING)
 ),
 cte AS (
  SELECT
    agg.time,
    agg.job_id,
    agg.job_time,
    agg.cloud_region,
    agg.action,
    agg.status,
    agg.error,
    agg.latency,
    COALESCE(NULLIF(public.average(agg.latency_stats), 'NaN')::bigint, 0) AS "latency_avg",
    COALESCE(NULLIF(public.stddev(agg.latency_stats),  'NaN')::bigint, 0) AS "latency_sd",
    COALESCE(mon_alert.id::text, '') AS "alert_id"
  FROM agg
    LEFT JOIN {schema}.mon_alert AS mon_alert ON mon_alert.job_id = agg.job_id
  WHERE agg.time > @since
  ORDER BY agg.time
)
SELECT cte.*,
  cte.latency-cte.latency_avg AS "distance",
  COALESCE(((cte.latency-cte.latency_avg)/NULLIF(cte.latency_sd,0)::float), 0) AS "zscore"
FROM cte`
	sqlSelectJobAgentUptime = `SELECT $1 - mon_agent.started AS "uptime"
FROM {schema}.mon_agent AS mon_agent
  LEFT JOIN {schema}.mon_job AS mon_job ON mon_job.agent_id = mon_agent.id
WHERE mon_job.id = $2`
)

// probe JobID will be JobId here for consistency with ProtoBuf.
type probeResult struct {
	Time        time.Time
	JobId       string
	JobTime     time.Time
	CloudRegion string
	Action      uint32
	Status      pb.ResultStatus
	Error       string
	Latency     time.Duration
	LatencyAvg  time.Duration
	LatencySD   time.Duration
	Distance    time.Duration
	Zscore      float32
	AlertId     string
}

func (j *probeJob) getProbeResults(ctx context.Context, until time.Time) ([]*probeResult, error) {
	rows, _ := j.dbc.Read.Query(ctx, db.WithSchemaAndTable(sqlSelectProbeResults, j.cloud, j.probe.Name), pgx.NamedArgs{
		"window":        j.zscoreWindow,
		"since_window":  j.checkpoint.Add(-j.zscoreWindow),
		"until":         until,
		"since":         j.checkpoint.Add(-j.checkWindow),
		"max_action_id": data.ProbeMaxDisplayActionId,
	})
	//nolint:wrapcheck // Error is properly logged by the caller.
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[probeResult])
}

func (j *probeJob) getJobAgentUptime(ctx context.Context, until time.Time, jobId string) (time.Duration, error) {
	var value time.Duration
	// 'until' is the same until as above, i.e. the current checkpoint
	err := j.dbc.Write.QueryRow(ctx, db.WithSchema(sqlSelectJobAgentUptime, j.cloud), until, jobId).Scan(&value)
	//nolint:wrapcheck // Error is properly logged by the caller.
	return value, err
}
