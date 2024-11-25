package analyst

import (
	"context"
	"errors"
	"log/slog"
	"math/rand"
	"time"

	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/msg"
	"cspage/pkg/worker"
)

const (
	minStartDelayMilliseconds = 100
	maxStartDelayMilliseconds = 5000
)

type taskCounter struct {
	Total int `json:"total"`
	Error int `json:"error"`
}

type analyzeFun = func(context.Context, time.Time) (time.Time, *data.Notification, bool)

type analystJob struct {
	cfg        *Config
	dbc        *db.Clients
	dbtx       db.Tx
	publisher  msg.Publisher
	baseLog    *slog.Logger
	taskLog    *slog.Logger
	cloud      string
	name       string
	probe      *data.ProbeDefinition
	checkpoint time.Time // Previous checkpoint
	tasks      taskCounter

	alertsChan    chan<- *data.Alert
	incidentsChan chan<- *data.Incident
}

func (j *analystJob) String() string {
	return j.name
}

func (j *analystJob) PreStart(_ context.Context) {}

func (j *analystJob) Start(ctx context.Context) {
	var err error
	if j.checkpoint, err = j.loadCheckpoint(ctx); err != nil {
		DieLog(j.baseLog, "Could not load checkpoint", "err", err)
	}

	j.baseLog.Info("Worker checkpoint loaded", "db_checkpoint", j.checkpoint)
	//nolint:gosec // So that not all cloud checks start at the same time.
	time.Sleep(time.Duration(minStartDelayMilliseconds+rand.Intn(maxStartDelayMilliseconds)) * time.Millisecond)
}

func (j *analystJob) Stop(ctx context.Context) {
	j.saveCheckpointNow(ctx, slog.LevelInfo)
	if j.alertsChan != nil {
		close(j.alertsChan)
	}
	if j.incidentsChan != nil {
		close(j.incidentsChan)
	}
}

func (j *analystJob) Shutdown(_ error) {}

func (j *analystJob) do(ctx context.Context, tick worker.Tick, analyze analyzeFun) {
	j.tasks.Total++
	now := tick.Time.Round(time.Microsecond)
	j.taskLog = j.baseLog.With(slog.Group("checkpoint", "previous", j.checkpoint, "current", now))

	var err error
	j.dbtx, err = j.dbc.Write.Begin(ctx)
	if err != nil {
		j.taskLog.Error("Could not begin transaction", "err", err)
		j.tasks.Error++
		return
	}

	//nolint:contextcheck // Rollback should not depend on the task context.
	defer j.rollback()

	//nolint:nestif // Nested ifs make sense here.
	if checkpoint, notification, ok := analyze(ctx, now); ok {
		if j.commit(ctx) {
			j.checkpoint = checkpoint
			if notification != nil {
				j.notify(ctx, notification)
			}
			if j.tasks.Total%maxUnsavedCheckpoints == 0 {
				j.saveCheckpointNow(ctx, slog.LevelDebug)
			}
			return // success
		}
	}
	j.tasks.Error++
}

func (j *analystJob) saveCheckpointNow(ctx context.Context, level slog.Level) {
	if err := j.saveCheckpoint(ctx, j.checkpoint); err == nil {
		logCtx := context.Background()
		//nolint:contextcheck // Logging context should be independent of the job context.
		j.baseLog.Log(logCtx, level, "Worker checkpoint saved", "db_checkpoint", j.checkpoint, "tasks", j.tasks)
	} else {
		j.baseLog.Error("Could not save checkpoint", "err", err)
	}
}

func (j *analystJob) commit(ctx context.Context) bool {
	if j.cfg.DryRun {
		return true
	}
	if err := j.dbtx.Commit(ctx); err != nil {
		j.taskLog.Error("Could not commit transaction", "err", err)
		return false
	}
	return true
}

func (j *analystJob) rollback() {
	if err := j.dbtx.Rollback(context.Background()); err != nil && !errors.Is(err, db.ErrTxClosed) {
		j.taskLog.Error("Could not rollback transaction", "err", err)
	}
}

func (j *analystJob) notify(ctx context.Context, notification *data.Notification) {
	for _, alert := range notification.Alerts {
		if alert.Id == "" {
			j.taskLog.Error("Cannot publish alert without an ID", "alert", alert)
		} else {
			j.publishAlert(ctx, alert)
		}
	}
	for _, inc := range notification.Incidents {
		if inc.Id == "" {
			j.taskLog.Error("Cannot publish incident without an ID", "incident", inc)
		} else {
			j.publishIncident(ctx, inc)
		}
	}
}

func (j *analystJob) publishAlert(ctx context.Context, alert *data.Alert) {
	if j.alertsChan != nil {
		j.alertsChan <- alert
	}
	message, err := alert.ToMessage(ctx, j.dbc.Read, j.cloud)
	if err != nil {
		j.taskLog.Error("Could not convert alert to proto message", "alert", alert, "err", err)
		return
	}
	//nolint:contextcheck // This means publish & forget as we are not waiting for the result
	//                       therefore we can't use the task context here.
	j.publisher.Publish(context.Background(), message)
}

func (j *analystJob) publishIncident(ctx context.Context, inc *data.Incident) {
	if j.incidentsChan != nil {
		j.incidentsChan <- inc
	}
	message, err := inc.ToMessage(ctx, j.dbc.Read, j.cloud)
	if err != nil {
		j.taskLog.Error("Could not convert incident to proto message", "incident", inc, "err", err)
		return
	}
	//nolint:contextcheck // This means publish & forget as we are not waiting for the result
	//                       therefore we can't use the task context here.
	j.publisher.Publish(context.Background(), message)
}
