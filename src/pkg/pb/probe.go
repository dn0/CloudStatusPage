package pb

import (
	"context"
	"fmt"

	"cspage/pkg/db"
	"cspage/pkg/msg"
)

const (
	sqlInsertProbe = `INSERT INTO {schema}.mon_probe_{table} (
  job_id, time, action, status, latency, error
 ) VALUES (
  $1, $2, $3, $4, $5, $6)`
)

func (p *Probe) Repr() string {
	return fmt.Sprintf("<Probe:%s:%s>", p.GetJob().GetName(), p.ID())
}

func (p *Probe) ID() string {
	return p.GetJob().ID()
}

//nolint:protogetter,wrapcheck // Running on a freshly marshaled Object and DB error is wrapped by caller.
func (p *Probe) Save(ctx context.Context, dbc db.Client, attrs *msg.Attrs) error {
	sqlJob, sqlJobArgs := p.Job.sqlInsert(attrs)
	sqlProbe := db.WithSchemaAndTable(sqlInsertProbe, attrs.Cloud, p.Job.Name)
	batch := &db.Batch{}
	batch.Queue(sqlJob, sqlJobArgs...)
	for _, result := range p.Result {
		if result.Status == ResultStatus_RESULT_UNKNOWN {
			continue
		}
		batch.Queue(
			sqlProbe,
			p.Job.Id,
			result.Time.AsTime(),
			result.Action,
			result.Status,
			result.Latency.AsDuration(),
			result.Error,
		)
	}
	return dbc.SendBatch(ctx, batch).Close()
}
