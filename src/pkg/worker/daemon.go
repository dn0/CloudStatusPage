package worker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"cspage/pkg/config"
)

var (
	errPreStartTimeoutExceeded     = errors.New("job.PreStart() deadline exceeded")
	errStartTimeoutExceeded        = errors.New("job.Start() deadline exceeded")
	errStopTimeoutExceeded         = errors.New("job.Stop() deadline exceeded")
	errTickerTaskTimeoutExceeded   = errors.New("job.Do() deadline exceeded")
	errConsumerTaskTimeoutExceeded = errors.New("job.Process() deadline exceeded")
	errShutdownInProgress          = errors.New("job.Shutdown() in progress")
)

type loopFunc = func(context.Context, context.Context)

type Daemon struct {
	ctx      context.Context // workerCtx (usually mainCtx)
	cfg      *config.BaseConfig
	job      job
	log      *slog.Logger
	enabled  bool
	loopFunc loopFunc
	loopSync sync.WaitGroup
}

func newDaemon(
	ctx context.Context,
	cfg *config.BaseConfig,
	job job,
	enabled bool,
	loopFunc loopFunc,
	loopDesc ...any,
) *Daemon {
	return &Daemon{
		ctx:      ctx,
		cfg:      cfg,
		job:      job,
		log:      slog.With(loopDesc...),
		enabled:  enabled,
		loopFunc: loopFunc,
	}
}

func (d *Daemon) Enabled() bool {
	return d.enabled
}

//nolint:contextcheck // runnerCtx is used "only" to process a shutdown signal.
func (d *Daemon) Run(runnerCtx context.Context, wg *Group) {
	if !d.Enabled() {
		d.log.Info("Worker is disabled")
		return
	}

	defer wg.Done()
	loopCtx, loopCancel := context.WithCancelCause(d.ctx)
	taskCtx, taskCancel := context.WithCancelCause(d.ctx)
	defer d.shutdown(wg, taskCtx, taskCancel, loopCancel)
	d.loopSync.Add(1)

	go func() {
		defer d.loopSync.Done()
		defer taskCancel(nil)
		d.log.Info("Worker is starting...")

		preStartCtx, preStartCancel := d.contextWithTimeout(d.cfg.WorkerStartTimeout, errPreStartTimeoutExceeded)
		d.job.PreStart(preStartCtx)
		preStartCancel()

		wg.Start()

		<-wg.Ready

		startCtx, startCancel := d.contextWithTimeout(d.cfg.WorkerStartTimeout, errStartTimeoutExceeded)
		d.job.Start(startCtx)
		startCancel()

		d.log.Info("Worker started")
		d.loopFunc(loopCtx, taskCtx) // 2 x mainCtx with cancel, runnerCtx is deliberately not passed down
	}()

	<-runnerCtx.Done()
}

//nolint:contextcheck,revive // The context is inherited from runnerCtx.
func (d *Daemon) shutdown(wg *Group, taskCtx context.Context, taskCancel, loopCancel context.CancelCauseFunc) {
	d.log.Info("Worker is shutting down...")
	wg.Stop()

	// Cancel the main loop and inform the running job
	d.job.Shutdown(errShutdownInProgress)
	loopCancel(fmt.Errorf("loop:%s: %w", d.job.String(), errShutdownInProgress))
	// A final task in the loop could be still in progress, let's give it a few seconds to finish
	select {
	case <-taskCtx.Done():
		d.log.Info("Worker task finished")
		break
	case <-time.After(d.cfg.WorkerShutdownTimeout):
		d.log.Warn("Worker task cancelled")
		taskCancel(fmt.Errorf("task:%s: %w", d.job.String(), errShutdownInProgress))
		break
	}

	d.loopSync.Wait()

	stopCtx, stopCancel := d.contextWithTimeout(d.cfg.WorkerStopTimeout, errStopTimeoutExceeded)
	defer stopCancel()
	d.job.Stop(stopCtx)
	d.log.Info("Worker stopped")
}

func (d *Daemon) contextWithTimeout(timeout time.Duration, cause error) (context.Context, context.CancelFunc) {
	return context.WithTimeoutCause(d.ctx, timeout, fmt.Errorf("%s: %w", d.job.String(), cause))
}
