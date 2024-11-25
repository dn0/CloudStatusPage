package pb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/jackc/pgx/v5"

	"cspage/pkg/db"
	"cspage/pkg/msg"
)

const (
	sqlSelectAgents = `SELECT
  mon_agent.id           AS "id",
  mon_agent.version      AS "version",
  mon_agent.hostname     AS "hostname",
  mon_agent.ip_address   AS "ip_address",
  mon_agent.cloud_region AS "cloud_region"
FROM {schema}.mon_agent AS mon_agent
WHERE mon_agent.status = $1`
	sqlInsertAgentStart = `INSERT INTO {schema}.mon_agent (
  id, status, started, version, hostname, ip_address, cloud_region, cloud_zone, sysinfo
 ) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9)`
	// Different stop statuses have different priorities, e.g. a manual stop can be overridden but not an auto stop.
	sqlUpdateAgentStop = `UPDATE {schema}.mon_agent SET status = $2, stopped = $3 WHERE id = $1 AND status <= $2`
)

var errUnknownAgentStatus = errors.New("unknown agent status")

type Agents []*Agent

//nolint:wrapcheck // Error is properly logged by the caller.
func (a Agents) Render(_ context.Context, w io.Writer) error {
	return json.NewEncoder(w).Encode(a)
}

func (a *Agent) Repr() string {
	return fmt.Sprintf("<Agent:%s>", a.ID())
}

func (a *Agent) ID() string {
	return a.GetId()
}

func (a *Agent) Save(ctx context.Context, dbc db.Client, attrs *msg.Attrs) error {
	var err error

	switch a.GetAction() {
	case AgentAction_AGENT_START:
		err = a.start(ctx, dbc, attrs.Cloud)
		if err == nil {
			if errr := a.consolidateRunningAgents(ctx, dbc, attrs.Cloud); errr != nil {
				slog.Error(
					"Could not consolidate running agents",
					"cloud", attrs.Cloud,
					"region", attrs.Region,
					"agent", a,
					"err", errr,
				)
			}
		}
	case AgentAction_AGENT_STOP, AgentAction_AGENT_STOPPING, AgentAction_AGENT_STOP_MANUAL:
		err = a.stop(ctx, dbc, attrs.Cloud)
	case AgentAction_AGENT_UNKNOWN:
		err = errUnknownAgentStatus
	}

	return err
}

//nolint:protogetter,wrapcheck // Running on a freshly marshaled object and DB error is wrapped by caller.
func (a *Agent) start(ctx context.Context, dbc db.Client, cloud string) error {
	_, err := dbc.Exec(ctx, db.WithSchema(sqlInsertAgentStart, cloud),
		a.Id,
		a.Action,
		a.Time.AsTime(),
		a.Version,
		a.Hostname,
		a.IpAddress,
		a.CloudRegion,
		a.CloudZone,
		a.Sysinfo,
	)
	return err
}

//nolint:protogetter,wrapcheck // Running on a freshly marshaled object and DB error is wrapped by caller.
func (a *Agent) stop(ctx context.Context, dbc db.Client, cloud string) error {
	_, err := dbc.Exec(ctx, db.WithSchema(sqlUpdateAgentStop, cloud),
		a.Id,
		a.Action,
		a.Time.AsTime(),
	)
	return err
}

//nolint:protogetter // Running on a freshly marshaled object.
func (a *Agent) consolidateRunningAgents(ctx context.Context, dbc db.Client, cloud string) error {
	region := a.GetCloudRegion()
	agents, err := GetRunningAgents(ctx, dbc, cloud, region)
	if err != nil {
		return err
	}

	var errs []error
	for _, agent := range agents {
		if agent.Id == a.Id {
			continue
		}
		slog.Warn("Stopping stale agent", "cloud", cloud, "region", region, "old_agent", agent, "new_agent", a)
		agent.Action = AgentAction_AGENT_STOP_MANUAL // Will be applied only if the agent's action==start in DB
		agent.Time = a.Time                          // Stale agent's stop time = new agent's start time
		if err = agent.stop(ctx, dbc, cloud); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func GetRunningAgents(ctx context.Context, dbc db.Client, cloud, region string) ([]*Agent, error) {
	query := db.WithSchema(sqlSelectAgents, cloud)
	args := []any{AgentAction_AGENT_START}
	if region != "" {
		query += " AND mon_agent.cloud_region = $2"
		args = append(args, region)
	}

	rows, _ := dbc.Query(ctx, query, args...)
	//nolint:wrapcheck // Error is properly logged by the caller.
	return pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[Agent])
}
