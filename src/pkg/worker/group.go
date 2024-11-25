package worker

import (
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
)

type Group struct {
	Ready    chan struct{}
	Required int32
	running  atomic.Int32
	group    sync.WaitGroup
}

func newGroup(numWorkers int) *Group {
	wg := &Group{
		Ready:    make(chan struct{}),
		Required: int32(numWorkers), //nolint:gosec // This integer should never overflow.
	}
	wg.group.Add(numWorkers)
	return wg
}

func (g *Group) Done() {
	g.group.Done()
}

func (g *Group) wait() {
	g.group.Wait()
}

func (g *Group) Start() {
	g.running.Add(1)

	if g.running.Load() == g.Required {
		// The /healthz endpoint now returns 200
		slog.Info("All workers are almost ready", "workers", g.String())
		close(g.Ready)
	}
}

func (g *Group) Stop() {
	g.running.Add(-1)

	if g.running.Load() == 0 {
		slog.Info("All workers are almost done", "workers", g.String())
	}
}

func (g *Group) Healthy() bool {
	return g.running.Load() == g.Required
}

func (g *Group) String() string {
	return fmt.Sprintf("%d/%d", g.running.Load(), g.Required)
}
