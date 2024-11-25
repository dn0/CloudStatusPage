package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"

	"cspage/pkg/mon/agent"
	"cspage/pkg/mon/agent/cloud/azure/storage"
	"cspage/pkg/msg"
)

const (
	storageServiceURL = "https://%s.blob.core.windows.net"
)

func NewStorageClient(
	cfg *agent.CloudConfig[agent.Azure],
	cred *azidentity.DefaultAzureCredential,
) *azblob.Client {
	serviceURL := fmt.Sprintf(storageServiceURL, cfg.Cloud.StorageAccountName)
	client, err := azblob.NewClient(serviceURL, cred, nil)
	if err != nil {
		agent.Die("Could not initialize Azure storage client", "err", err)
	}
	return client
}

func NewStorageObjectProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.Azure],
	topic *msg.PubsubTopic,
	storageClient *azblob.Client,
) *agent.Worker {
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeInterval(cfg.Cloud.StorageBlobProbeInterval),
		storage.NewBlobProbe(cfg, storageClient),
	)
}

func NewStorageContainerProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.Azure],
	topic *msg.PubsubTopic,
	storageClient *azblob.Client,
) *agent.Worker {
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeInterval(cfg.Cloud.StorageContainerProbeInterval),
		storage.NewContainerProbe(cfg, storageClient),
	)
}
