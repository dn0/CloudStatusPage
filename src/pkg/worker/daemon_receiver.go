package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"cspage/pkg/config"
)

const (
	ReceiverJobPrefix = "receiver:"
)

type receiver[T any] struct {
	job     ReceiverJob[T]
	channel <-chan T
}

func NewReceiver[T any](
	ctx context.Context,
	cfg *config.BaseConfig,
	job ReceiverJob[T],
	channel <-chan T,
	loopDesc ...any,
) *Daemon {
	loopDesc = append([]any{
		"job", job.String(),
		"channel", fmt.Sprintf("%T", channel),
	}, loopDesc...)
	worker := &receiver[T]{
		job:     job,
		channel: channel,
	}
	return newDaemon(ctx, cfg, job, job.Enabled(), worker.loop, loopDesc...)
}

func (r *receiver[T]) loop(loopCtx, taskCtx context.Context) {
	log := slog.With("job", r.job.String())
	for {
		select {
		case <-loopCtx.Done():
			return
		case event := <-r.channel:
			now := time.Now()
			log.Debug("Received event", "event", event)
			err := r.job.Process(taskCtx, event)
			if err == nil {
				log.Debug("Processed event", "event", event, "took", time.Since(now))
			} else {
				log.Error("Failed to process event", "event", event, "took", time.Since(now), "err", err)
			}
		}
	}
}
