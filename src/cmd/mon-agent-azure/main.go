package main

import (
	"context"
	"log/slog"

	"cspage/pkg/http"
	"cspage/pkg/mon/agent"
	"cspage/pkg/mon/agent/cloud/azure"
	"cspage/pkg/msg"
	"cspage/pkg/worker"
)

const svcName = "mon-agent-azure"

func main() {
	ctx := context.Background()
	cfg := agent.NewConfig[agent.Azure]()

	pubsubClient := msg.NewPubsubClient(ctx, cfg.PubsubProjectID, msg.NewPubsubConfig(&cfg.PubsubPublisherConfig, nil))
	pubsubPingTopic := msg.NewPubsubTopic(pubsubClient, cfg.PubsubPingTopic)
	pubsubProbeTopic := msg.NewPubsubTopic(pubsubClient, cfg.PubsubProbeTopic)

	azCredential := azure.NewCredential()
	storageClient := azure.NewStorageClient(cfg, azCredential)
	computeFactory := azure.NewComputeClientFactory(cfg, azCredential)
	servicebusClient := azure.NewServiceBusClient(cfg, azCredential)

	worker.Run(
		ctx,
		slog.With("name", svcName, "env", cfg.Env),
		func() {
			_ = servicebusClient.Close(ctx)

			msg.ClosePubsubTopic(pubsubProbeTopic)
			msg.ClosePubsubTopic(pubsubPingTopic)
			msg.ClosePubsubClient(pubsubClient)
		},
		http.NewSimpleServer(ctx, &cfg.HTTPConfig, cfg.Debug),
		agent.NewPingTicker(ctx, &cfg.Config, pubsubPingTopic),

		azure.NewComputeVMStandardProbe(ctx, cfg, pubsubProbeTopic, computeFactory),
		azure.NewComputeVMSpotProbe(ctx, cfg, pubsubProbeTopic, computeFactory),
		azure.NewComputeVMMetadataProbe(ctx, cfg, pubsubProbeTopic),
		azure.NewComputeVHDSnapshotProbe(ctx, cfg, pubsubProbeTopic, computeFactory),

		azure.NewStorageObjectProbe(ctx, cfg, pubsubProbeTopic, storageClient),
		azure.NewStorageContainerProbe(ctx, cfg, pubsubProbeTopic, storageClient),

		azure.NewServiceBusQueueMessageProbe(ctx, cfg, pubsubProbeTopic, servicebusClient),

		azure.NewVPCInterPingProbe(ctx, cfg, pubsubProbeTopic),
	)
}
