package vpc

import (
	"context"
	"time"

	"github.com/google/uuid"

	"cspage/pkg/mon/agent"
	"cspage/pkg/msg"
	"cspage/pkg/worker"
)

type PingProbeJob[T agent.Cloud] struct {
	*agent.ProbeJob[T]
	Probe *PingProbe[T]
}

// NewPingProbeJob is used for when PingProbe is used as a side job in another job.
func NewPingProbeJob[T agent.Cloud](
	cfg *agent.CloudConfig[T],
	topic *msg.PubsubTopic,
	probe *PingProbe[T],
) *PingProbeJob[T] {
	job := agent.NewProbeJob(cfg, msg.NewPublisher(&cfg.PubsubPublisherConfig, topic), probe, false)
	return &PingProbeJob[T]{ProbeJob: job, Probe: probe}
}

func (j *PingProbeJob[T]) Run(ctx context.Context, host, addr string) error {
	tick := worker.Tick{
		ID:    uuid.Must(uuid.NewV7()).String(),
		Time:  time.Now(),
		Drift: 0,
	}

	j.Probe.Host = host
	j.Probe.Addr = addr

	// We are doing this for Azure only to avoid negative DNS caching issues.
	// Azure after a VM is created does not update the DNS (unlike AWS and GCP),
	select {
	case <-ctx.Done():
		//nolint:wrapcheck // The error is ignored by the caller anyway.
		return ctx.Err()
	case <-time.After(j.Probe.cfg.ResolveDelay):
		break
	}

	result, err := j.Probe.run(ctx)
	if err != nil {
		// The error can come either from pingPong() which can be a
		//  - resolve timeout,
		//  - no such host,
		// or from pingRun() which can be one or more:
		//  - packet loss,
		//  - i/o timeout.
		// In any of those cases it meas that the host is unavailable and the parent job did a poor job
		// creating and keeping the host up&running, i.e. it should deal with the error. It also means
		// that we should not publish such result as it has no value when the host is not available.
		return err
	}

	j.PublishResult(ctx, tick, result)

	return nil
}
