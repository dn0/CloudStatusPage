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

const (
	missingPingDurationBufferDiv = 2
	missingPingAlertTrigger      = "nodata(agent.ping) > %.0fs"
)

func NewPingTicker(
	ctx context.Context,
	cfg *Config,
	dbc *db.Clients,
	cloud string,
	topic *msg.PubsubTopic,
) *worker.Daemon {
	probe := data.NewPingProbeDefinition()
	name := worker.TickerJobPrefix + probe.Name + ":" + cloud
	job := &pingJob{
		analystJob: analystJob{
			cfg:       cfg,
			dbc:       dbc,
			publisher: msg.NewPublisher(&cfg.PubsubPublisherConfig, topic),
			baseLog:   slog.With("job", name),
			cloud:     cloud,
			name:      name,
			probe:     probe,
		},
	}

	job.durationBuffer = cfg.PingCheckInterval / missingPingDurationBufferDiv
	job.alertTrigger = fmt.Sprintf(missingPingAlertTrigger, cfg.PingLagThreshold.Seconds())

	return worker.NewTicker(ctx, &cfg.BaseConfig, job, cfg.PingCheckInterval)
}
