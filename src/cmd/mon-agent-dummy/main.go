package main

import (
	"context"
	"log/slog"

	"cspage/pkg/http"
	"cspage/pkg/mon/agent"
	"cspage/pkg/mon/agent/cloud/dummy"
	"cspage/pkg/msg"
	"cspage/pkg/worker"
)

const svcName = "mon-agent-dummy"

func main() {
	ctx := context.Background()
	cfg := agent.NewConfig[agent.Dummy]()

	pubsubClient := msg.NewPubsubClient(ctx, cfg.PubsubProjectID, nil)
	pubsubPingTopic := msg.NewPubsubTopic(pubsubClient, cfg.PubsubPingTopic)
	pubsubProbeTopic := msg.NewPubsubTopic(pubsubClient, cfg.PubsubProbeTopic)

	worker.Run(
		ctx,
		slog.With("name", svcName, "env", cfg.Env),
		func() {
			msg.ClosePubsubTopic(pubsubPingTopic)
			msg.ClosePubsubClient(pubsubClient)
		},
		http.NewSimpleServer(ctx, &cfg.HTTPConfig, cfg.Debug),
		agent.NewPingTicker(ctx, &cfg.Config, pubsubPingTopic),

		dummy.NewAgentPingProbe(ctx, cfg, pubsubProbeTopic),
		dummy.NewInternetPingProbe(ctx, cfg, pubsubProbeTopic),
	)
}
