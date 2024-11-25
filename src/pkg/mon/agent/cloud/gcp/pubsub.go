package gcp

import (
	"context"

	"cspage/pkg/mon/agent"
	"cspage/pkg/mon/agent/cloud/gcp/pubsub"
	"cspage/pkg/msg"
)

func NewPubsubMessageProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.GCP],
	topic *msg.PubsubTopic,
) *agent.Worker {
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeInterval(cfg.Cloud.PubsubMessageInterval),
		pubsub.NewMessageProbe(cfg),
	)
}
