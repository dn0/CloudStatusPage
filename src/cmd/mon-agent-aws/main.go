package main

import (
	"context"
	"log/slog"

	"cspage/pkg/http"
	"cspage/pkg/mon/agent"
	"cspage/pkg/mon/agent/cloud/aws"
	"cspage/pkg/msg"
	"cspage/pkg/worker"
)

const svcName = "mon-agent-aws"

func main() {
	ctx := context.Background()
	cfg := agent.NewConfig[agent.AWS]()

	pubsubClient := msg.NewPubsubClient(ctx, cfg.PubsubProjectID, msg.NewPubsubConfig(&cfg.PubsubPublisherConfig, nil))
	pubsubPingTopic := msg.NewPubsubTopic(pubsubClient, cfg.PubsubPingTopic)
	pubsubProbeTopic := msg.NewPubsubTopic(pubsubClient, cfg.PubsubProbeTopic)

	awsConfig := aws.LoadConfig(ctx, cfg)

	worker.Run(
		ctx,
		slog.With("name", svcName, "env", cfg.Env),
		func() {
			msg.ClosePubsubTopic(pubsubProbeTopic)
			msg.ClosePubsubTopic(pubsubPingTopic)
			msg.ClosePubsubClient(pubsubClient)
		},
		http.NewSimpleServer(ctx, &cfg.HTTPConfig, cfg.Debug),
		agent.NewPingTicker(ctx, &cfg.Config, pubsubPingTopic),

		aws.NewEC2VMStandardProbe(ctx, cfg, pubsubProbeTopic, awsConfig),
		aws.NewEC2VMSpotProbe(ctx, cfg, pubsubProbeTopic, awsConfig),
		aws.NewEC2VMMetadataProbe(ctx, cfg, pubsubProbeTopic),
		aws.NewEC2EBSSnapshotProbe(ctx, cfg, pubsubProbeTopic, awsConfig),

		aws.NewS3ObjectProbe(ctx, cfg, pubsubProbeTopic, awsConfig),
		aws.NewS3BucketProbe(ctx, cfg, pubsubProbeTopic, awsConfig),

		aws.NewSQSMessageProbe(ctx, cfg, pubsubProbeTopic, awsConfig),

		aws.NewVPCInterPingProbe(ctx, cfg, pubsubProbeTopic),
	)
}
