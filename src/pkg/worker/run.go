package worker

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type Worker interface {
	Run(context.Context, *Group)
	Enabled() bool
}

func Run(mainCtx context.Context, log *slog.Logger, closeFun func(), workers ...Worker) {
	log.Info("Service is starting...")

	numWorkers := 0
	for _, worker := range workers {
		if worker.Enabled() {
			numWorkers++
		}
	}

	ctx, cancel := context.WithCancel(mainCtx)
	wg := newGroup(numWorkers)

	for _, worker := range workers {
		go worker.Run(ctx, wg)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown

	log.Info("Service is shutting down...")
	cancel()
	wg.wait()
	closeFun()
	log.Info("Service stopped")
}
