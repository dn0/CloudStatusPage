package vpc

import (
	"context"
	"log/slog"
	"sync"

	"cspage/pkg/http"
	"cspage/pkg/mon/agent"
	"cspage/pkg/mon/agent/cloud/common"
	"cspage/pkg/pb"
)

type AgentPingProbe[T agent.Cloud] struct {
	cfg     *agent.Config
	name    string
	actions map[string]*PingAction
	client  *http.Client
	log     *slog.Logger
}

func NewAgentPingProbe[T agent.Cloud](
	cfg *agent.Config,
	pname string,
	actions map[string]*PingAction,
) *AgentPingProbe[T] {
	return &AgentPingProbe[T]{
		cfg:     cfg,
		name:    pname,
		actions: actions,
		client:  http.NewClient(),
		log:     slog.With("probe", pname),
	}
}

func (p *AgentPingProbe[T]) String() string {
	return p.name
}

func (p *AgentPingProbe[T]) Start(_ context.Context) {
	if p.cfg.VPCInterPing.Count < 1 {
		p.log.Info("Probe suspended", "actions", p.actions)
	} else {
		p.log.Info("Probe initialized", "actions", p.actions)
	}
}

//nolint:varnamelen // Variable a(gent) looks OK to me in this context.
func (p *AgentPingProbe[T]) Do(ctx context.Context) []*pb.Result {
	agents, err := common.GetRunningAgents(ctx, p.cfg, p.client)
	if err != nil {
		p.log.Error("Could not fetch list of running agents", "err", err)
		return nil
	}

	var wg sync.WaitGroup
	probes := make([]*PingProbe[T], len(agents))
	p.log.Debug("Agent ping probes starting...", "agents", agents)

	for i, a := range agents {
		region := a.GetCloudRegion()
		if p.cfg.Env.Region == region {
			continue
		}

		action, ok := p.actions[region]
		if !ok {
			p.log.Error(
				"Could not find agent ping action",
				"region", region,
				"hostname", a.GetHostname(),
				"addr", a.GetIpAddress(),
			)
			continue
		}
		probes[i] = NewPingProbe[T](&p.cfg.VPCInterPing, p.name+":"+region, a.GetHostname(), a.GetIpAddress(), action)
		wg.Add(1)
		go probes[i].Run(ctx, &wg)
	}

	wg.Wait()
	p.log.Debug("Agent ping probes finished", "agents", agents)
	res := make([]*pb.Result, len(probes))

	for i, probe := range probes {
		if probe != nil && len(probe.Results) > 0 {
			res[i] = probe.Results[0]
		}
	}

	return res
}

func (p *AgentPingProbe[T]) Stop(_ context.Context) {}
