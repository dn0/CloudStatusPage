package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"spock/pkg/db"
	"spock/pkg/pb"
)

const (
	sqlSelectJob = `SELECT
  mon_job.agent_id       AS "agent_id",
  mon_job.id             AS "id",
  mon_job.time           AS "time",
  mon_job.drift          AS "drift",
  mon_job.took           AS "took",
  mon_job.name           AS "name",
  mon_job.errors         AS "errors",
  mon_agent.cloud_region AS "cloud_region"
FROM {schema}.mon_job AS mon_job
  LEFT JOIN {schema}.mon_agent AS mon_agent ON mon_agent.id = mon_job.agent_id
WHERE mon_job.id = $1`
	sqlSelectLatestJob = `SELECT
  mon_job.agent_id       AS "agent_id",
  mon_job.id             AS "id",
  mon_job.time           AS "time",
  mon_job.drift          AS "drift",
  mon_job.took           AS "took",
  mon_job.name           AS "name",
  mon_job.errors         AS "errors",
  mon_agent.cloud_region AS "cloud_region"
FROM {schema}.mon_job AS mon_job
  LEFT JOIN {schema}.mon_agent AS mon_agent ON mon_agent.id = mon_job.agent_id
WHERE mon_job.name = $1
  AND mon_agent.cloud_region = $2
  AND mon_job.time >= NOW() - INTERVAL '14 days'
ORDER BY mon_job.time DESC
LIMIT 1`
)

// Job ID will be Id here for consistency with ProtoBuf.
type Job struct {
	AgentId     string        `json:"agent_id"`
	Id          string        `json:"id"`
	Time        time.Time     `json:"time"`
	Drift       time.Duration `json:"drift"`
	Took        time.Duration `json:"took"`
	Name        string        `json:"name"`
	Errors      uint32        `json:"errors"`
	CloudRegion string        `json:"cloud_region"`
}

func (j *Job) IsPing() bool {
	return j.Name == pb.JobNamePing
}

func (j *Job) toProto() *pb.Job {
	return &pb.Job{
		AgentId: j.AgentId,
		Id:      j.Id,
		Time:    timestamppb.New(j.Time),
		Drift:   durationpb.New(j.Drift),
		Took:    durationpb.New(j.Took),
		Name:    j.Name,
		Errors:  j.Errors,
	}
}

func getJob(ctx context.Context, dbc db.Conn, cloud, jobId string) (*Job, error) {
	rows, _ := dbc.Query(ctx, db.WithSchema(sqlSelectJob, cloud), jobId)
	//nolint:wrapcheck // Error is properly logged by the caller.
	return pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[Job])
}

func GetLastJob(ctx context.Context, dbc db.Conn, cloud, name, region string) (*Job, error) {
	rows, _ := dbc.Query(ctx, db.WithSchema(sqlSelectLatestJob, cloud), name, region)
	obj, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[Job])
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return nil, &db.ObjectNotFoundError{Object: fmt.Sprintf("job=%s cloud=%s region=%s", name, cloud, region)}
	}
	//nolint:wrapcheck // Error is properly logged by the caller.
	return obj, err
}
