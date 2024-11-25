package gcp

import (
	"context"

	"cspage/pkg/mon/agent"
	"cspage/pkg/mon/agent/cloud/gcp/vpc"
	"cspage/pkg/msg"
)

func NewVPCInterPingProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.GCP],
	topic *msg.PubsubTopic,
) *agent.Worker {
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeInterval(cfg.Cloud.VPCInterPingInterval),
		vpc.NewInterPingProbe(cfg),
	)
}
