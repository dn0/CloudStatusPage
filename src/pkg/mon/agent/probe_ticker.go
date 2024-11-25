package agent

import (
	"context"
	"time"

	"cspage/pkg/msg"
	"cspage/pkg/pb"
	"cspage/pkg/worker"
)

type Probe[T Cloud] interface {
	String() string
	Start(context.Context)
	Do(context.Context) []*pb.Result
	Stop(context.Context)
}

type Worker = worker.Daemon

func NewProbeTicker[T Cloud](
	ctx context.Context,
	cfg *CloudConfig[T],
	topic *msg.PubsubTopic,
	interval time.Duration,
	probe Probe[T],
) *Worker {
	job := NewProbeJob(cfg, msg.NewPublisher(&cfg.PubsubPublisherConfig, topic), probe, true)
	return worker.NewTicker(ctx, &cfg.BaseConfig, job, cfg.ProbeInterval(interval))
}
