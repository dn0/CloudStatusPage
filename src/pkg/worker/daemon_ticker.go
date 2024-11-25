package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"cspage/pkg/config"
)

const (
	TickerJobPrefix = "ticker:"
)

type ticker struct {
	job        TickerJob
	timeout    time.Duration
	errTimeout error
	interval   time.Duration
}

func NewTicker(
	ctx context.Context,
	cfg *config.BaseConfig,
	job TickerJob,
	interval time.Duration,
	loopDesc ...any,
) *Daemon {
	if interval < 0 {
		config.Die("Ticker interval cannot be negative", "job", job.String())
	}
	loopDesc = append([]any{
		"job", job.String(),
		"interval", interval,
	}, loopDesc...)
	worker := &ticker{
		job:        job,
		timeout:    cfg.WorkerTaskTimeout,
		errTimeout: fmt.Errorf("%s: %w", job.String(), errTickerTaskTimeoutExceeded),
		interval:   interval,
	}
	return newDaemon(ctx, cfg, job, interval != 0, worker.loop, loopDesc...)
}

func (t *ticker) loop(loopCtx, taskCtx context.Context) {
	tickr := time.NewTicker(t.interval)
	prevT := time.Now()

	for {
		select {
		case <-loopCtx.Done():
			tickr.Stop()
			return
		case tickT := <-tickr.C:
			// I know, this look weird, but we need to double-check whether the main loop does not want us to stop,
			// i.e. loopCancel wasn't called. This helps to protect against situations where we have to stop but
			// the tickr is "overflown" due to a current long-running task and already wants to run the next tick job.
			select {
			case <-loopCtx.Done():
				tickr.Stop()
				return
			default:
				break
			}

			tick := Tick{
				ID:    uuid.Must(uuid.NewV7()).String(),
				Time:  tickT,
				Drift: tickT.Sub(prevT) - t.interval,
			}
			log := slog.With("job", t.job.String())
			log.Debug("Worker task started", "tick", tick)
			innerTaskCtx, innerTaskCancel := context.WithTimeoutCause(taskCtx, t.timeout, t.errTimeout)
			t.job.Do(innerTaskCtx, tick)
			innerTaskCancel()
			tick.Took = time.Since(tickT)
			log.Debug("Worker task finished", "tick", tick)
			prevT = tickT
		}
	}
}
