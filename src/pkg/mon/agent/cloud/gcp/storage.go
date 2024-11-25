package gcp

import (
	"context"
	"log/slog"

	cloudStorage "cloud.google.com/go/storage"

	"cspage/pkg/mon/agent"
	"cspage/pkg/mon/agent/cloud/gcp/storage"
	"cspage/pkg/msg"
)

func CloseStorageClient(client *cloudStorage.Client) {
	slog.Debug("Closing GCP storage client...")
	if client != nil {
		if err := client.Close(); err != nil {
			slog.Error("Could not close GCP storage client", "err", err)
		}
	}
}

func NewStorageClient(ctx context.Context) *cloudStorage.Client {
	client, err := cloudStorage.NewClient(ctx)
	if err != nil {
		agent.Die("Could not initialize GCP storage client", "err", err)
	}
	return client
}

func NewStorageObjectProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.GCP],
	topic *msg.PubsubTopic,
	storageClient *cloudStorage.Client,
) *agent.Worker {
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeInterval(cfg.Cloud.StorageObjectProbeInterval),
		storage.NewObjectProbe(cfg, storageClient),
	)
}

func NewStorageBucketProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.GCP],
	topic *msg.PubsubTopic,
	storageClient *cloudStorage.Client,
) *agent.Worker {
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeInterval(cfg.Cloud.StorageBucketProbeInterval),
		storage.NewBucketProbe(cfg, storageClient),
	)
}
