package analyst

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/pb"
)

const (
	maxUnsavedCheckpoints = 240
	checkpointFallback    = 1 * time.Hour
	sqlUpsertCheckpoint   = `INSERT INTO {schema}.mon_checkpoint VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET time = $2`
	sqlSelectCheckpoint   = `SELECT time FROM {schema}.mon_checkpoint WHERE name = $1`
)

func (j *analystJob) saveCheckpoint(ctx context.Context, value time.Time) error {
	if j.cfg.DryRun {
		return nil
	}
	_, err := j.dbc.Write.Exec(ctx, db.WithSchema(sqlUpsertCheckpoint, j.cloud), j.String(), value)
	//nolint:wrapcheck // Error is properly logged by the caller.
	return err
}

func (j *analystJob) loadCheckpoint(ctx context.Context) (time.Time, error) {
	var value time.Time
	err := j.dbc.Write.QueryRow(ctx, db.WithSchema(sqlSelectCheckpoint, j.cloud), j.String()).Scan(&value)

	if errors.Is(err, pgx.ErrNoRows) {
		fallback := time.Now().Add(-checkpointFallback).Round(time.Microsecond)
		j.baseLog.Warn("Checkpoint does not exist", "fallback", -checkpointFallback)
		return fallback, nil
	}

	//nolint:wrapcheck // Error is properly logged by the caller.
	return value, err
}

func (j *analystJob) getOpenAlerts(ctx context.Context, t []pb.AlertType, pname string) (map[string]*data.Alert, error) {
	filters := map[string]any{"mon_alert.type = ANY($2)": t}
	alerts, err := data.GetAlerts(ctx, j.dbc.Read, j.cloud, pb.AlertStatus_OPEN, filters)
	if err != nil {
		//nolint:wrapcheck // Error is properly logged by the caller.
		return nil, err
	}

	alertMap := map[string]*data.Alert{}
	for _, a := range alerts {
		if a.ProbeName == pname { // Include only alerts that match our probe name or 'ping'
			alertMap[a.Id] = a
		}
	}

	return alertMap, nil
}

func (j *analystJob) createAlert(ctx context.Context, alert *data.Alert) error {
	alert.Created = time.Now()
	//nolint:wrapcheck // Error is properly logged by the caller.
	return data.CreateAlert(ctx, j.dbtx, j.cloud, alert)
}

func (j *analystJob) updateAlert(ctx context.Context, alert *data.Alert, major bool) error {
	if major {
		// Otherwise it's only a checkpoint update.
		alert.Updated = time.Now()
	}
	//nolint:wrapcheck // Error is properly logged by the caller.
	return data.UpdateAlert(ctx, j.dbtx, j.cloud, alert)
}
