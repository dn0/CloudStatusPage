package pb

import (
	"cspage/pkg/db"
	"cspage/pkg/msg"
)

const (
	sqlInsertJob = `INSERT INTO {schema}.mon_job (agent_id, id, time, drift, took, name, errors)
VALUES ($1, $2, $3, $4, $5, $6, $7)`
)

func (j *Job) ID() string {
	return j.GetId()
}

//nolint:protogetter,nonamedreturns // Running on a freshly marshaled Object.
func (j *Job) sqlInsert(attrs *msg.Attrs) (query string, args []any) {
	return db.WithSchema(sqlInsertJob, attrs.Cloud), []any{
		j.AgentId,
		j.Id,
		j.Time.AsTime(),
		j.Drift.AsDuration(),
		j.Took.AsDuration(),
		j.Name,
		j.Errors,
	}
}
