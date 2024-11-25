package main

import (
	"context"
	"log/slog"

	"cspage/pkg/http"
	"cspage/pkg/mon/agent"
	"cspage/pkg/mon/agent/cloud/gcp"
	"cspage/pkg/msg"
	"cspage/pkg/worker"
)

const svcName = "mon-agent-gcp"

func main() {
	ctx := context.Background()
	cfg := agent.NewConfig[agent.GCP]()

	pubsubClient := msg.NewPubsubClient(ctx, cfg.PubsubProjectID, msg.NewPubsubConfig(&cfg.PubsubPublisherConfig, nil))
	pubsubPingTopic := msg.NewPubsubTopic(pubsubClient, cfg.PubsubPingTopic)
	pubsubProbeTopic := msg.NewPubsubTopic(pubsubClient, cfg.PubsubProbeTopic)

	storageClient := gcp.NewStorageClient(ctx)

	worker.Run(
		ctx,
		slog.With("name", svcName, "env", cfg.Env),
		func() {
			gcp.CloseStorageClient(storageClient)

			msg.ClosePubsubTopic(pubsubProbeTopic)
			msg.ClosePubsubTopic(pubsubPingTopic)
			msg.ClosePubsubClient(pubsubClient)
		},
		http.NewSimpleServer(ctx, &cfg.HTTPConfig, cfg.Debug),
		agent.NewPingTicker(ctx, &cfg.Config, pubsubPingTopic),

		gcp.NewComputeVMStandardProbe(ctx, cfg, pubsubProbeTopic),
		gcp.NewComputeVMSpotProbe(ctx, cfg, pubsubProbeTopic),
		gcp.NewComputeVMMetadataProbe(ctx, cfg, pubsubProbeTopic),
		gcp.NewComputeDiskSnapshotProbe(ctx, cfg, pubsubProbeTopic),

		gcp.NewStorageObjectProbe(ctx, cfg, pubsubProbeTopic, storageClient),
		gcp.NewStorageBucketProbe(ctx, cfg, pubsubProbeTopic, storageClient),

		gcp.NewPubsubMessageProbe(ctx, cfg, pubsubProbeTopic),

		gcp.NewVPCInterPingProbe(ctx, cfg, pubsubProbeTopic),
	)
}
