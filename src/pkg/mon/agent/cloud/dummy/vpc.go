package dummy

import (
	"context"

	"cspage/pkg/mon/agent"
	"cspage/pkg/mon/agent/cloud/dummy/vpc"
	"cspage/pkg/msg"
)

func NewAgentPingProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.Dummy],
	topic *msg.PubsubTopic,
) *agent.Worker {
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeInterval(cfg.Cloud.AgentPingProbeInterval),
		vpc.NewAgentPingProbe(&cfg.Config),
	)
}

func NewInternetPingProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.Dummy],
	topic *msg.PubsubTopic,
) *agent.Worker {
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeInterval(cfg.Cloud.InternetPingProbeInterval),
		vpc.NewInternetPingProbe(&cfg.Config),
	)
}
