package speaker

import (
	"context"
	"log/slog"

	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/worker"
)

const (
	maxPendingIncidents = 8
)

// NewIncidentChannel create a new channel which the analyst sends new or updated Incidents to.
func NewIncidentChannel() chan *data.Incident {
	return make(chan *data.Incident, maxPendingIncidents)
}

// NewIncidentReceiver creates a channel receiver that is running inside analyst.
func NewIncidentReceiver(
	ctx context.Context,
	cfg *Config,
	dbc *db.Clients,
	cloud string,
	channel <-chan *data.Incident,
) *worker.Daemon {
	name := worker.ReceiverJobPrefix + "incident:" + cloud
	job := &incidentJob{
		cfg:     cfg,
		dbc:     dbc,
		log:     slog.With("job", name),
		cloud:   cloud,
		name:    name,
		twitter: newTwitterClient(ctx, cfg, cloud),
	}

	return worker.NewReceiver(ctx, cfg.BaseConfig, job, channel)
}
