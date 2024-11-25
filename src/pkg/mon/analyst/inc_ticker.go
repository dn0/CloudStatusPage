package analyst

import (
	"context"
	"fmt"
	"log/slog"

	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/msg"
	"cspage/pkg/worker"
)

func NewIncidentTicker(
	ctx context.Context,
	cfg *Config,
	dbc *db.Clients,
	cloud string,
	topic *msg.PubsubTopic,
	channel chan<- *data.Incident,
) *worker.Daemon {
	name := worker.TickerJobPrefix + "incident:" + cloud
	job := &incidentJob{
		analystJob: analystJob{
			cfg:           cfg,
			dbc:           dbc,
			publisher:     msg.NewPublisher(&cfg.PubsubPublisherConfig, topic),
			baseLog:       slog.With("job", name),
			cloud:         cloud,
			name:          name,
			incidentsChan: channel,
		},
	}

	return worker.NewTicker(ctx, &cfg.BaseConfig, job, cfg.IncidentCheckInterval, "sender", fmt.Sprintf("%T", channel))
}
