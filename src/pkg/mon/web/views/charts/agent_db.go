package charts

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/jackc/pgx/v5"

	"cspage/pkg/data"
	"cspage/pkg/db"
)

const (
	sqlSelectJobStats = `, mon_job.time AS "time"
FROM {schema}.mon_job AS mon_job
  LEFT JOIN {schema}.mon_agent  AS mon_agent ON mon_agent.id = mon_job.agent_id
WHERE mon_job.time > @since AND mon_job.time <= @until
  AND mon_job.name = @probe_name
  AND mon_agent.cloud_region = @region
`
	sqlSelectAgentStats = `, mon_job.time AS "time"
FROM {schema}.mon_ping AS mon_ping
  LEFT JOIN {schema}.mon_job    AS mon_job   ON mon_job.id = mon_ping.job_id
  LEFT JOIN {schema}.mon_agent  AS mon_agent ON mon_agent.id = mon_job.agent_id
WHERE mon_ping.time > @since AND mon_ping.time <= @until
  AND mon_job.time > @since AND mon_job.time <= @until
  AND mon_agent.cloud_region = @region
ORDER BY mon_ping.time ASC
`
)

type chartData struct {
	metrics []*chartMetric
	series  map[string][]opts.LineData
}

func (v *agentView) getData(
	ctx context.Context,
	cloud string,
	region *data.CloudRegion,
	probe *data.ProbeDefinition,
	form *agentViewForm,
	baseQuery string,
	metrics []*chartMetric,
) (map[string][]opts.LineData, error) {
	cdata := &chartData{
		metrics: metrics,
		series:  make(map[string][]opts.LineData),
	}

	var fields []string
	for _, metric := range metrics {
		for _, field := range metric.fields {
			cdata.series[field.alias] = []opts.LineData{}
			fields = append(fields, fmt.Sprintf(`%s AS %q`, field.column, field.alias))
		}
	}

	query := "SELECT " + strings.Join(fields, ", ") + baseQuery
	rows, _ := v.dbc.Read.Query(ctx, db.WithSchemaAndTable(query, cloud, probe.Name), pgx.NamedArgs{
		"since":      form.Since,
		"until":      form.Until,
		"region":     region.Name,
		"probe_name": probe.Name,
	})
	_, err := pgx.CollectRows(rows, cdata.rowToMetrics)
	if err != nil {
		//nolint:wrapcheck // Error is properly logged by the caller.
		return nil, err
	}

	return cdata.series, nil
}

func (d *chartData) rowToMetrics(row pgx.CollectableRow) (map[string]any, error) {
	dataRow, err := pgx.RowToMap(row)
	if err != nil {
		//nolint:wrapcheck // Error is properly logged by the caller.
		return nil, err
	}

	for _, metric := range d.metrics {
		for _, field := range metric.fields {
			//nolint:forcetypeassert // This is OK.
			d.series[field.alias] = append(d.series[field.alias], opts.LineData{Value: []any{
				dataRow["time"].(time.Time),
				dataRow[field.alias],
			}})
		}
	}

	return dataRow, nil
}
