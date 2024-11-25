package vpc

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"cspage/pkg/mon/agent"
	"cspage/pkg/pb"
)

const (
	PingSize64B   = 56
	PingSize1500B = 1472
)

type PingProbe[T agent.Cloud] struct {
	cfg     *agent.PingConfig
	name    string
	actions []*PingAction

	Host string
	Addr string

	Results []*pb.Result
}

func NewPingProbe[T agent.Cloud](cfg *agent.PingConfig, pname, host, addr string, actions ...*PingAction) *PingProbe[T] {
	return &PingProbe[T]{
		cfg:     cfg,
		name:    pname,
		actions: actions,
		Host:    host,
		Addr:    addr,
	}
}

func (p *PingProbe[T]) String() string {
	return p.name
}

func (p *PingProbe[T]) Start(_ context.Context) {
	if p.cfg.Count < 1 {
		p.log().Info("Probe suspended", "actions", p.actions)
	} else {
		p.log().Info("Probe initialized", "actions", p.actions)
	}
}

func (p *PingProbe[T]) Do(ctx context.Context) []*pb.Result {
	results, _ := p.run(ctx)
	return results
}

func (p *PingProbe[T]) Stop(_ context.Context) {}

func (p *PingProbe[T]) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	p.Results, _ = p.run(ctx)
}

func (p *PingProbe[T]) run(ctx context.Context) ([]*pb.Result, error) {
	if p.cfg.Count < 1 {
		return nil, nil
	}

	res := make([]*pb.Result, len(p.actions))
	var errs []error

	for i, action := range p.actions {
		res[i] = pb.NewResult(action.Id)
		if i == 0 { // warm up
			if ret, err := p.pingPong(ctx, action); err != nil {
				res[i].Store(ret, err)
				return res, err
			}
		}
		ret, err := p.pingRun(ctx, action)
		res[i].Store(ret, err)
		if err != nil {
			errs = append(errs, err)
		}
	}

	return res, errors.Join(errs...)
}

func (p *PingProbe[T]) pingRun(ctx context.Context, action *PingAction) (pb.ResultTime, error) {
	ret := pb.EmptyResultTime()
	p.log().Debug("Ping starting...", "pinger", action.String())
	stats, err := action.pingRun(ctx, p.cfg, p.Host, p.Addr)

	if err == nil {
		p.log().Debug("Ping finished", "pinger", action.String(), "stats", stats)
	} else {
		p.log().Error("Ping failed", "pinger", action.String(), "stats", stats, "err", err)
	}
	if stats != nil {
		ret.Took = stats.AvgRtt
	}

	action.pingReset()
	return ret, err
}

func (p *PingProbe[T]) pingPong(ctx context.Context, action *PingAction) (pb.ResultTime, error) {
	ret := pb.EmptyResultTime()
	p.log().Debug("Ping pong starting...", "pinger", action.String())
	stats, err := action.pingPong(ctx, p.cfg, p.Host, p.Addr)

	if err == nil {
		p.log().Debug("Ping pong finished", "pinger", action.String(), "stats", stats)
	} else {
		p.log().Error("Ping pong failed", "pinger", action.String(), "stats", stats, "err", err)
	}

	return ret, err
}

func (p *PingProbe[T]) log() *slog.Logger {
	return slog.With("probe", p.String(), "host", p.Host, "addr", p.Addr)
}
