package analyst

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"cspage/pkg/db"
	"cspage/pkg/pb"
)

const (
	sqlSelectFirstPingSince = `SELECT mon_ping.time
FROM {schema}.mon_ping AS mon_ping
  LEFT JOIN {schema}.mon_job AS mon_job ON mon_job.id = mon_ping.job_id
  LEFT JOIN {schema}.mon_agent AS mon_agent ON mon_agent.id = mon_job.agent_id
WHERE mon_ping.time > $1
  AND mon_job.time > $1
  AND mon_agent.cloud_region = $2
ORDER BY mon_ping.time ASC LIMIT 1`
	sqlSelectMissingPings = `WITH
  agg AS (
    SELECT
      public.heartbeat_agg(mon_ping.time, @agg_start, @agg_duration, @heartbeat_liveness) AS health,
      mon_agent.cloud_region AS "cloud_region"
    FROM {schema}.mon_ping AS mon_ping
      LEFT JOIN {schema}.mon_job AS mon_job ON mon_job.id = mon_ping.job_id
      LEFT JOIN {schema}.mon_agent AS mon_agent ON mon_agent.id = mon_job.agent_id
    WHERE mon_ping.time >= @agg_start
      AND mon_job.time >= @agg_start
    GROUP BY mon_agent.cloud_region
  ),
  cte AS (
    SELECT
	  t.cloud_region AS "cloud_region",
	  t.start AS "begin",
	  t.end AS "end",
	  t.start - @heartbeat_liveness AS "time_begin",
	  t.end + @heartbeat_liveness AS "time_end"
    FROM (SELECT (dead_ranges(health)).*, cloud_region FROM agg) AS t
    WHERE (t.end - t.start) > @heartbeat_liveness
)
SELECT
  cte.*,
  mon_ping.job_id AS "job_id",
  mon_alert.id    AS "alert_id"
FROM cte
  LEFT JOIN {schema}.mon_ping  AS mon_ping  ON mon_ping.time = cte.time_begin
  LEFT JOIN {schema}.mon_alert AS mon_alert ON mon_alert.job_id = mon_ping.job_id
  LEFT JOIN {schema}.mon_job   AS mon_job   ON mon_job.id = mon_ping.job_id
  LEFT JOIN {schema}.mon_agent AS mon_agent ON mon_agent.id = mon_job.agent_id
WHERE mon_ping.time >= @agg_start
  AND mon_job.time >= @agg_start
  AND mon_agent.cloud_region = cte.cloud_region -- just to be sure
  AND mon_agent.status = @agent_running -- temporary as it happens during agent updates
  ORDER BY cte.time_begin`
)

// missingPing JobID and AgentID will be JobId and AgentId here to be consistent with ProtoBuf.
type missingPing struct {
	CloudRegion string    `json:"cloud_region"`
	Begin       time.Time `json:"begin"`
	End         time.Time `json:"end"`
	TimeBegin   time.Time `json:"time_begin"`
	TimeEnd     time.Time `json:"time_end"`
	JobId       *string   `json:"job_id"`
	AlertId     *string   `json:"alert_id"`
}

func (j *pingJob) getMissingPings(ctx context.Context, until time.Time) ([]*missingPing, error) {
	since := j.checkpoint.Add(-j.cfg.PingLagWindow)
	rows, _ := j.dbc.Read.Query(ctx, db.WithSchema(sqlSelectMissingPings, j.cloud), pgx.NamedArgs{
		"agg_start":          since,
		"agg_duration":       until.Sub(since).Round(time.Microsecond) + j.durationBuffer,
		"heartbeat_liveness": j.cfg.PingLagThreshold,
		"agent_running":      pb.AgentAction_AGENT_START,
	})
	//nolint:wrapcheck // Error is properly logged by the caller.
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[missingPing])
}

func (j *pingJob) getFirstPingSince(ctx context.Context, since time.Time, region string) (time.Time, error) {
	var value time.Time
	query := db.WithSchema(sqlSelectFirstPingSince, j.cloud)
	err := j.dbc.Write.QueryRow(ctx, query, since, region).Scan(&value)
	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, nil
	}
	//nolint:wrapcheck // Error is properly logged by the caller.
	return value, err
}
