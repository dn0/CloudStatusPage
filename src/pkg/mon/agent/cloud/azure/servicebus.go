package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"

	"cspage/pkg/mon/agent"
	"cspage/pkg/mon/agent/cloud/azure/servicebus"
	"cspage/pkg/msg"
)

const (
	serviceBusURI = "%s.servicebus.windows.net"
)

func NewServiceBusClient(
	cfg *agent.CloudConfig[agent.Azure],
	cred *azidentity.DefaultAzureCredential,
) *azservicebus.Client {
	fullyQualifiedNamespace := fmt.Sprintf(serviceBusURI, cfg.Cloud.ServiceBusNamespace)
	client, err := azservicebus.NewClient(fullyQualifiedNamespace, cred, nil)
	if err != nil {
		agent.Die("Could not initialize Azure storage client", "err", err)
	}
	return client
}

func NewServiceBusQueueMessageProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.Azure],
	topic *msg.PubsubTopic,
	busClient *azservicebus.Client,
) *agent.Worker {
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeInterval(cfg.Cloud.ServiceBusQueueMessageInterval),
		servicebus.NewQueueMessageProbe(cfg, busClient),
	)
}
