package analyst

import (
	"context"
	"log/slog"
	"time"

	"cspage/pkg/data"
	"cspage/pkg/db"
	"cspage/pkg/msg"
	"cspage/pkg/worker"
)

func NewProbeTickers(
	ctx context.Context,
	cfg *Config,
	dbc *db.Clients,
	cloud string,
	topic *msg.PubsubTopic,
) []*worker.Daemon {
	probes, err := data.GetProbeDefinitions(ctx, dbc.Read, cloud)
	if err != nil {
		Die("Could not fetch probes from config table", "cloud", cloud, "err", err)
	}

	workers := make([]*worker.Daemon, len(probes))
	for i, probed := range probes {
		workers[i] = newProbeTicker(ctx, cfg, dbc, topic, cloud, probed)
	}
	return workers
}

func newProbeTicker(
	ctx context.Context,
	cfg *Config,
	dbc *db.Clients,
	topic *msg.PubsubTopic,
	cloud string,
	probe *data.ProbeDefinition,
) *worker.Daemon {
	name := worker.TickerJobPrefix + "probe:" + probe.Name
	job := &probeJob{
		analystJob: analystJob{
			cfg:       cfg,
			dbc:       dbc,
			publisher: msg.NewPublisher(&cfg.PubsubPublisherConfig, topic),
			baseLog:   slog.With("job", name, "probe", probe.Name),
			cloud:     cloud,
			name:      name,
			probe:     probe,
		},
	}

	// This should help us get a roughly uniform number of probe results across different probes
	switch probe.Config.IntervalType {
	case data.ProbeIntervalStandard:
		job.checkWindow = cfg.ProbeCheckWindow
		job.zscoreWindow = cfg.ProbeZscoreWindow
	case data.ProbeIntervalLong:
		job.checkWindow = cfg.ProbeLongCheckWindow
		job.zscoreWindow = cfg.ProbeLongZscoreWindow
	default:
		DieLog(job.baseLog, "Unexpected probe interval type", "probe_interval_type", probe.Config.IntervalType)
	}
	//nolint:gosec // This integer should not overflow.
	job.failureThreshold = uint32(cfg.ProbeFailureThreshold)
	job.zscoreThresholdOn = float32(cfg.ProbeZscoreThresholdOn)
	job.zscoreThresholdOff = float32(cfg.ProbeZscoreThresholdOff)
	job.distanceFromAvgDiv = time.Duration(cfg.ProbeDistanceFromAvgDivThreshold)
	job.distanceFromAvgMin = cfg.ProbeDistanceFromAvgMinThreshold

	return worker.NewTicker(ctx, &cfg.BaseConfig, job, cfg.ProbeCheckInterval)
}
