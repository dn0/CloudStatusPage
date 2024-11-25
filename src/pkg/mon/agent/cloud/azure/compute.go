package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"

	"cspage/pkg/mon/agent"
	"cspage/pkg/mon/agent/cloud/azure/compute"
	"cspage/pkg/mon/agent/cloud/azure/vpc"
	"cspage/pkg/msg"
)

func NewComputeClientFactory(
	cfg *agent.CloudConfig[agent.Azure],
	cred *azidentity.DefaultAzureCredential,
) *armcompute.ClientFactory {
	factory, err := armcompute.NewClientFactory(cfg.Cloud.SubscriptionID, cred, nil)
	if err != nil {
		agent.Die("Could not initialize Azure client factory", "err", err)
	}
	return factory
}

func NewComputeVMStandardProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.Azure],
	topic *msg.PubsubTopic,
	factory *armcompute.ClientFactory,
) *agent.Worker {
	pingJob := vpc.NewIntraPingProbeJob(cfg, topic, cfg.Cloud.ComputeVMPrefix+"-TBD")
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeLongInterval(cfg.Cloud.ComputeVMProbeInterval),
		compute.NewVMProbe(cfg, factory, compute.VMProvisioningModelStandard, pingJob),
	)
}

func NewComputeVMSpotProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.Azure],
	topic *msg.PubsubTopic,
	factory *armcompute.ClientFactory,
) *agent.Worker {
	pingJob := vpc.NewIntraPingProbeJob(cfg, topic, cfg.Cloud.ComputeVMSpotPrefix+"-TBD")
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeLongInterval(cfg.Cloud.ComputeVMSpotProbeInterval),
		compute.NewVMProbe(cfg, factory, compute.VMProvisioningModelSpot, pingJob),
	)
}

func NewComputeVMMetadataProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.Azure],
	topic *msg.PubsubTopic,
) *agent.Worker {
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeInterval(cfg.Cloud.ComputeVMMetadataProbeInterval),
		compute.NewVMMetadataProbe(cfg),
	)
}

func NewComputeVHDSnapshotProbe(
	ctx context.Context,
	cfg *agent.CloudConfig[agent.Azure],
	topic *msg.PubsubTopic,
	factory *armcompute.ClientFactory,
) *agent.Worker {
	return agent.NewProbeTicker(
		ctx,
		cfg,
		topic,
		cfg.ProbeLongInterval(cfg.Cloud.ComputeVHDSnapshotProbeInterval),
		compute.NewVHDSnapshotProbe(cfg, factory),
	)
}
