package agent

import (
	"context"

	"cspage/pkg/msg"
	"cspage/pkg/worker"
)

func NewPingTicker(ctx context.Context, cfg *Config, topic *msg.PubsubTopic) *worker.Daemon {
	job := &pingJob{
		cfg:       cfg,
		publisher: msg.NewPublisher(&cfg.PubsubPublisherConfig, topic),
	}
	return worker.NewTicker(ctx, &cfg.BaseConfig, job, cfg.PingInterval)
}
