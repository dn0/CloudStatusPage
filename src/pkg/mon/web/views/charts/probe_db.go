package charts

import (
	"context"
	"time"

	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/jackc/pgx/v5"

	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/pb"
)

const (
	sqlSelectProbeResults = `SELECT
  mon_job.time      AS "time",
  mon_probe.action  AS "action",
  mon_probe.status  AS "status",
  ROUND(mon_probe.latency / 1000000::numeric, @rounding) AS "latency",
  mon_probe.error   AS "error"
FROM {schema}.mon_probe_{table} AS mon_probe
  LEFT JOIN {schema}.mon_job    AS mon_job   ON mon_job.id = mon_probe.job_id
  LEFT JOIN {schema}.mon_agent  AS mon_agent ON mon_agent.id = mon_job.agent_id
WHERE mon_probe.time > @since AND mon_probe.time <= @until
  AND mon_job.time > @since AND mon_job.time <= @until_job
  AND mon_agent.cloud_region = @region
  AND mon_probe.action = ANY(@actions)
ORDER BY mon_probe.time ASC
`
)

type probeResult struct {
	Time    time.Time
	Action  uint32
	Status  pb.ResultStatus
	Latency any // int64 or float64 depending on the @rounding
	Error   string
}

func (v *probeView) getProbeResults(
	ctx context.Context,
	cloud string,
	region *data.CloudRegion,
	probe *data.ProbeDefinition,
	form *probeViewForm,
	probeResults map[uint32]*probeAction,
) error {
	_, rounding := probe.LatencyRounding()
	rows, _ := v.dbc.Read.Query(ctx, db.WithSchemaAndTable(sqlSelectProbeResults, cloud, probe.Name), pgx.NamedArgs{
		"since":     form.Since,
		"until":     form.Until,
		"until_job": form.Until.Add(probe.Config.Timeout()),
		"region":    region.Name,
		"actions":   form.Actions,
		"rounding":  rounding,
	})
	results, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[probeResult])
	if err != nil {
		//nolint:wrapcheck // Error is properly logged by the caller.
		return err
	}

	for _, res := range results {
		pres := probeResults[res.Action]
		switch res.Status {
		case pb.ResultStatus_RESULT_UNKNOWN:
			continue
		case pb.ResultStatus_RESULT_FAILURE, pb.ResultStatus_RESULT_TIMEOUT:
			pres.errors = append(pres.errors, opts.ScatterData{
				Value: []any{
					res.Time,
					res.Latency,
					data.ParseProbeError(res.Error),
				},
			})
		case pb.ResultStatus_RESULT_SUCCESS:
			pres.series = append(pres.series, opts.LineData{Value: []any{
				res.Time,
				res.Latency,
			}})
		}
	}

	return nil
}
