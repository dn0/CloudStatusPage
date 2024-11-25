package agent

import (
	"context"
	"log/slog"
	"math/rand"
	"sync"
	"time"

	"cspage/pkg/msg"
	"cspage/pkg/pb"
	"cspage/pkg/worker"
)

const (
	minProbeStartDelaySeconds = 1
	maxProbeStartDelaySeconds = 10
)

type ProbeJob[T Cloud] struct {
	cfg        *CloudConfig[T]
	publisher  msg.Publisher
	probe      Probe[T]
	startSleep time.Duration

	shutdownMutex sync.RWMutex
	shutdownCause error
}

func NewProbeJob[T Cloud](cfg *CloudConfig[T], publisher msg.Publisher, probe Probe[T], startDelay bool) *ProbeJob[T] {
	startSleep := time.Duration(0)
	if startDelay {
		//nolint:gosec // So that not all jobs (probes) start at the same time.
		startSleep = time.Duration(minProbeStartDelaySeconds+rand.Intn(maxProbeStartDelaySeconds)) * time.Second
	}
	return &ProbeJob[T]{
		cfg:        cfg,
		publisher:  publisher,
		probe:      probe,
		startSleep: startSleep,
	}
}

func (j *ProbeJob[T]) String() string {
	return worker.TickerJobPrefix + j.probe.String()
}

func (j *ProbeJob[T]) PreStart(_ context.Context) {}

func (j *ProbeJob[T]) Start(ctx context.Context) {
	time.Sleep(j.startSleep)
	j.probe.Start(ctx)
}

func (j *ProbeJob[T]) Do(ctx context.Context, tick worker.Tick) {
	j.PublishResult(ctx, tick, j.probe.Do(ctx))
}

func (j *ProbeJob[T]) Stop(ctx context.Context) {
	defer j.publisher.Close()

	j.probe.Stop(ctx)
}

func (j *ProbeJob[T]) Shutdown(cause error) {
	j.shutdownMutex.Lock()
	defer j.shutdownMutex.Unlock()
	j.shutdownCause = cause
}

func (j *ProbeJob[T]) isShuttingDown() error {
	j.shutdownMutex.RLock()
	defer j.shutdownMutex.RUnlock()
	return j.shutdownCause
}

func (j *ProbeJob[T]) PublishResult(ctx context.Context, tick worker.Tick, result []*pb.Result) {
	if result == nil {
		return // Probe didn't produce any results
	}
	if cause := j.isShuttingDown(); cause != nil {
		slog.Warn("Ignoring task result", "job", j.String(), "tick", tick, "cause", cause)
		return // Main loop was interrupted => the result can be corrupted
	}
	if ctx.Err() != nil {
		slog.Warn("Ignoring task result", "job", j.String(), "tick", tick, "cause", context.Cause(ctx))
		return // Task was interrupted => the result can be corrupted
	}
	// Publish & forget => not using task context
	j.publishMessage(tick, result)
}

func (j *ProbeJob[T]) publishMessage(tick worker.Tick, result []*pb.Result) {
	message := &pb.Probe{
		Job:    newJobMessage(&j.cfg.Env, &tick, j.probe.String()),
		Result: []*pb.Result{},
	}
	for _, res := range result {
		if res == nil || res.GetStatus() == pb.ResultStatus_RESULT_UNKNOWN {
			continue
		}
		message.Result = append(message.Result, res)
		if res.Failed() {
			message.Job.Errors++
		}
	}
	// NOTE: this means publish & forget as we are not waiting for the result
	//       therefore we can't use the task context here
	j.publisher.Publish(context.Background(), msg.NewMessage(message, msg.NewAttrs(
		msg.TypeProbe,
		j.cfg.Env.Cloud,
		j.cfg.Env.Region,
	)))
}
